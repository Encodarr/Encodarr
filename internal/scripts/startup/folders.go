package startup

import (
	"os"
	"path/filepath"
	"time"
	"transfigurr/internal/interfaces/repositories"
	"transfigurr/internal/models"
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

func writeUptimeToDB(systemRepo repositories.SystemRepositoryInterface) error {
	uptimeDate := time.Now()
	// Assuming you have a function to set system data

	system := models.System{Id: "start_time", Value: uptimeDate.Format(time.RFC3339)}

	systemRepo.UpsertSystem("start_time", system)
	return nil
}
