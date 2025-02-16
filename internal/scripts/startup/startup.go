package startup

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"time"
	"transfigurr/internal/config"
	"transfigurr/internal/database/migrations"
	"transfigurr/internal/repositories"
	"transfigurr/internal/services"
	"transfigurr/internal/types"
	"transfigurr/internal/utils"

	_ "modernc.org/sqlite"
)

func Startup() (sqlDB *sql.DB, servicesContainer *types.Services, repos *types.Repositories) {
	// Ensure the database path exists
	if err := ensureDbPathExists(config.DbPath); err != nil {
		log.Fatal(err)
	}

	// Increase timeout for initial setup
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Configure database with optimized settings
	dsn := fmt.Sprintf("%s?"+
		"_journal=WAL&"+
		"_timeout=30000&"+ // Increase timeout to 30 seconds
		"_busy_timeout=30000&"+ // Increase busy timeout
		"cache=shared&"+
		"_synchronous=NORMAL&"+
		"_journal_mode=WAL&"+
		"_locking_mode=NORMAL&"+
		"_page_size=4096&"+
		"_mmap_size=67108864&"+ // Reduce mmap size to 64MB
		"_max_page_count=2147483646", // Max page count
		config.DbPath)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Print("Error opening database: %v", err)
		log.Fatal(err)
	}

	// Configure connection pool with higher limits
	db.SetMaxOpenConns(4)            // Allow more concurrent connections
	db.SetMaxIdleConns(2)            // Keep more idle connections
	db.SetConnMaxLifetime(time.Hour) // Keep connections alive longer

	// Enable Write-Ahead Logging
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		log.Fatal(err)
	}

	// Optimize performance
	pragmas := []string{
		"PRAGMA synchronous=NORMAL",
		"PRAGMA temp_store=MEMORY",
		"PRAGMA cache_size=10000",
		"PRAGMA foreign_keys=ON",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			log.Printf("Warning: failed to set %s: %v", pragma, err)
		}
	}

	migrationService := services.NewMigrationService(db)

	migrations := []services.Migration{
		{
			ID:   1,
			Name: "initial_tables",
			SQL:  mustLoadMigration("0000001_initial_tables.sql"),
		},
		{
			ID:   2,
			Name: "default_data",
			SQL:  mustLoadMigration("0000002_default_data.sql"),
		},
		{
			ID:   3,
			Name: "insert_test_data",
			SQL:  mustLoadMigration("0000003_sub_data.sql"),
		},
	}
	log.Print("Running migrations...")

	if err := migrationService.RunMigrations(ctx, migrations); err != nil {
		log.Fatal(err)
	}
	log.Print("Migrations completed successfully")

	repos = repositories.NewRepositories(db)
	if repos == nil {
		log.Fatal("Failed to initialize repositories")
	}

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

	return sqlDB, servicesContainer, repos
}

func mustLoadMigration(filename string) string {
	sql, err := migrations.LoadMigration(filename)
	if err != nil {
		log.Printf("Error loading migration: %v", err) // Add debug logging
		log.Fatalf("Failed to load migration %s: %v", filename, err)
	}
	return sql
}
