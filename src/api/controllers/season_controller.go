package controllers

import (
	"net/http"
	"strconv"
	"transfigurr/interfaces"
	"transfigurr/models"

	"github.com/gin-gonic/gin"
)

type SeasonController struct {
	Repo interfaces.SeasonRepositoryInterface
}

func NewSeasonController(repo interfaces.SeasonRepositoryInterface) *SeasonController {
	return &SeasonController{
		Repo: repo,
	}
}

func (ctrl *SeasonController) GetSeasons(c *gin.Context) {
	seriesId := c.Param("seriesId")

	seasons, err := ctrl.Repo.GetSeasons(seriesId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving seasons"})
		return
	}

	c.IndentedJSON(http.StatusOK, seasons)
}

func (ctrl *SeasonController) UpsertSeason(c *gin.Context) {
	var inputSeason models.Season
	seriesId := c.Param("seriesId")
	seasonNumber := c.Param("seasonNumber")

	seasonNum, err := strconv.Atoi(seasonNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season number"})
		return
	}

	if err := c.ShouldBindJSON(&inputSeason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	season, err := ctrl.Repo.UpsertSeason(seriesId, seasonNum, inputSeason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error upserting season"})
		return
	}

	c.IndentedJSON(http.StatusOK, season)
}

func (ctrl *SeasonController) GetSeasonById(c *gin.Context) {
	seriesId := c.Param("seriesId")
	seasonNumber := c.Param("seasonNumber")

	seasonNum, err := strconv.Atoi(seasonNumber)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season number"})
		return
	}

	season, err := ctrl.Repo.GetSeasonById(seriesId, seasonNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving season"})
		return
	}

	c.IndentedJSON(http.StatusOK, season)
}

func (ctrl *SeasonController) DeleteSeasonById(c *gin.Context) {
	seriesId := c.Param("seriesId")
	seasonNumber := c.Param("seasonNumber")

	seasonNum, err := strconv.Atoi(seasonNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season number"})
		return
	}

	err = ctrl.Repo.DeleteSeasonById(seriesId, seasonNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting season"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Season deleted successfully"})
}
