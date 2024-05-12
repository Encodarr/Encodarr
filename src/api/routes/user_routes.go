package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/repository"

	"github.com/gin-gonic/gin"
)

func UserRoutes(rg *gin.RouterGroup, userRepo *repository.UserRepository) {
	controller := controllers.NewUserController(userRepo)
	rg.GET("/user", controller.GetUsers)
	rg.POST("/user", controller.UpdateUser)
}
