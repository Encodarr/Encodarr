package interfaces

import "transfigurr/models"

type ScanServiceInterface interface {
	Startup()
	Enqueue(item models.Item)
	EnqueueAll()
	EnqueueAllSeries()
	EnqueueAllMovies()
}
