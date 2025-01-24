package startup

import (
	"database/sql"
	"log"
	"path/filepath"
	"transfigurr/internal/config"
	"transfigurr/internal/repositories"
	"transfigurr/internal/services"
	"transfigurr/internal/types"
	"transfigurr/internal/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Startup() (sqlDB *sql.DB, servicesContainer *types.Services, repos *types.Repositories) {

	// Ensure the database path exists
	if err := ensureDbPathExists(config.DbPath); err != nil {
		return
	}

	// Ensure the database path exists
	if err := ensureDbPathExists(config.DbPath); err != nil {
		return
	}

	db, err := gorm.Open(sqlite.Open(config.DbPath), &gorm.Config{
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

	repos = repositories.NewRepositories(db)

	// services
	eventService := services.NewEventService(repos.EventRepo, 100)
	eventService.Startup("debug")

	metadataService := services.NewMetadataService(eventService, repos)
	metadataService.Startup()

	encodeService := services.NewEncodeService(eventService, repos)
	encodeService.Startup()

	scanService := services.NewScanService(eventService, metadataService, encodeService, repos)
	scanService.Startup()

	// Create and return the service container
	servicesContainer = &types.Services{
		ScanService:     scanService,
		EncodeService:   encodeService,
		MetadataService: metadataService,
	}

	currentDir := getParentDir()

	// Write uptime to db
	if err := writeUptimeToDB(repos.SystemRepo); err != nil {
		return
	}

	// Create an instance for movies
	moviesWatchdogService := services.NewWatchdogService(scanService)
	moviesWatchdogService.Startup(filepath.Join(currentDir, "movies"), "movies")
	// Create an instance for series
	seriesWatchdogService := services.NewWatchdogService(scanService)
	seriesWatchdogService.Startup(filepath.Join(currentDir, "series"), "series")

	// scan system
	utils.ScanSystem(repos.SeriesRepo, repos.SystemRepo)

	sqlDB, err = db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB from gorm.DB: %v", err)
	}

	return sqlDB, servicesContainer, repos
}
