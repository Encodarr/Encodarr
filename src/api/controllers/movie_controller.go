package controllers

import (
	"errors"
	"net/http"
	"transfigurr/interfaces"
	"transfigurr/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MovieController struct {
	Repo interfaces.MovieRepositoryInterface
}

func NewMovieController(repo interfaces.MovieRepositoryInterface) *MovieController {
	return &MovieController{
		Repo: repo,
	}
}

func (ctrl *MovieController) GetMovies(c *gin.Context) {
	movieList, err := ctrl.Repo.GetMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving movie"})
		return
	}

	c.IndentedJSON(http.StatusOK, movieList)
}

func (ctrl *MovieController) UpsertMovie(c *gin.Context) {
	var inputMovie models.Movie
	id := c.Param("movieId")

	if err := c.ShouldBindJSON(&inputMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	movie, err := ctrl.Repo.UpsertMovie(id, inputMovie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error upserting movie"})
		return
	}

	c.IndentedJSON(http.StatusOK, movie)
}

func (ctrl *MovieController) GetMovieByID(c *gin.Context) {
	id := c.Param("movieId")

	movie, err := ctrl.Repo.GetMovieById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving movie"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, movie)
}

func (ctrl *MovieController) DeleteMovieByID(c *gin.Context) {
	id := c.Param("movieId")

	err := ctrl.Repo.DeleteMovieById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
}
