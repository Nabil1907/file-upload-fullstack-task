package config

import (
	"janan_csv_service/pkg/helpers"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	RedisAddr     string
	RedisPassword string
	APIKey        string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:        helpers.SafeGetEnv("DB_HOST"),
		DBPort:        helpers.SafeGetEnv("DB_PORT"),
		DBUser:        helpers.SafeGetEnv("DB_USER"),
		DBPassword:    helpers.SafeGetEnv("DB_PASSWORD"),
		DBName:        helpers.SafeGetEnv("DB_NAME"),
		RedisAddr:     helpers.SafeGetEnv("REDIS_ADDR"),
		RedisPassword: helpers.SafeGetEnv("REDIS_PASSWORD"),
		APIKey:        helpers.SafeGetEnv("API_KEY"),
	}
}
