package controllers

import (
	"errors"
	"net/http"
	"transfigurr/interfaces"
	"transfigurr/models"
	"transfigurr/repository"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	Repo interfaces.EventRepositoryInterface
}

func NewEventController(repo interfaces.EventRepositoryInterface) *EventController {
	return &EventController{
		Repo: repo,
	}
}

func (ctrl EventController) GetEvents(c *gin.Context) {
	events, err := ctrl.Repo.GetEvents()
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Events not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving events"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, events)
}

func (ctrl EventController) UpsertEvent(c *gin.Context) {
	var inputEvent models.Event
	eventId := c.Param("eventId")

	if err := c.ShouldBindJSON(&inputEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	event, err := ctrl.Repo.GetEventById(eventId)
	if err != nil && errors.Is(err, repository.ErrRecordNotFound) {
		event = inputEvent
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving event"})
		return
	}

	if err := ctrl.Repo.UpsertEventById(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error upserting event"})
		return
	}

	c.IndentedJSON(http.StatusOK, event)
}

func (ctrl EventController) GetEventById(c *gin.Context) {
	eventId := c.Param("eventId")
	event, err := ctrl.Repo.GetEventById(eventId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving event"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, event)
}

func (ctrl EventController) DeleteEventById(c *gin.Context) {
	eventId := c.Param("eventId")
	event, err := ctrl.Repo.GetEventById(eventId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err := ctrl.Repo.DeleteEventById(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}
