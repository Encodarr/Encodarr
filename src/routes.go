package main

import (
	"net/http"
	"transfigurr/api/routes"
	"transfigurr/constants"
	"transfigurr/interfaces"
	"transfigurr/repository"

	"github.com/gin-gonic/gin"
)

func SetupRouter(scanService interfaces.ScanServiceInterface, seriesRepo *repository.SeriesRepository, seasonRepo *repository.SeasonRepository, episodeRepo *repository.EpisodeRepository, movieRepo *repository.MovieRepository, settingRepo *repository.SettingRepository, systemRepo *repository.SystemRepository, profileRepo *repository.ProfileRepository, authRepo *repository.AuthRepository, userRepo *repository.UserRepository, historyRepo *repository.HistoryRepository, eventRepo *repository.EventRepository, codecRepo *repository.CodecRepository) *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	api := router.Group("/api")
	routes.SeriesRoutes(api, seriesRepo)
	routes.SeasonRoutes(api, seasonRepo)
	routes.EpisodeRoutes(api, episodeRepo)
	routes.MovieRoutes(api, movieRepo)
	routes.SettingRoutes(api, settingRepo)
	routes.SystemRoutes(api, systemRepo)
	routes.ProfileRoutes(api, profileRepo)
	routes.AuthRoutes(api, authRepo)
	routes.UserRoutes(api, userRepo)
	routes.HistoryRoutes(api, historyRepo)
	routes.EventRoutes(api, eventRepo)
	routes.CodecRoutes(api, codecRepo)
	routes.ActionRoutes(api, scanService)
	routes.ArtworkRoutes(api)

	router.StaticFS("/app", http.Dir(constants.FrontendDistPath))
	router.NoRoute(func(c *gin.Context) {
		c.File(constants.IndexFilePath)
	})
	return router
}
