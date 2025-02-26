package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"io"
	"janan_csv_service/internal/csvprocessor/repository"
	"janan_csv_service/internal/models"
	"janan_csv_service/pkg/progress"
	"janan_csv_service/pkg/redis"
	"log"
	"strconv"
	"sync"
	"time"
)

type CSVProcessorUsecase struct {
	repo     *repository.CSVProcessorRepository
	redis    *redis.RedisClient
	apiKey   string
	progress *progress.Tracker
}

func NewCSVProcessorUsecase(repo *repository.CSVProcessorRepository, redis *redis.RedisClient, apiKey string, progress *progress.Tracker) *CSVProcessorUsecase {
	return &CSVProcessorUsecase{repo: repo, redis: redis, apiKey: apiKey, progress: progress}
}

func (uc *CSVProcessorUsecase) ProcessCSV(file io.Reader, fileName string, taskID string) (string, error) {
	// Initialize progress
	uc.progress.UpdateProgress(taskID, 0, 0, "uploading")

	// Generate a hash for the file
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	fileHash := hex.EncodeToString(hasher.Sum(nil))

	// Reset the file pointer to the beginning
	if seeker, ok := file.(io.Seeker); ok {
		_, err := seeker.Seek(0, 0)
		if err != nil {
			log.Printf("Error resetting file pointer: %v", err)
			return "", err
		}
	}

	// Check if the file has been processed before
	ctx := context.Background()
	// existingHashes, err := uc.redis.GetFileHashes(ctx, uc.apiKey)
	// if err != nil {
	// 	return "", err
	// }

	// for _, hash := range existingHashes {
	// 	if hash == fileHash {
	// 		log.Printf("File %s has already been processed", fileName)
	// 		return "", nil
	// 	}
	// }

	// Process the file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading CSV file: %v", err)
		return "", err
	}

	// Check if the file is empty
	if len(records) == 0 {
		log.Printf("File %s is empty", fileName)
		return "", nil
	}

	totalRecords := len(records)

	// Update total records in progress
	uc.progress.UpdateProgress(taskID, 0, totalRecords, "processing")

	var wg sync.WaitGroup
	for _, record := range records {
		wg.Add(1)
		go func(record []string) {
			defer wg.Done()

			// Convert record[3] (grade) to float64
			grade, err := strconv.ParseFloat(record[3], 64)
			if err != nil {
				grade = 0.00
			}

			student := models.Student{
				StudentID:   record[0],
				StudentName: record[1],
				Subject:     record[2],
				Grade:       grade,
			}

			if err := uc.repo.SaveStudent(student); err != nil {
				log.Printf("Error saving student: %v", err)
			}

			// Update processed records in progress
			currentProgress, _ := uc.progress.GetProgress(taskID)
			uc.progress.UpdateProgress(taskID, currentProgress.Processed+1, currentProgress.Total, "processing")
		}(record)
	}

	wg.Wait()

	uc.progress.UpdateProgress(taskID, totalRecords, totalRecords, "completed")

	// Store the file hash under the API key in Redis
	if err := uc.redis.StoreFileHash(ctx, uc.apiKey, fileHash); err != nil {
		log.Printf("Error storing file hash in Redis: %v", err)
	}

	return taskID, nil
}

// generateTaskID generates a unique task ID.
func (uc *CSVProcessorUsecase) GenerateTaskID(fileName string) string {
	//TODO assume user id will be 1234
	return fileName + "-" + time.Now().Format("20060102150405")
}
