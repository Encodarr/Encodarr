package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/interfaces"

	"github.com/gin-gonic/gin"
)

func SystemRoutes(rg *gin.RouterGroup, systemRepo interfaces.SystemRepositoryInterface) {
	controller := controllers.NewSystemController(systemRepo)
	rg.GET("/systems", controller.GetSystems)
	rg.GET("/systems/:systemId", controller.GetSystemById)
	rg.PUT("/systems/:systemId", controller.UpsertSystem)
	rg.DELETE("/systems/:systemId", controller.DeleteSystemById)
}
