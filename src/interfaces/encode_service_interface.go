package interfaces

import "transfigurr/models"

type EncodeServiceInterface interface {
	Enqueue(item models.Item)
	Startup()
	GetQueue() models.QueueStatus
}
