package repository

import (
	"errors"
	"transfigurr/models"

	"github.com/jinzhu/gorm"
)

type EventRepository struct {
	DB *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{
		DB: db,
	}
}

func (repo *EventRepository) GetEvents() ([]models.Event, error) {
	var events []models.Event
	if err := repo.DB.Find(&events).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return events, nil
}

func (repo *EventRepository) GetEventById(id string) (models.Event, error) {
	var event models.Event
	if err := repo.DB.Where("id = ?", id).First(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Event{}, ErrRecordNotFound
		}
		return event, err
	}
	return event, nil
}

func (repo *EventRepository) UpsertEventById(event models.Event) error {
	if err := repo.DB.Save(event).Error; err != nil {
		return err
	}
	return nil
}

func (repo *EventRepository) DeleteEventById(event models.Event) error {
	if err := repo.DB.Delete(event).Error; err != nil {
		return err
	}
	return nil
}
