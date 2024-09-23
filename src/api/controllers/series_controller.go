package controllers

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"transfigurr/interfaces"
	"transfigurr/models"

	"github.com/gin-gonic/gin"
)

type SeriesController struct {
	Repo interfaces.SeriesRepositoryInterface
}

func NewSeriesController(repo interfaces.SeriesRepositoryInterface) *SeriesController {
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

	// Log the incoming request body for debugging
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Print("Error reading request body: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading request body"})
		return
	}
	log.Printf("Request body: %s", body)

	// Reset the request body so it can be read again by ShouldBindJSON
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if err := c.ShouldBindJSON(&inputSeries); err != nil {
		log.Print("Error binding JSON: ", err)
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
