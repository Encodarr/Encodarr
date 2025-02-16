package database

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

const (
	maxRetries = 5
	retryDelay = 100 * time.Millisecond
)

func OpenDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath+"?_txlock=immediate&_timeout=5000&_journal_mode=WAL&_synchronous=NORMAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}

	// Set connection pool limits
	db.SetMaxOpenConns(1) // Single connection to prevent lock conflicts
	db.SetMaxIdleConns(1) // Keep connection alive

	// Configure SQLite for better concurrency and performance
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA foreign_keys=ON",
		"PRAGMA busy_timeout=5000", // Reduced from 10000
		"PRAGMA cache_size=-64000",
		"PRAGMA temp_store=MEMORY",
		"PRAGMA mmap_size=268435456",
		"PRAGMA page_size=4096",
		"PRAGMA wal_autocheckpoint=1000",
		"PRAGMA busy_handler=1",
		// Remove exclusive locking mode
		// "PRAGMA locking_mode=EXCLUSIVE",
		"PRAGMA journal_size_limit=67108864",
	}

	// Apply pragmas with retry logic
	for _, pragma := range pragmas {
		for i := 0; i < maxRetries; i++ {
			if _, err := db.Exec(pragma); err == nil {
				break
			} else if i == maxRetries-1 {
				return nil, err
			}
			time.Sleep(retryDelay)
		}
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
