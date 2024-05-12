package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/interfaces"

	"github.com/gin-gonic/gin"
)

func ActionRoutes(rg *gin.RouterGroup, scanService interfaces.ScanServiceInterface) {
	controller := controllers.NewActionController(scanService)
	rg.POST("/actions/restart", controller.Restart)
	rg.POST("/actions/shutdown", controller.Shutdown)
	rg.POST("/actions/refresh/metadata", controller.RefreshMetadata)
	rg.POST("/actions/scan", controller.Scan)
	rg.POST("/actions/refresh/metadata/series/:series_id", controller.RefreshSeriesMetadata)
	rg.POST("/actions/scan/series/:series_id", controller.ScanSeries)
	rg.POST("/actions/refresh/metadata/movie/:movie_id", controller.RefreshMovieMetadata)
	rg.POST("/actions/scan/movie/:movie_id", controller.ScanMovie)
}
