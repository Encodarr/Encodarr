package routes

import (
	"transfigurr/api/controllers"
	"transfigurr/interfaces"

	"github.com/gin-gonic/gin"
)

func EventRoutes(rg *gin.RouterGroup, eventRepo interfaces.EventRepositoryInterface) {
	controller := controllers.NewEventController(eventRepo)
	rg.GET("/events", controller.GetEvents)
	rg.GET("/events/:eventId", controller.GetEventById)
	rg.PUT("/events/:eventId", controller.UpsertEvent)
	rg.DELETE("/events/:eventId", controller.DeleteEventById)
}
