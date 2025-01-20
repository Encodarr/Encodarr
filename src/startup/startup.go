package startup

import (
	"database/sql"
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
	"transfigurr/types"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getParentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
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

func writeUptimeToDB(systemRepo interfaces.SystemRepositoryInterface) error {
	uptimeDate := time.Now()
	// Assuming you have a function to set system data

	system := models.System{Id: "start_time", Value: uptimeDate.Format(time.RFC3339)}

	systemRepo.UpsertSystem("start_time", system)
	return nil
}

func Startup() (sqlDB *sql.DB, servicesContainer *types.Services, repositories *types.Repositories) {

	// Ensure the database path exists
	if err := ensureDbPathExists(constants.DbPath); err != nil {
		return
	}

	// Ensure the database path exists
	if err := ensureDbPathExists(constants.DbPath); err != nil {
		return
	}

	db, err := gorm.Open(sqlite.Open(constants.DbPath), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
	})
	if err != nil {
		return
	}
	InitDB(db)
	SeedDB(db)
	//db.LogMode(true)
	// repos

	repositories = repository.NewRepositories(db)

	// services
	eventService := services.NewEventService(repositories.EventRepo, 100)
	eventService.Startup("debug")

	metadataService := services.NewMetadataService(eventService, repositories)
	metadataService.Startup()

	encodeService := services.NewEncodeService(eventService, repositories)
	encodeService.Startup()

	scanService := services.NewScanService(eventService, metadataService, encodeService, repositories)
	scanService.Startup()

	// Create and return the service container
	servicesContainer = &types.Services{
		ScanService:     scanService,
		EncodeService:   encodeService,
		MetadataService: metadataService,
	}

	currentDir := getParentDir()

	// Write uptime to db
	if err := writeUptimeToDB(repositories.SystemRepo); err != nil {
		return
	}

	// Create an instance for movies
	moviesWatchdogService := services.NewWatchdogService(scanService)
	moviesWatchdogService.Startup(filepath.Join(currentDir, "movies"), "movies")
	// Create an instance for series
	seriesWatchdogService := services.NewWatchdogService(scanService)
	seriesWatchdogService.Startup(filepath.Join(currentDir, "series"), "series")

	// scan system
	tasks.ScanSystem(repositories.SeriesRepo, repositories.SystemRepo)

	sqlDB, err = db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB from gorm.DB: %v", err)
	}

	return sqlDB, servicesContainer, repositories
}
