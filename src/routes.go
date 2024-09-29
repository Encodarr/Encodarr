package main

import (
	"net/http"
	"transfigurr/api/routes"
	"transfigurr/constants"
	"transfigurr/interfaces"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(scanService interfaces.ScanServiceInterface, encodeService interfaces.EncodeServiceInterface, metadataService interfaces.MetadataServiceInterface, seriesRepo interfaces.SeriesRepositoryInterface, seasonRepo interfaces.SeasonRepositoryInterface, episodeRepo interfaces.EpisodeRepositoryInterface, movieRepo interfaces.MovieRepositoryInterface, settingRepo interfaces.SettingRepositoryInterface, systemRepo interfaces.SystemRepositoryInterface, profileRepo interfaces.ProfileRepositoryInterface, authRepo interfaces.AuthRepositoryInterface, userRepo interfaces.UserRepositoryInterface, historyRepo interfaces.HistoryRepositoryInterface, eventRepo interfaces.EventRepositoryInterface, codecRepo interfaces.CodecRepositoryInterface) *gin.Engine {
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
	routes.ActionRoutes(api, scanService, metadataService)
	routes.ArtworkRoutes(api)

	routes.WebsocketRoutes(api, encodeService, seriesRepo, movieRepo, profileRepo, settingRepo, systemRepo, historyRepo, eventRepo, codecRepo)

	router.StaticFS("/app", http.Dir(constants.FrontendDistPath))
	router.NoRoute(func(c *gin.Context) {
		c.File(constants.IndexFilePath)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
