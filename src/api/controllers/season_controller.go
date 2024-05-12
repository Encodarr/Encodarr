package controllers

import (
	"log"
	"net/http"
	"transfigurr/models"
	"transfigurr/repository"

	"github.com/gin-gonic/gin"
)

type SeasonController struct {
	Repo *repository.SeasonRepository
}

func NewSeasonController(repo *repository.SeasonRepository) *SeasonController {
	return &SeasonController{
		Repo: repo,
	}
}

func (ctrl *SeasonController) GetSeasons(c *gin.Context) {
	seriesId := c.Param("seriesId")

	seasons, err := ctrl.Repo.GetSeasons(seriesId)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving seasons"})
		return
	}

	c.IndentedJSON(http.StatusOK, seasons)
}

func (ctrl *SeasonController) UpsertSeason(c *gin.Context) {
	var inputSeason models.Season
	seriesId := c.Param("seriesId")
	seasonNumber := c.Param("seasonNumber")

	if err := c.ShouldBindJSON(&inputSeason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	season, err := ctrl.Repo.UpsertSeason(seriesId, seasonNumber, inputSeason)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error upserting season"})
		return
	}

	c.IndentedJSON(http.StatusOK, season)
}

func (ctrl *SeasonController) GetSeasonById(c *gin.Context) {
	seriesId := c.Param("seriesId")
	seasonNumber := c.Param("seasonNumber")

	season, err := ctrl.Repo.GetSeasonById(seriesId, seasonNumber)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving season"})
		return
	}

	c.IndentedJSON(http.StatusOK, season)
}

func (ctrl *SeasonController) DeleteSeasonById(c *gin.Context) {
	seriesId := c.Param("seriesId")
	seasonNumber := c.Param("seasonNumber")

	err := ctrl.Repo.DeleteSeasonById(seriesId, seasonNumber)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting season"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Season deleted successfully"})
}
