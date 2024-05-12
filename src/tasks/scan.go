package tasks

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"transfigurr/constants"
	"transfigurr/interfaces"
	"transfigurr/models"
)

func ScanMovie(movieID string, movieRepo interfaces.MovieRepositoryInterface, settingRepo interfaces.SettingRepositoryInterface, profileRepo interfaces.ProfileRepositoryInterface) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("An error occurred scanning series %s: %v", movieID, r)
		}
	}()

	log.Printf("Scanning movie: %s", movieID)
	movie := models.Movie{Id: movieID}
	moviesPath := filepath.Join(constants.MoviesPath, movieID)

	if _, err := os.Stat(moviesPath); os.IsNotExist(err) {
		return
	}

	missingMetadata := false
	settings, settingsErr := settingRepo.GetAllSettings()
	if settingsErr != nil {
		log.Print(settingsErr)
	}

	if movie.ProfileID == 0 {
		log.Print(settings)
	}
	if movie.Name == "" {
		missingMetadata = true
	}
	log.Print("here")
	profile, profileErr := profileRepo.GetProfileById(fmt.Sprint(movie.ProfileID))
	if profileErr != nil {
		log.Print(profileErr)
		log.Print(profile)
	}

	err := filepath.Walk(moviesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if movie.Name == "" {
			missingMetadata = true
		}

		movie.Filename = info.Name()
		movie.Path = path
		return nil
	})

	if err != nil {
		log.Printf("Error walking the path %v: %v\n", moviesPath, err)
		return
	}
	log.Print(missingMetadata)
	movieRepo.UpsertMovie(movie.Id, movie)

	//if movie.Missing && movie.Monitored != 0 {
	//encodeServiceEnqueue(movie)
	//}

	//if missingMetadata {
	//metadataServiceEnqueue(movieID, "movie")
	//}
}
