package repository

import (
	"log"
	"time"
	"transfigurr/models"

	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func (repo *EventRepository) DeleteEventById(event models.Event) error {
	panic("unimplemented")
}

func (repo *EventRepository) GetEventById(id string) (models.Event, error) {
	panic("unimplemented")
}

func (repo *EventRepository) GetEvents() ([]models.Event, error) {
	var events []models.Event
	result := repo.db.Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

func (repo *EventRepository) UpsertEventById(event models.Event) error {
	panic("unimplemented")
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (repo *EventRepository) Log(level, service, message string) error {
	log.Print("Attempting to log entry")
	log.Printf("Level: %s, Service: %s, Message: %s", level, service, message)
	eventEntry := models.Event{
		Timestamp: time.Now().Format("2006-01-02T15:04:05.000"),
		Level:     level,
		Service:   service,
		Message:   message,
	}

	if err := repo.db.Create(&eventEntry).Error; err != nil {
		log.Printf("Error creating event entry: %v", err)
		return err
	}

	log.Print("Event entry successfully created")
	return nil
}
