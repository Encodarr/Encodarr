package controllers

import (
	"log"
	"net/http"
	"transfigurr/models"
	"transfigurr/repository"

	"github.com/gin-gonic/gin"
)

type SeriesController struct {
	Repo *repository.SeriesRepository
}

func NewSeriesController(repo *repository.SeriesRepository) *SeriesController {
	return &SeriesController{
		Repo: repo,
	}
}

func (ctrl *SeriesController) GetSeries(c *gin.Context) {
	seriesList, err := ctrl.Repo.GetSeries()
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving series"})
		return
	}

	c.IndentedJSON(http.StatusOK, seriesList)
}

func (ctrl *SeriesController) UpsertSeries(c *gin.Context) {
	var inputSeries models.Series
	id := c.Param("seriesId")

	if err := c.ShouldBindJSON(&inputSeries); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	series, err := ctrl.Repo.UpsertSeries(id, inputSeries)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error upserting series"})
		return
	}

	c.IndentedJSON(http.StatusOK, series)
}

func (ctrl *SeriesController) GetSeriesByID(c *gin.Context) {
	id := c.Param("seriesId")

	series, err := ctrl.Repo.GetSeriesByID(id)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving series"})
		return
	}

	c.IndentedJSON(http.StatusOK, series)
}

func (ctrl *SeriesController) DeleteSeriesByID(c *gin.Context) {
	id := c.Param("seriesId")

	err := ctrl.Repo.DeleteSeriesByID(id)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting series"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Series deleted successfully"})
}
