package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/repository"

	"github.com/gin-gonic/gin"
)

func HistoryRoutes(rg *gin.RouterGroup, historyRepo *repository.HistoryRepository) {
	controller := controllers.NewHistoryController(historyRepo)
	rg.GET("/historys", controller.GetHistories)
	rg.GET("/historys/:historyId", controller.GetHistoryById)
	rg.PUT("/historys/:historyId", controller.UpsertHistory)
	rg.DELETE("/historys/:historyId", controller.DeleteHistoryById)
}
