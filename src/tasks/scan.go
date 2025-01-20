package tasks

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"syscall"
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

	if movieID == "" {
		return
	}

	moviesPath := filepath.Join(constants.MoviesPath, movieID)
	if _, err := os.Stat(moviesPath); os.IsNotExist(err) {
		return
	}

	movie, err := movieRepo.GetMovieById(movieID)
	if err != nil {
		return
	}

	// Initialize a new movie if it doesn't exist
	if movie.Id == "" {
		movie = models.Movie{
			Id: movieID,
		}
	}

	defaultProfile, settingsErr := settingRepo.GetSettingById("defaultProfile")
	if settingsErr != nil {
		return
	}

	if movie.ProfileID == 0 {
		profileID, err := strconv.Atoi(defaultProfile.Value)
		if err != nil {
			return
		}
		movie.ProfileID = profileID
	}

	profile, profileErr := profileRepo.GetProfileById(movie.ProfileID)
	if profileErr != nil {
		return
	}

	fileFound := false
	err = filepath.Walk(moviesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		fileFound = true
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

	if !fileFound {
		movie.Filename = ""
		movie.Path = ""
		movie.Size = 0
		movie.VideoCodec = ""
		movie.OriginalSize = 0
		movie.VideoCodec = ""
	}

	if err != nil {
		return
	}

	_, err = movieRepo.UpsertMovie(movie.Id, movie)
	if err != nil {
		return
	}

}

func ScanSeries(encodeService interfaces.EncodeServiceInterface, seriesID string, seriesRepo interfaces.SeriesRepositoryInterface, seasonRepo interfaces.SeasonRepositoryInterface, episodeRepo interfaces.EpisodeRepositoryInterface, settingRepo interfaces.SettingRepositoryInterface, profileRepo interfaces.ProfileRepositoryInterface) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	series, err := seriesRepo.GetSeriesByID(seriesID)
	if err != nil {
	}
	series.MissingEpisodes = 0
	if series.Id == "" {
		series.Id = seriesID
	}
	seriesPath := filepath.Join(constants.SeriesPath, seriesID)

	if _, err := os.Stat(seriesPath); os.IsNotExist(err) {
		return
	}

	defaultProfile, settingsErr := settingRepo.GetSettingById("defaultProfile")
	if settingsErr != nil {
	}
	if series.ProfileID == 0 {
		profileID, _ := strconv.Atoi(defaultProfile.Value)
		series.ProfileID = profileID
	}
	profile, profileErr := profileRepo.GetProfileById(series.ProfileID)
	if profileErr != nil {
		return
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

		if episode.VideoCodec != profile.Codec && profile.Codec != "Any" {
			episode.Missing = true
			series.MissingEpisodes += 1
			if _, ok := seasons[seasonID]; ok {
				seasons[seasonID].MissingEpisodes += 1
			}
		}
		if episode.Missing && series.Monitored {
			encodeService.Enqueue(models.Item{Type: "episode", Id: episode.Id, ProfileId: series.ProfileID, SeriesId: seriesID, SeasonNumber: seasonNumber, EpisodeNumber: episodeNumber, Name: series.Id, Size: episode.Size, Codec: episode.VideoCodec})
		}

		episodeRepo.UpsertEpisode(seriesID, seasonNumber, episodeNumber, *episode)

		return nil
	})
	if err != nil {
		return
	}

	// Update series properties before saving
	series.SeasonsCount = len(seasons)
	series.EpisodeCount = 0
	series.Size = 0
	series.SpaceSaved = 0

	for _, season := range seasons {
		seasonRepo.UpsertSeason(seriesID, season.SeasonNumber, *season)
		series.EpisodeCount += season.EpisodeCount
		series.Size += season.Size
		series.SpaceSaved += season.SpaceSaved
	}

	seriesRepo.UpsertSeries(series.Id, series)
}

func getDiskSpace(path string) (uint64, uint64, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return 0, 0, err
	}
	free := stat.Bavail * uint64(stat.Bsize)
	total := stat.Blocks * uint64(stat.Bsize)
	return free, total, nil
}

func ScanSystem(seriesRepo interfaces.SeriesRepositoryInterface, systemRepo interfaces.SystemRepositoryInterface) {

	series, err := seriesRepo.GetSeries()
	if err != nil {
		return
	}

	seriesCount := 0
	episodeCount := 0
	fileCount := 0
	sizeOnDisk := 0
	monitoredCount := 0
	unmonitoredCount := 0
	endedCount := 0
	continuingCount := 0
	spaceSaved := 0

	seriesFreeSpace, seriesTotalSpace, err := getDiskSpace(constants.SeriesPath)
	if err != nil {
		log.Print("seriesDisk", err)
		return
	}

	moviesFreeSpace, moviesTotalSpace, err := getDiskSpace(constants.MoviesPath)
	if err != nil {
		log.Print("moviesDisk", err)
		return
	}
	configFreeSpace, configTotalSpace, err := getDiskSpace(constants.ConfigPath)
	if err != nil {
		log.Print("configDisk", err)
		return
	}
	transcodeFreeSpace, transcodeTotalSpace, err := getDiskSpace(constants.TranscodeFolder)

	if err != nil {
		log.Print("transDisk", err)
		return
	}

	for id := range series {
		s := series[id]
		seriesCount++
		sizeOnDisk += s.Size
		spaceSaved += s.SpaceSaved
		episodeCount += s.EpisodeCount
		fileCount += s.EpisodeCount
		if s.Monitored {
			monitoredCount++
		} else {
			unmonitoredCount++
		}
		if s.Status == "Ended" {
			endedCount++
		} else {
			continuingCount++
		}
	}

	systemRepo.UpsertSystem("series_count", models.System{Id: "series_count", Value: strconv.Itoa(seriesCount)})
	systemRepo.UpsertSystem("episode_count", models.System{Id: "episode_count", Value: strconv.Itoa(episodeCount)})
	systemRepo.UpsertSystem("files_count", models.System{Id: "files_count", Value: strconv.Itoa(fileCount)})
	systemRepo.UpsertSystem("size_on_disk", models.System{Id: "size_on_disk", Value: strconv.Itoa(sizeOnDisk)})
	systemRepo.UpsertSystem("space_saved", models.System{Id: "space_saved", Value: strconv.Itoa(spaceSaved)})
	systemRepo.UpsertSystem("monitored_count", models.System{Id: "monitored_count", Value: strconv.Itoa(monitoredCount)})
	systemRepo.UpsertSystem("unmonitored_count", models.System{Id: "unmonitored_count", Value: strconv.Itoa(unmonitoredCount)})
	systemRepo.UpsertSystem("ended_count", models.System{Id: "ended_count", Value: strconv.Itoa(endedCount)})
	systemRepo.UpsertSystem("continuing_count", models.System{Id: "continuing_count", Value: strconv.Itoa(continuingCount)})
	systemRepo.UpsertSystem("series_total_space", models.System{Id: "series_total_space", Value: strconv.Itoa(int(seriesTotalSpace))})
	systemRepo.UpsertSystem("series_free_space", models.System{Id: "series_free_space", Value: strconv.Itoa(int(seriesFreeSpace))})
	systemRepo.UpsertSystem("movies_total_space", models.System{Id: "movies_total_space", Value: strconv.Itoa(int(moviesTotalSpace))})
	systemRepo.UpsertSystem("movies_free_space", models.System{Id: "movies_free_space", Value: strconv.Itoa(int(moviesFreeSpace))})
	systemRepo.UpsertSystem("config_total_space", models.System{Id: "config_total_space", Value: strconv.Itoa(int(configTotalSpace))})
	systemRepo.UpsertSystem("config_free_space", models.System{Id: "config_free_space", Value: strconv.Itoa(int(configFreeSpace))})
	systemRepo.UpsertSystem("transcode_total_space", models.System{Id: "transcode_total_space", Value: strconv.Itoa(int(transcodeTotalSpace))})
	systemRepo.UpsertSystem("transcode_free_space", models.System{Id: "transcode_free_space", Value: strconv.Itoa(int(transcodeFreeSpace))})
}
