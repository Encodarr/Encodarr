package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/interfaces"

	"github.com/gin-gonic/gin"
)

func EpisodeRoutes(rg *gin.RouterGroup, episodeRepo interfaces.EpisodeRepositoryInterface) {
	controller := controllers.NewEpisodeController(episodeRepo)
	rg.GET("/series/:seriesId/seasons/:seasonNumber/episodes", controller.GetEpisodes)
	rg.GET("/series/:seriesId/seasons/:seasonNumber/episodes/:episodeNumber", controller.GetEpisodeBySeriesSeasonEpisode)
	rg.PUT("/series/:seriesId/seasons/:seasonNumber/episodes/:episodeNumber", controller.UpsertEpisode)
	rg.DELETE("/series/:seriesId/seasons/:seasonNumber/episodes/:episodeNumber", controller.DeleteEpisodeById)
}
