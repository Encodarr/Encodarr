package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/repository"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup, authRepo *repository.AuthRepository) {
	controller := controllers.NewAuthController(authRepo)
	rg.GET("/activated", controller.GetActivated)
	rg.POST("/register", controller.Register)
	rg.POST("/login", controller.Login)
	rg.POST("/logintoken", controller.LoginToken)
}
