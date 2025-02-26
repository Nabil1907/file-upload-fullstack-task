package repository

import (
	"janan_csv_service/internal/models"
	"log"

	"gorm.io/gorm"
)

type CSVProcessorRepository struct {
	db *gorm.DB
}

func NewCSVProcessorRepository(db *gorm.DB) *CSVProcessorRepository {
	return &CSVProcessorRepository{db: db}
}
func (r *CSVProcessorRepository) SaveStudent(student models.Student) error {
	// Save the student record using GORM
	result := r.db.Create(&student)
	if result.Error != nil {
		log.Printf("Error saving student: %v", result.Error)
		return result.Error
	}
	return nil
}
