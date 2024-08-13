package controllers

import (
	"os"
	"transfigurr/interfaces"
	"transfigurr/models"

	"github.com/gin-gonic/gin"
)

type ActionController struct {
	scanService interfaces.ScanServiceInterface
}

func NewActionController(scanService interfaces.ScanServiceInterface) *ActionController {
	return &ActionController{
		scanService: scanService,
	}
}

func (ctrl ActionController) Restart(c *gin.Context) {
}

func (ctrl ActionController) Shutdown(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Application is shutting down"})
	os.Exit(0)
}

func (ctrl ActionController) RefreshMetadata(c *gin.Context) {
}

func (ctrl ActionController) Scan(c *gin.Context) {

}

func (ctrl ActionController) RefreshSeriesMetadata(c *gin.Context) {

}

func (ctrl ActionController) ScanSeries(c *gin.Context) {
	ctrl.scanService.Enqueue(models.Item{Id: c.Param("series_id"), Type: "series"})
	c.JSON(200, gin.H{"message": "Scan enqueued"})

}

func (ctrl ActionController) RefreshMovieMetadata(c *gin.Context) {

}

func (ctrl ActionController) ScanMovie(c *gin.Context) {
	ctrl.scanService.Enqueue(models.Item{Id: c.Param("movie_id"), Type: "movie"})
	c.JSON(200, gin.H{"message": "Scan enqueued"})
}
