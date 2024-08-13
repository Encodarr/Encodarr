package tasks

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"transfigurr/constants"
	"transfigurr/interfaces"
	"transfigurr/models"
	"transfigurr/utils"
)

var (
	seasonPattern  = regexp.MustCompile(`\d+`)
	episodePattern = regexp.MustCompile(`(?:S(\d{2})E(\d{2})|E(\d{2}))`)
)

func parseEpisodeAndSeasonNumber(file string, folder string) (int, int) {
	match := episodePattern.FindStringSubmatch(file)
	if match == nil {
		return 0, 0
	}
	if match[1] != "" {
		season, _ := strconv.Atoi(match[1])
		episode, _ := strconv.Atoi(match[2])
		return season, episode
	} else {
		parent := filepath.Base(folder)
		seasonNumber := seasonPattern.FindStringSubmatch(parent)
		if parent == "specials" {
			episode, _ := strconv.Atoi(match[3])
			return 0, episode
		}
		if seasonNumber != nil {
			season, _ := strconv.Atoi(seasonNumber[0])
			episode, _ := strconv.Atoi(match[3])
			return season, episode
		} else {
			episode, _ := strconv.Atoi(match[3])
			return 0, episode
		}
	}
}

func ScanMovie(movieID string, movieRepo interfaces.MovieRepositoryInterface, settingRepo interfaces.SettingRepositoryInterface, profileRepo interfaces.ProfileRepositoryInterface) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("An error occurred scanning series %s: %v", movieID, r)
		}
	}()

	movie, err := movieRepo.GetMovieById(movieID)
	if err != nil {
		log.Printf("Error getting movie: %v\n", err)
	}
	if movie.Id == "" {
		movie.Id = movieID
	}
	moviesPath := filepath.Join(constants.MoviesPath, movieID)

	if _, err := os.Stat(moviesPath); os.IsNotExist(err) {
		return
	}

	defaultProfile, settingsErr := settingRepo.GetSettingById("default_profile")
	if settingsErr != nil {
		log.Print(settingsErr)
	}
	if movie.ProfileID == 0 {
		profileID, _ := strconv.Atoi(defaultProfile.Value)
		movie.ProfileID = profileID
	}
	profile, profileErr := profileRepo.GetProfileById(movie.ProfileID)
	if profileErr != nil {
		log.Print(profileErr, profile)
	}
	err = filepath.Walk(moviesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		movie.Filename = info.Name()
		movie.Path = path
		movie.Size = int(info.Size())

		if movie.OriginalSize == 0 {
			movie.OriginalSize = movie.Size
		}

		vcodec, err := utils.AnalyzeMediaFile(path)
		if err != nil {
			return nil
		}
		movie.VideoCodec = vcodec
		if movie.VideoCodec != profile.Codec && profile.Codec != "Any" {
			movie.Missing = true
		}
		return nil
	})

	if err != nil {
		log.Printf("Error walking the path %v: %v\n", moviesPath, err)
		return
	}

	movieRepo.UpsertMovie(movie.Id, movie)
}

func ScanSeries(encodeService interfaces.EncodeServiceInterface, seriesID string, seriesRepo interfaces.SeriesRepositoryInterface, seasonRepo interfaces.SeasonRepositoryInterface, episodeRepo interfaces.EpisodeRepositoryInterface, settingRepo interfaces.SettingRepositoryInterface, profileRepo interfaces.ProfileRepositoryInterface) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("An error occurred scanning series %s: %v", seriesID, r)
		}
	}()

	series, err := seriesRepo.GetSeriesByID(seriesID)
	if err != nil {
		log.Printf("Error getting series: %v\n", err)
	}
	if series.Id == "" {
		series.Id = seriesID
	}
	seriesPath := filepath.Join(constants.SeriesPath, seriesID)

	if _, err := os.Stat(seriesPath); os.IsNotExist(err) {
		return
	}

	defaultProfile, settingsErr := settingRepo.GetSettingById("default_profile")
	if settingsErr != nil {
		log.Print(settingsErr)
	}
	if series.ProfileID == 0 {
		profileID, _ := strconv.Atoi(defaultProfile.Value)
		series.ProfileID = profileID
	}
	profile, profileErr := profileRepo.GetProfileById(series.ProfileID)
	if profileErr != nil {
		log.Print(profileErr, profile)
	}
	seasons := make(map[string]*models.Season)

	err = filepath.Walk(seriesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		seasonNumber, episodeNumber := parseEpisodeAndSeasonNumber(info.Name(), filepath.Dir(path))
		seasonID := seriesID + strconv.Itoa(seasonNumber)

		episode := &models.Episode{
			Id:            seriesID + strconv.Itoa(seasonNumber) + strconv.Itoa(episodeNumber),
			Filename:      info.Name(),
			Path:          path,
			Size:          int(info.Size()),
			EpisodeNumber: episodeNumber,
			SeasonNumber:  seasonNumber,
			SeasonId:      seasonID,
			SeasonName:    filepath.Base(filepath.Dir(path)),
		}

		if episode.OriginalSize == 0 {
			episode.OriginalSize = episode.Size
		}

		vcodec, err := utils.AnalyzeMediaFile(path)
		if err != nil {
			return nil
		}
		episode.VideoCodec = vcodec
		if episode.VideoCodec != profile.Codec && profile.Codec != "Any" {
			episode.Missing = true
			series.MissingEpisodes += 1
		}
		if episode.Missing && series.Monitored {
			encodeService.Enqueue(models.Item{Type: "episode", Id: episode.Id})
		}

		episodeRepo.UpsertEpisode(seriesID, seasonNumber, episodeNumber, *episode)
		if _, ok := seasons[seasonID]; !ok {
			seasons[seasonID] = &models.Season{
				SeasonNumber: seasonNumber,
				Name:         episode.SeasonName,
				EpisodeCount: 1,
				Size:         episode.Size,
			}
		} else {
			seasons[seasonID].EpisodeCount++
			seasons[seasonID].Size += episode.Size
			seasons[seasonID].SpaceSaved += episode.Size
		}

		return nil
	})

	if err != nil {
		log.Printf("Error walking the path %v: %v\n", seriesPath, err)
	}

	for _, data := range seasons {
		seasonRepo.UpsertSeason(seriesID, data.SeasonNumber, *data)
		series.SeasonsCount += 1
		series.EpisodeCount += data.EpisodeCount
		series.Size += data.Size
		series.SpaceSaved += data.SpaceSaved
	}
	seriesRepo.UpsertSeries(series.Id, series)
}
