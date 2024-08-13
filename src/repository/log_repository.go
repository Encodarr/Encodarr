package repository

import (
	"log"
	"time"
	"transfigurr/models"

	"github.com/jinzhu/gorm"
)

type LoggingRepository struct {
	db *gorm.DB
}

func NewLoggingRepository(db *gorm.DB) *LoggingRepository {
	return &LoggingRepository{db: db}
}

func (repo *LoggingRepository) Log(level, service, message string) error {
	log.Print("Attempting to log entry")
	log.Printf("Level: %s, Service: %s, Message: %s", level, service, message)
	logEntry := models.Log{
		Timestamp: time.Now().Format("2006-01-02T15:04:05.000"),
		Level:     level,
		Service:   service,
		Message:   message,
	}

	if err := repo.db.Create(&logEntry).Error; err != nil {
		log.Printf("Error creating log entry: %v", err)
		return err
	}

	log.Print("Log entry successfully created")
	return nil
}
