package routes

import (
	"transfigurr/api/controllers"

	"github.com/gin-gonic/gin"
)

func ArtworkRoutes(rg *gin.RouterGroup) {
	controller := controllers.NewArtworkController()
	rg.GET("/series/:seriesId/backdrop", controller.GetSeriesBackdrop)
	rg.GET("/series/:seriesId/poster", controller.GetSeriesPoster)
	rg.GET("/movies/:movieId/backdrop", controller.GetMovieBackdrop)
	rg.GET("/movies/:movieId/poster", controller.GetMoviePoster)
}
