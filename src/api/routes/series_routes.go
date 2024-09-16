package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/interfaces"

	"github.com/gin-gonic/gin"
)

func SeriesRoutes(rg *gin.RouterGroup, seriesRepo interfaces.SeriesRepositoryInterface) {
	controller := controllers.NewSeriesController(seriesRepo)
	rg.GET("/series", controller.GetSeries)
	rg.GET("/series/:seriesId", controller.GetSeriesByID)
	rg.PUT("/series/:seriesId", controller.UpsertSeries)
	rg.DELETE("/series/:seriesId", controller.DeleteSeriesByID)
}
