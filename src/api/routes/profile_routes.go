package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/repository"

	"github.com/gin-gonic/gin"
)

func ProfileRoutes(rg *gin.RouterGroup, profileRepo *repository.ProfileRepository) {
	controller := controllers.NewProfileController(profileRepo)
	rg.GET("/profiles", controller.GetProfiles)
	rg.GET("/profiles/:profileId", controller.GetProfileById)
	rg.PUT("/profiles/:profileId", controller.UpsertProfile)
	rg.DELETE("/profiles/:profileId", controller.DeleteProfileById)
}
