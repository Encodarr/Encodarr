package controllers

import (
	"os"
	"transfigurr/interfaces"
	"transfigurr/models"

	"github.com/gin-gonic/gin"
)

type ActionController struct {
	scanService     interfaces.ScanServiceInterface
	metadataService interfaces.MetadataServiceInterface
}

func NewActionController(scanService interfaces.ScanServiceInterface, metadataService interfaces.MetadataServiceInterface) *ActionController {
	return &ActionController{
		scanService:     scanService,
		metadataService: metadataService,
	}
}

func (ctrl ActionController) Restart(c *gin.Context) {
}

func (ctrl ActionController) Shutdown(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Application is shutting down"})
	os.Exit(0)
}

func (ctrl ActionController) RefreshMetadata(c *gin.Context) {
	ctrl.metadataService.EnqueueAll()
	c.JSON(200, gin.H{"message": "Metadata refresh enqueued"})

}

func (ctrl ActionController) Scan(c *gin.Context) {
	ctrl.scanService.EnqueueAll()
	c.JSON(200, gin.H{"message": "Scan enqueued"})
}

func (ctrl ActionController) RefreshSeriesMetadata(c *gin.Context) {
	ctrl.metadataService.Enqueue(models.Item{Id: c.Param("series_id"), Type: "series"})
	c.JSON(200, gin.H{"message": "Refresh enqueued"})

}

func (ctrl ActionController) ScanSeries(c *gin.Context) {
	ctrl.scanService.Enqueue(models.Item{Id: c.Param("series_id"), Type: "series"})
	c.JSON(200, gin.H{"message": "Scan enqueued"})

}

func (ctrl ActionController) RefreshMovieMetadata(c *gin.Context) {
	ctrl.metadataService.Enqueue(models.Item{Id: c.Param("movie_id"), Type: "movie"})
	c.JSON(200, gin.H{"message": "Refresh enqueued"})
}

func (ctrl ActionController) ScanMovie(c *gin.Context) {
	ctrl.scanService.Enqueue(models.Item{Id: c.Param("movie_id"), Type: "movie"})
	c.JSON(200, gin.H{"message": "Scan enqueued"})
}
