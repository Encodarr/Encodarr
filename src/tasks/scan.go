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

	"github.com/shirou/gopsutil/disk"
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

func getDiskSpace(path string) (uint64, uint64, error) {
	log.Printf("Checking disk space for path: %s", path)
	usage, err := disk.Usage(path)
	if err != nil {
		log.Printf("Error getting disk usage for path %s: %v", path, err)
		return 0, 0, err
	}
	log.Printf("Disk usage for path %s: Free: %d, Total: %d", path, usage.Free, usage.Total)
	return usage.Free, usage.Total, nil
}

func ScanSystem(seriesRepo interfaces.SeriesRepositoryInterface, systemRepo interfaces.SystemRepositoryInterface) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("An error occurred scanning system: %v", r)
		}
	}()

	log.Println("Scanning System")

	series, err := seriesRepo.GetSeries()
	if err != nil {
		log.Printf("Error fetching series: %v", err)
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
		log.Printf("Error fetching series disk space: %v", err)
		return
	}

	moviesFreeSpace, moviesTotalSpace, err := getDiskSpace(constants.MoviesPath)
	if err != nil {
		log.Printf("Error fetching Movies disk space: %v", err)
		return
	}
	configFreeSpace, configTotalSpace, err := getDiskSpace(constants.ConfigPath)
	if err != nil {
		log.Printf("Error fetching Config disk space: %v", err)
		return
	}
	transcodeFreeSpace, transcodeTotalSpace, err := getDiskSpace(constants.TranscodeFolder)

	if err != nil {
		log.Printf("Error fetching Transcode disk space: %v", err)
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
