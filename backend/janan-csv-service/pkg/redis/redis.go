package redis

import (
	"context"
	"encoding/json"
	"janan_csv_service/config"
	"log"
	"time"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg *config.Config) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		Password:     cfg.RedisPassword,
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 2,
		MaxConnAge:   5 * time.Minute,
	})

	if err := client.Ping().Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &RedisClient{Client: client}
}

// StoreFileHash stores a file hash under the API key in Redis.
func (r *RedisClient) StoreFileHash(ctx context.Context, apiKey, fileHash string) error {
	// Get the existing array of file hashes
	hashes, err := r.GetFileHashes(ctx, apiKey)
	if err != nil && err != redis.Nil {
		return err
	}

	hashes = append(hashes, fileHash)
	hashesJSON, err := json.Marshal(hashes)
	if err != nil {
		return err
	}

	return r.Client.Set(apiKey, hashesJSON, 0).Err()
}

// GetFileHashes retrieves the array of file hashes for the given API key.
func (r *RedisClient) GetFileHashes(ctx context.Context, apiKey string) ([]string, error) {
	hashesJSON, err := r.Client.Get(apiKey).Result()
	if err == redis.Nil {
		return []string{}, nil // Return empty array if key doesn't exist
	}
	if err != nil {
		return nil, err
	}

	var hashes []string
	if err := json.Unmarshal([]byte(hashesJSON), &hashes); err != nil {
		return nil, err
	}

	return hashes, nil
}

// StoreProgress stores the progress of a task in Redis.
func (r *RedisClient) StoreProgress(ctx context.Context, taskID string, progress map[string]interface{}) error {
	progressJSON, err := json.Marshal(progress)
	if err != nil {
		return err
	}
	return r.Client.Set(taskID, progressJSON, 0).Err()
}

// GetProgress retrieves the progress of a task from Redis.
func (r *RedisClient) GetProgress(ctx context.Context, taskID string) (map[string]interface{}, error) {
	progressJSON, err := r.Client.Get(taskID).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var progress map[string]interface{}
	if err := json.Unmarshal([]byte(progressJSON), &progress); err != nil {
		return nil, err
	}

	return progress, nil
}
