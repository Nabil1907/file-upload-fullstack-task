package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Student struct {
	UUID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"uuid"`
	StudentID   string    `gorm:"not null" json:"student_id"`
	StudentName string    `gorm:"not null" json:"student_name"`
	Subject     string    `gorm:"not null" json:"subject"`
	Grade       float64   `gorm:"not null" json:"grade"`
}

// BeforeCreate is a GORM hook that generates a UUID for the student before creating the record.
func (s *Student) BeforeCreate(tx *gorm.DB) (err error) {
	s.UUID = uuid.NewV4()
	return
}
