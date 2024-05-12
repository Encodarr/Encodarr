package main

import (
	"log"
	"transfigurr/constants"
	"transfigurr/repository"
	"transfigurr/services"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := gorm.Open(constants.DbDriverName, constants.DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	InitDB(db)
	SeedDB(db)

	// repos
	seriesRepo := repository.NewSeriesRepository(db)
	seasonRepo := repository.NewSeasonRepository(db)
	episodeRepo := repository.NewEpisodeRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	settingRepo := repository.NewSettingRepository(db)
	systemRepo := repository.NewSystemRepository(db)
	profileRepo := repository.NewProfileRepository(db)
	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	historyRepo := repository.NewHistoryRepository(db)
	eventRepo := repository.NewEventRepository(db)
	codecRepo := repository.NewCodecRepository()

	// services
	scanService := services.NewItemScanService(seriesRepo, seasonRepo, episodeRepo, movieRepo, settingRepo, systemRepo, profileRepo, authRepo, userRepo, historyRepo, eventRepo, codecRepo)
	scanService.Startup()

	router := SetupRouter(scanService, seriesRepo, seasonRepo, episodeRepo, movieRepo, settingRepo, systemRepo, profileRepo, authRepo, userRepo, historyRepo, eventRepo, codecRepo)
	router.Run("0.0.0.0:7889")
}
