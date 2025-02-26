package http

import (
	"bytes"
	"context"
	"io"
	"janan_csv_service/internal/csvprocessor/usecase"
	"janan_csv_service/pkg/progress"
	"janan_csv_service/pkg/redis"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CSVProcessorHandler struct {
	usecase  *usecase.CSVProcessorUsecase
	redis    *redis.RedisClient
	progress *progress.Tracker
}

func NewCSVProcessorHandler(usecase *usecase.CSVProcessorUsecase, redis *redis.RedisClient, progress *progress.Tracker) *CSVProcessorHandler {
	return &CSVProcessorHandler{usecase: usecase, redis: redis, progress: progress}
}

func (h *CSVProcessorHandler) UploadCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	// Read the file content into memory
	fileContent, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate a task ID
	taskID := h.usecase.GenerateTaskID(file.Filename)

	// Start processing the file in a goroutine
	go func(content []byte, fileName string, taskID string) {
		_, err := h.usecase.ProcessCSV(bytes.NewReader(content), fileName, taskID)
		if err != nil {
			log.Printf("Error processing file: %v", err)
			return
		}
		log.Printf("File processing completed for task ID: %s", taskID)
	}(fileContent, file.Filename, taskID)
	// Return the task ID to the client immediately
	c.JSON(http.StatusOK, gin.H{
		"message": "File upload started",
		"task_id": taskID, // Include the taskId in the response
	})
}

func (h *CSVProcessorHandler) CheckProgress(c *gin.Context) {
	taskID := c.Param("task_id")
	progress, err := h.redis.GetProgress(context.Background(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if progress == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

func (h *CSVProcessorHandler) StreamProgress(c *gin.Context) {
	taskID := c.Param("task_id")

	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "http://localhost:3000")

	// Create a channel to listen for progress updates
	progressChan := h.progress.Subscribe(taskID)
	defer h.progress.Unsubscribe(taskID)

	// Stream progress updates to the client
	for {
		select {
		case progress := <-progressChan:
			// Send progress update to the client
			c.SSEvent("message", gin.H{
				"processed": progress.Processed,
				"total":     progress.Total,
				"status":    progress.Status,
			})
			c.Writer.Flush()

			// Stop streaming if the task is completed
			if progress.Status == "completed" {
				return
			}
		case <-c.Writer.CloseNotify():
			// Stop streaming if the client disconnects
			return
		case <-time.After(30 * time.Second):
			// Send a heartbeat to keep the connection alive
			c.SSEvent("heartbeat", nil)
			c.Writer.Flush()
		}
	}
}
