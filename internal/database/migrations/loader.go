package migrations

import (
	"embed"
	"fmt"
	"log"
	"strings"
)

func init() {
	// Add debug logging to verify embedding
	entries, _ := sqlFiles.ReadDir("sql")
	for _, entry := range entries {
		log.Printf("Embedded SQL file: %s", entry.Name())
	}
}

//go:embed sql/*.sql
var sqlFiles embed.FS

func LoadMigration(filename string) (string, error) {

	data, err := sqlFiles.ReadFile(fmt.Sprintf("sql/%s", filename))
	if err != nil {
		return "", fmt.Errorf("failed to read migration file %s: %w", filename, err)
	}
	return strings.TrimSpace(string(data)), nil
}
