package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/interfaces"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup, authRepo interfaces.AuthRepositoryInterface) {
	controller := controllers.NewAuthController(authRepo)
	rg.GET("/activated", controller.GetActivated)
	rg.POST("/register", controller.Register)
	rg.POST("/login", controller.Login)
	rg.POST("/logintoken", controller.LoginToken)
}
