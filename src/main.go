package main

import (
	"log"
	"os"
	"path/filepath"
	"transfigurr/constants"
	"transfigurr/repository"
	"transfigurr/services"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func getParentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(dir)
}

func main() {

	db, err := gorm.Open(constants.DbDriverName, constants.DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	InitDB(db)
	SeedDB(db)
	//db.LogMode(true)
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
	logRepo := repository.NewLoggingRepository(db)

	// services
	metadataService := services.NewMetadataService(seriesRepo, seasonRepo, episodeRepo, movieRepo, settingRepo, systemRepo, profileRepo, authRepo, userRepo, historyRepo, eventRepo, codecRepo)
	metadataService.Startup()

	encodeService := services.NewEncodeService(seriesRepo, seasonRepo, episodeRepo, movieRepo, settingRepo, systemRepo, profileRepo, authRepo, userRepo, historyRepo, eventRepo, codecRepo)
	encodeService.Startup()

	scanService := services.NewScanService(metadataService, encodeService, seriesRepo, seasonRepo, episodeRepo, movieRepo, settingRepo, systemRepo, profileRepo, authRepo, userRepo, historyRepo, eventRepo, codecRepo)
	scanService.Startup()

	logService := services.NewLogService(logRepo, 100)
	logService.Startup("debug")
	currentDir := getParentDir()

	// Create an instance for movies
	moviesWatchdogService := services.NewWatchdogService(100)
	moviesWatchdogService.Startup(filepath.Join(currentDir, "movies"), "movies")
	log.Print("movies")
	// Create an instance for series
	seriesWatchdogService := services.NewWatchdogService(100)
	seriesWatchdogService.Startup(filepath.Join(currentDir, "series"), "series")

	router := SetupRouter(scanService, encodeService, seriesRepo, seasonRepo, episodeRepo, movieRepo, settingRepo, systemRepo, profileRepo, authRepo, userRepo, historyRepo, eventRepo, codecRepo)
	router.Run("0.0.0.0:7889")
}
