package repository

import (
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
	eventEntry := models.Event{
		Timestamp: time.Now().Format("2006-01-02T15:04:05.000"),
		Level:     level,
		Service:   service,
		Message:   message,
	}

	if err := repo.db.Create(&eventEntry).Error; err != nil {
		return err
	}

	return nil
}
