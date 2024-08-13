package interfaces

import "transfigurr/models"

type MetadataServiceInterface interface {
	Startup()
	Enqueue(item models.Item)
}
