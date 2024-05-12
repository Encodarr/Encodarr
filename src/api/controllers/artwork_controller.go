package controllers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"transfigurr/constants"

	"github.com/gin-gonic/gin"
)

type ArtworkController struct{}

func NewArtworkController() *ArtworkController {
	return &ArtworkController{}
}

func (a *ArtworkController) GetSeriesBackdrop(c *gin.Context) {
	seriesId := c.Param("seriesId")
	filePath := filepath.Join(constants.ArtworkPath, "series", seriesId, "backdrop.webp")
	log.Print(filePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Backdrop not found"})
		return
	}
	c.File(filePath)
}

func (a *ArtworkController) GetSeriesPoster(c *gin.Context) {
	seriesId := c.Param("seriesId")
	filePath := filepath.Join(constants.ArtworkPath, "series", seriesId, "poster.webp")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Poster not found"})
		return
	}
	c.File(filePath)
}

func (a *ArtworkController) GetMovieBackdrop(c *gin.Context) {
	movieId := c.Param("movieId")
	filePath := filepath.Join(constants.ArtworkPath, "movies", movieId, "backdrop.webp")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Backdrop not found"})
		return
	}
	c.File(filePath)
}

func (a *ArtworkController) GetMoviePoster(c *gin.Context) {
	movieId := c.Param("movieId")
	filePath := filepath.Join(constants.ArtworkPath, "movies", movieId, "poster.webp")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Poster not found"})
		return
	}
	c.File(filePath)
}
