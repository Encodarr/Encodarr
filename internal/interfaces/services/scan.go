package services

import "transfigurr/internal/models"

type ScanServiceInterface interface {
	Startup()
	Enqueue(item models.Item)
	EnqueueAll()
	EnqueueAllSeries()
	EnqueueAllMovies()
}
