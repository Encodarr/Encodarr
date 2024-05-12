// setting_controller.go
package controllers

import (
	"errors"
	"log"
	"net/http"
	"transfigurr/models"
	"transfigurr/repository"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type SettingController struct {
	Repo *repository.SettingRepository
}

func NewSettingController(repo *repository.SettingRepository) *SettingController {
	return &SettingController{
		Repo: repo,
	}
}

func (ctrl *SettingController) GetSettings(c *gin.Context) {
	settings, err := ctrl.Repo.GetAllSettings()
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving settings"})
		return
	}

	c.IndentedJSON(http.StatusOK, settings)
}

func (ctrl *SettingController) UpsertSetting(c *gin.Context) {
	var inputSetting models.Setting
	settingId := c.Param("settingId")

	if err := c.ShouldBindJSON(&inputSetting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	setting, err := ctrl.Repo.GetSettingById(settingId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create a new setting if no setting exists with the provided ID
		err = ctrl.Repo.CreateSetting(inputSetting)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating setting"})
			return
		}
		c.IndentedJSON(http.StatusCreated, inputSetting)
	} else if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving setting"})
		return
	} else {
		ctrl.Repo.UpdateSetting(inputSetting)
		err = ctrl.Repo.UpdateSetting(setting)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating setting"})
			return
		}
		c.IndentedJSON(http.StatusOK, setting)
	}
}

func (ctrl *SettingController) GetSettingById(c *gin.Context) {
	settingId := c.Param("settingId")

	setting, err := ctrl.Repo.GetSettingById(settingId)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving setting"})
		return
	}

	c.IndentedJSON(http.StatusOK, setting)
}

func (ctrl *SettingController) DeleteSettingById(c *gin.Context) {
	settingId := c.Param("settingId")

	setting, err := ctrl.Repo.GetSettingById(settingId)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}

	err = ctrl.Repo.DeleteSetting(setting)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting setting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Setting deleted successfully"})
}
