package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/interfaces"

	"github.com/gin-gonic/gin"
)

func ActionRoutes(rg *gin.RouterGroup, scanService interfaces.ScanServiceInterface, metadataService interfaces.MetadataServiceInterface) {
	controller := controllers.NewActionController(scanService, metadataService)
	rg.POST("/actions/restart", controller.Restart)
	rg.POST("/actions/shutdown", controller.Shutdown)
	rg.POST("/actions/refresh/metadata", controller.RefreshMetadata)
	rg.POST("/actions/scan", controller.Scan)
	rg.POST("/actions/refresh/metadata/series/:series_id", controller.RefreshSeriesMetadata)
	rg.POST("/actions/scan/series/:series_id", controller.ScanSeries)
	rg.POST("/actions/refresh/metadata/movies/:movie_id", controller.RefreshMovieMetadata)
	rg.POST("/actions/scan/movies/:movie_id", controller.ScanMovie)
}
