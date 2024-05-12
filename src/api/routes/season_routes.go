package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/repository"

	"github.com/gin-gonic/gin"
)

func SeasonRoutes(rg *gin.RouterGroup, seasonRepo *repository.SeasonRepository) {
	controller := controllers.NewSeasonController(seasonRepo)
	rg.GET("/series/:seriesId/seasons", controller.GetSeasons)
	rg.GET("/series/:seriesId/seasons/:seasonNumber", controller.GetSeasonById)
	rg.PUT("/series/:seriesId/seasons/:seasonNumber", controller.UpsertSeason)
	rg.DELETE("/series/:seriesId/seasons/:seasonNumber", controller.DeleteSeasonById)
}
