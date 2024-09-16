package interfaces

import "transfigurr/models"

type EventRepositoryInterface interface {
	GetEvents() ([]models.Event, error)
	GetEventById(id string) (models.Event, error)
	UpsertEventById(event models.Event) error
	DeleteEventById(event models.Event) error
	Log(level, service, message string) error
}
