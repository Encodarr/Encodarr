package controllers

import (
	"log"
	"net/http"
	"strconv"
	"transfigurr/models"
	"transfigurr/repository"

	"github.com/gin-gonic/gin"
)

type EpisodeController struct {
	Repo *repository.EpisodeRepository
}

func NewEpisodeController(repo *repository.EpisodeRepository) *EpisodeController {
	return &EpisodeController{
		Repo: repo,
	}
}

func (ctrl *EpisodeController) GetEpisodes(c *gin.Context) {
	seriesId := c.Param("seriesId")
	seasonNumber := c.Param("seasonNumber")

	seasonNum, err := strconv.Atoi(seasonNumber)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season number"})
		return
	}

	episodes, err := ctrl.Repo.GetEpisodes(seriesId, seasonNum)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving episodes"})
		return
	}

	c.IndentedJSON(http.StatusOK, episodes)
}

func (ctrl *EpisodeController) UpsertEpisode(c *gin.Context) {
	var inputEpisode models.Episode
	seriesId := c.Param("seriesId")
	seasonNumber := c.Param("seasonNumber")
	episodeNumber := c.Param("episodeNumber")

	seasonNum, err := strconv.Atoi(seasonNumber)

	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season number"})
		return
	}

	episodeNum, err := strconv.Atoi(episodeNumber)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode number"})
		return
	}

	if err := c.ShouldBindJSON(&inputEpisode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	episode, err := ctrl.Repo.UpsertEpisode(seriesId, seasonNum, episodeNum, inputEpisode)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error upserting episode"})
		return
	}

	c.IndentedJSON(http.StatusOK, episode)
}

func (ctrl *EpisodeController) GetEpisodeBySeriesSeasonEpisode(c *gin.Context) {
	seriesId := c.Param("seriesId")
	seasonNumber := c.Param("seasonNumber")
	episodeNumber := c.Param("episodeNumber")

	seasonNum, err := strconv.Atoi(seasonNumber)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season number"})
		return
	}

	episodeNum, err := strconv.Atoi(episodeNumber)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode number"})
		return
	}

	episode, err := ctrl.Repo.GetEpisodeBySeriesSeasonEpisode(seriesId, seasonNum, episodeNum)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving episode"})
		return
	}

	c.IndentedJSON(http.StatusOK, episode)
}

func (ctrl *EpisodeController) DeleteEpisodeById(c *gin.Context) {
	seriesId := c.Param("seriesId")
	seasonNumber := c.Param("seasonNumber")
	episodeNumber := c.Param("episodeNumber")

	seasonNum, err := strconv.Atoi(seasonNumber)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season number"})
		return
	}

	episodeNum, err := strconv.Atoi(episodeNumber)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode number"})
		return
	}

	err = ctrl.Repo.DeleteEpisodeById(seriesId, seasonNum, episodeNum)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting episode"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Episode deleted successfully"})
}
