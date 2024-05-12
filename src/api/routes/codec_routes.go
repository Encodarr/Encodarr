package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/repository"

	"github.com/gin-gonic/gin"
)

func CodecRoutes(rg *gin.RouterGroup, codecRepo *repository.CodecRepository) {
	controller := controllers.NewCodecController(codecRepo)
	rg.GET("/codecs", controller.GetCodecs)
	rg.GET("/containers", controller.GetContainers)
	rg.GET("/encoders", controller.GetEncoders)
}
