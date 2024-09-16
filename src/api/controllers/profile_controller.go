package controllers

import (
	"log"
	"net/http"
	"strconv"
	"transfigurr/interfaces"
	"transfigurr/models"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	Repo interfaces.ProfileRepositoryInterface
}

func NewProfileController(repo interfaces.ProfileRepositoryInterface) *ProfileController {
	return &ProfileController{
		Repo: repo,
	}
}

func (ctrl *ProfileController) GetProfiles(c *gin.Context) {
	profiles, err := ctrl.Repo.GetAllProfiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving profiles"})
		return
	}

	c.IndentedJSON(http.StatusOK, profiles)
}

func (ctrl *ProfileController) GetProfileById(c *gin.Context) {
	profileId := c.Param("profileId")
	profileIdInt, err := strconv.Atoi(profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID"})
		return
	}

	profile, err := ctrl.Repo.GetProfileById(profileIdInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, profile)
}

func (ctrl *ProfileController) UpsertProfile(c *gin.Context) {
	var inputProfile models.Profile
	profileId := c.Param("profileId")

	profileIdInt, err := strconv.Atoi(profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID"})
		return
	}

	if err := c.ShouldBindJSON(&inputProfile); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	profile, err := ctrl.Repo.UpsertProfile(profileIdInt, inputProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating profile"})
		return
	}

	c.IndentedJSON(http.StatusOK, profile)
}

func (ctrl *ProfileController) DeleteProfileById(c *gin.Context) {
	profileId := c.Param("profileId")

	profileIdInt, err := strconv.Atoi(profileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID"})
		return
	}

	err = ctrl.Repo.DeleteProfileById(profileIdInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
}
