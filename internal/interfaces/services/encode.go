package services

import "transfigurr/internal/models"

type EncodeServiceInterface interface {
	Enqueue(item models.Item)
	Startup()
	GetQueue() models.QueueStatus
}
