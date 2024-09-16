package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/interfaces"

	"github.com/gin-gonic/gin"
)

func CodecRoutes(rg *gin.RouterGroup, codecRepo interfaces.CodecRepositoryInterface) {
	controller := controllers.NewCodecController(codecRepo)
	rg.GET("/codecs", controller.GetCodecs)
	rg.GET("/containers", controller.GetContainers)
	rg.GET("/encoders", controller.GetEncoders)
}
