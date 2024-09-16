package main

import (
	"transfigurr/startup"
)

func main() {

	db, scanService, encodeService, seriesRepo, seasonRepo, episodeRepo, movieRepo, settingRepo, systemRepo, profileRepo, authRepo, userRepo, historyRepo, eventRepo, codecRepo := startup.Startup()
	router := SetupRouter(scanService, encodeService, seriesRepo, seasonRepo, episodeRepo, movieRepo, settingRepo, systemRepo, profileRepo, authRepo, userRepo, historyRepo, eventRepo, codecRepo)
	defer db.Close()

	router.Run("0.0.0.0:7889")

}
