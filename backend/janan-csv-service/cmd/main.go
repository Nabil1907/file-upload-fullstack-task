// cmd/main.go
package main

import (
	"janan_csv_service/config"
	"janan_csv_service/internal/csvprocessor/delivery/http"
	"janan_csv_service/internal/csvprocessor/repository"
	"janan_csv_service/internal/csvprocessor/usecase"
	"janan_csv_service/pkg/postgres"
	"janan_csv_service/pkg/progress"
	"janan_csv_service/pkg/redis"
	"log"
	"net"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Test DNS resolution
	_, err := net.LookupHost("localhost")
	if err != nil {
		log.Fatalf("DNS resolution failed: %v", err)
	}
	log.Println("DNS resolution succeeded")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.LoadConfig()

	// Initialize the progress tracker
	progressTracker := progress.NewTracker()

	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		panic(err)
	}

	csvRepo := repository.NewCSVProcessorRepository(db)
	redisClient := redis.NewRedisClient(cfg)
	csvUsecase := usecase.NewCSVProcessorUsecase(csvRepo, redisClient, cfg.APIKey, progressTracker)
	csvHandler := http.NewCSVProcessorHandler(csvUsecase, redisClient, progressTracker)

	r := gin.Default()

	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow requests from the frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Group routes under /api/v1
	api := r.Group("/api/v1")
	{
		api.POST("/upload", csvHandler.UploadCSV)
		api.GET("/progress/:task_id", csvHandler.CheckProgress)
		api.GET("/progress-stream/:task_id", csvHandler.StreamProgress)
	}

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
