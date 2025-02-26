package postgres

import (
	"janan_csv_service/config"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := "host=" + cfg.DBHost + " user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " port=" + cfg.DBPort + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)

	sqlDB.SetMaxOpenConns(100)

	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
