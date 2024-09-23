package controllers

import (
	"errors"
	"log"
	"net/http"
	"transfigurr/interfaces"
	"transfigurr/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SystemController struct {
	Repo interfaces.SystemRepositoryInterface
}

func NewSystemController(repo interfaces.SystemRepositoryInterface) *SystemController {
	return &SystemController{
		Repo: repo,
	}
}

func (ctrl *SystemController) GetSystems(c *gin.Context) {
	systems, err := ctrl.Repo.GetSystems()
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving systems"})
		return
	}

	c.IndentedJSON(http.StatusOK, systems)
}

func (ctrl *SystemController) UpsertSystem(c *gin.Context) {
	var inputSystem models.System
	systemId := c.Param("systemId")

	if err := c.ShouldBindJSON(&inputSystem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	system, err := ctrl.Repo.UpsertSystem(systemId, inputSystem)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error upserting system"})
		return
	}

	c.IndentedJSON(http.StatusOK, system)
}

func (ctrl *SystemController) GetSystemById(c *gin.Context) {
	systemId := c.Param("systemId")

	system, err := ctrl.Repo.GetSystemById(systemId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "System not found"})
		} else {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving system"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, system)
}

func (ctrl *SystemController) DeleteSystemById(c *gin.Context) {
	systemId := c.Param("systemId")

	err := ctrl.Repo.DeleteSystemById(systemId)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting system"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "System deleted successfully"})
}
