package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/repository"

	"github.com/gin-gonic/gin"
)

func EpisodeRoutes(rg *gin.RouterGroup, episodeRepo *repository.EpisodeRepository) {
	controller := controllers.NewEpisodeController(episodeRepo)
	rg.GET("/series/:seriesId/seasons/:seasonNumber/episodes", controller.GetEpisodes)
	rg.GET("/series/:seriesId/seasons/:seasonNumber/episodes/:episodeNumber", controller.GetEpisodeById)
	rg.PUT("/series/:seriesId/seasons/:seasonNumber/episodes/:episodeNumber", controller.UpsertEpisode)
	rg.DELETE("/series/:seriesId/seasons/:seasonNumber/episodes/:episodeNumber", controller.DeleteEpisodeById)
}
