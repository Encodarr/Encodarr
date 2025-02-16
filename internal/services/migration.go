package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Migration struct {
	ID        int
	Name      string
	SQL       string
	Timestamp time.Time
}

type MigrationService struct {
	db *sql.DB
}

func NewMigrationService(db *sql.DB) *MigrationService {
	return &MigrationService{db: db}
}

func (m *MigrationService) createMigrationsTable(ctx context.Context) error {
	query := `
        CREATE TABLE IF NOT EXISTS migrations (
            id INTEGER PRIMARY KEY,
            name TEXT NOT NULL,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `

	_, err := m.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	return nil
}

func (m *MigrationService) RunMigrations(ctx context.Context, migrations []Migration) error {
	log.Print("Creating migrations table...")
	if err := m.createMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Begin transaction
	log.Print("Beginning migration transaction...")
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get applied migrations
	log.Print("Checking for previously applied migrations...")
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Run each migration in order
	for _, migration := range migrations {
		log.Printf("Processing migration %d: %s", migration.ID, migration.Name)
		if !isApplied(applied, migration.Name) {
			log.Printf("Executing migration: %s", migration.Name)
			if err := m.executeMigration(ctx, tx, migration); err != nil {
				return fmt.Errorf("migration %s failed: %w", migration.Name, err)
			}
			log.Printf("Successfully executed migration: %s", migration.Name)
		} else {
			log.Printf("Skipping already applied migration: %s", migration.Name)
		}
	}

	// Commit transaction
	log.Print("Committing migrations...")
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migrations: %w", err)
	}

	return nil
}

func (m *MigrationService) getAppliedMigrations(ctx context.Context) ([]string, error) {
	rows, err := m.db.QueryContext(ctx, "SELECT name FROM migrations ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applied []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		applied = append(applied, name)
	}
	return applied, nil
}

func (m *MigrationService) executeMigration(ctx context.Context, tx *sql.Tx, migration Migration) error {
	// Execute the migration SQL
	if _, err := tx.ExecContext(ctx, migration.SQL); err != nil {
		return err
	}

	// Record the migration
	if _, err := tx.ExecContext(ctx,
		"INSERT INTO migrations (name) VALUES (?)",
		migration.Name,
	); err != nil {
		return err
	}

	return nil
}

func isApplied(applied []string, name string) bool {
	for _, a := range applied {
		if a == name {
			return true
		}
	}
	return false
}
