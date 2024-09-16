package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/interfaces"

	"github.com/gin-gonic/gin"
)

func MovieRoutes(rg *gin.RouterGroup, movieRepo interfaces.MovieRepositoryInterface) {
	controller := controllers.NewMovieController(movieRepo)
	rg.GET("/movies", controller.GetMovies)
	rg.GET("/movie/:movieId", controller.GetMovieByID)
	rg.PUT("/movie/:movieId", controller.UpsertMovie)
	rg.DELETE("/movie/:movieId", controller.DeleteMovieByID)
}
