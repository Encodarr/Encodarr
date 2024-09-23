package startup

import (
	"log"
	"os"
	"path/filepath"
	"time"
	"transfigurr/constants"
	"transfigurr/interfaces"
	"transfigurr/models"
	"transfigurr/repository"
	"transfigurr/services"
	"transfigurr/tasks"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getParentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(dir)
}

func ensureDbPathExists(dbPath string) error {

	dir := filepath.Dir(dbPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

// func customOpenAPI() map[string]interface{} {
// 	// Assuming you have a function to get the Gin engine
// 	app := getApp()
// 	if app == nil {
// 		return nil
// 	}
// 	openapiSchema := ginSwagger.WrapHandler(swaggerFiles.Handler)
// 	return openapiSchema
// }

// func writeAPI() error {
// 	openapiSchema := customOpenAPI()
// 	openapiPath := "src/Transfigurr.API.V1"
// 	if _, err := os.Stat(openapiPath); os.IsNotExist(err) {
// 		if err := os.MkdirAll(openapiPath, os.ModePerm); err != nil {
// 			return err
// 		}
// 	}
// 	file, err := os.Create(filepath.Join(openapiPath, "openapi.json"))
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()
// 	return json.NewEncoder(file).Encode(openapiSchema)
// }

func writeUptimeToDB(systemRepo interfaces.SystemRepositoryInterface) error {
	uptimeDate := time.Now()
	// Assuming you have a function to set system data

	system := models.System{Id: "start_time", Value: uptimeDate.Format(time.RFC3339)}

	systemRepo.UpsertSystem("start_time", system)
	return nil
}

func Startup() (db *gorm.DB, scanService interfaces.ScanServiceInterface, encodeService interfaces.EncodeServiceInterface, metadataService interfaces.MetadataServiceInterface, seriesRepo interfaces.SeriesRepositoryInterface, seasonRepo interfaces.SeasonRepositoryInterface, episodeRepo interfaces.EpisodeRepositoryInterface, movieRepo interfaces.MovieRepositoryInterface, settingRepo interfaces.SettingRepositoryInterface, systemRepo interfaces.SystemRepositoryInterface, profileRepo interfaces.ProfileRepositoryInterface, authRepo interfaces.AuthRepositoryInterface, userRepo interfaces.UserRepositoryInterface, historyRepo interfaces.HistoryRepositoryInterface, eventRepo interfaces.EventRepositoryInterface, codecRepo interfaces.CodecRepositoryInterface) {

	// Ensure the database path exists
	if err := ensureDbPathExists(constants.DbPath); err != nil {
		log.Fatalf("Failed to ensure database path exists: %v", err)
	}

	// Ensure the database path exists
	if err := ensureDbPathExists(constants.DbPath); err != nil {
		log.Fatalf("Failed to ensure database path exists: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(constants.DbPath), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
	})
	if err != nil {
		log.Fatal(err)
	}
	InitDB(db)
	SeedDB(db)
	//db.LogMode(true)
	// repos
	seriesRepo = repository.NewSeriesRepository(db)
	seasonRepo = repository.NewSeasonRepository(db)
	episodeRepo = repository.NewEpisodeRepository(db)
	movieRepo = repository.NewMovieRepository(db)
	settingRepo = repository.NewSettingRepository(db)
	systemRepo = repository.NewSystemRepository(db)
	profileRepo = repository.NewProfileRepository(db)
	authRepo = repository.NewAuthRepository(db)
	userRepo = repository.NewUserRepository(db)
	historyRepo = repository.NewHistoryRepository(db)
	eventRepo = repository.NewEventRepository(db)
	codecRepo = repository.NewCodecRepository()

	// services
	eventService := services.NewEventService(eventRepo, 100)
	eventService.Startup("debug")

	metadataService = services.NewMetadataService(eventService, seriesRepo, seasonRepo, episodeRepo, movieRepo, settingRepo, systemRepo, profileRepo, authRepo, userRepo, historyRepo, eventRepo, codecRepo)
	metadataService.Startup()

	encodeService = services.NewEncodeService(eventService, seriesRepo, seasonRepo, episodeRepo, movieRepo, settingRepo, systemRepo, profileRepo, authRepo, userRepo, historyRepo, eventRepo, codecRepo)
	encodeService.Startup()

	scanService = services.NewScanService(eventService, metadataService, encodeService, seriesRepo, seasonRepo, episodeRepo, movieRepo, settingRepo, systemRepo, profileRepo, authRepo, userRepo, historyRepo, eventRepo, codecRepo)
	scanService.Startup()

	currentDir := getParentDir()
	log.Print("currentDir", currentDir)

	// Write uptime to db
	if err := writeUptimeToDB(systemRepo); err != nil {
		log.Fatalf("Failed to write uptime to db: %v", err)
	}

	// Create an instance for movies
	moviesWatchdogService := services.NewWatchdogService(scanService)
	moviesWatchdogService.Startup(filepath.Join(currentDir, "movies"), "movies")
	log.Print("movies")
	// Create an instance for series
	seriesWatchdogService := services.NewWatchdogService(scanService)
	seriesWatchdogService.Startup(filepath.Join(currentDir, "series"), "series")

	// scan system
	tasks.ScanSystem(seriesRepo, systemRepo)
	return db, scanService, encodeService, metadataService, seriesRepo, seasonRepo, episodeRepo, movieRepo, settingRepo, systemRepo, profileRepo, authRepo, userRepo, historyRepo, eventRepo, codecRepo
}
