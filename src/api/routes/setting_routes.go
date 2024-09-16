package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/interfaces"

	"github.com/gin-gonic/gin"
)

func SettingRoutes(rg *gin.RouterGroup, settingRepo interfaces.SettingRepositoryInterface) {
	controller := controllers.NewSettingController(settingRepo)
	rg.GET("/settings", controller.GetSettings)
	rg.GET("/settings/:settingId", controller.GetSettingById)
	rg.PUT("/settings/:settingId", controller.UpsertSetting)
	rg.DELETE("/settings/:settingId", controller.DeleteSettingById)
}
