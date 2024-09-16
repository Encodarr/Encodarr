package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/interfaces"

	"github.com/gin-gonic/gin"
)

func UserRoutes(rg *gin.RouterGroup, userRepo interfaces.UserRepositoryInterface) {
	controller := controllers.NewUserController(userRepo)
	rg.GET("/user", controller.GetUsers)
	rg.POST("/user", controller.UpdateUser)
}
