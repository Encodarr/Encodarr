package services

import "transfigurr/internal/models"

type MetadataServiceInterface interface {
	Startup()
	Enqueue(item models.Item)
	EnqueueAll()
}
