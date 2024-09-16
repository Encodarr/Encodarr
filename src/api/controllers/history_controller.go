package controllers

import (
	"errors"
	"log"
	"net/http"
	"transfigurr/interfaces"
	"transfigurr/models"
	"transfigurr/repository"

	"github.com/gin-gonic/gin"
)

type HistoryController struct {
	Repo interfaces.HistoryRepositoryInterface
}

func NewHistoryController(repo interfaces.HistoryRepositoryInterface) *HistoryController {
	return &HistoryController{
		Repo: repo,
	}
}

func (ctrl HistoryController) GetHistories(c *gin.Context) {
	histories, err := ctrl.Repo.GetHistories()
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Histories not found"})
		} else {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving histories"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, histories)
}

func (ctrl HistoryController) UpsertHistory(c *gin.Context) {
	var inputHistory models.History
	historyId := c.Param("historyId")

	if err := c.ShouldBindJSON(&inputHistory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	history, err := ctrl.Repo.GetHistoryById(historyId)
	if err != nil && errors.Is(err, repository.ErrRecordNotFound) {
		history = inputHistory
	} else if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving history"})
		return
	}

	if err := ctrl.Repo.UpsertHistoryById(history); err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error upserting history"})
		return
	}

	c.IndentedJSON(http.StatusOK, history)
}

func (ctrl HistoryController) GetHistoryById(c *gin.Context) {
	historyId := c.Param("historyId")
	history, err := ctrl.Repo.GetHistoryById(historyId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "History not found"})
		} else {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving history"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, history)
}

func (ctrl HistoryController) DeleteHistoryById(c *gin.Context) {
	historyId := c.Param("historyId")
	history, err := ctrl.Repo.GetHistoryById(historyId)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "History not found"})
		return
	}

	if err := ctrl.Repo.DeleteHistoryById(history); err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "History deleted successfully"})
}
