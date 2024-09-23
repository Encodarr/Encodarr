package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/interfaces"

	"github.com/gin-gonic/gin"
)

func MovieRoutes(rg *gin.RouterGroup, movieRepo interfaces.MovieRepositoryInterface) {
	controller := controllers.NewMovieController(movieRepo)
	rg.GET("/movies", controller.GetMovies)
	rg.GET("/movies/:movieId", controller.GetMovieByID)
	rg.PUT("/movies/:movieId", controller.UpsertMovie)
	rg.DELETE("/movies/:movieId", controller.DeleteMovieByID)
}
