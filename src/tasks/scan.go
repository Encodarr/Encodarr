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
			log.Printf("An error occurred scanning movie %s: %v", movieID, r)
		}
	}()

	log.Printf("Starting scan for movie ID: %s", movieID)

	if movieID == "" {
		log.Printf("Movie ID is empty, aborting scan.")
		return
	}

	moviesPath := filepath.Join(constants.MoviesPath, movieID)
	if _, err := os.Stat(moviesPath); os.IsNotExist(err) {
		log.Printf("Movies path does not exist: %s", moviesPath)
		return
	}

	movie, err := movieRepo.GetMovieById(movieID)
	if err != nil {
		log.Printf("Error getting movie: %v", err)
	}

	// Initialize a new movie if it doesn't exist
	if movie.Id == "" {
		log.Printf("Movie with ID %s does not exist, initializing a new movie.", movieID)
		movie = models.Movie{
			Id: movieID,
		}
	}

	log.Printf("Fetched movie: %+v", movie)

	defaultProfile, settingsErr := settingRepo.GetSettingById("defaultProfile")
	if settingsErr != nil {
		log.Printf("Error getting default profile setting: %v", settingsErr)
		return
	}
	log.Printf("Default profile: %+v", defaultProfile)

	if movie.ProfileID == 0 {
		profileID, err := strconv.Atoi(defaultProfile.Value)
		if err != nil {
			log.Printf("Error converting default profile value to int: %v", err)
			return
		}
		movie.ProfileID = profileID
	}
	log.Printf("Movie Profile ID: %d", movie.ProfileID)

	profile, profileErr := profileRepo.GetProfileById(movie.ProfileID)
	if profileErr != nil {
		log.Printf("Error getting profile by ID: %v", profileErr)
		return
	}
	log.Printf("Profile: %+v", profile)

	fileFound := false
	err = filepath.Walk(moviesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error walking the path %s: %v", path, err)
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
			log.Printf("Error analyzing media file %s: %v", path, err)
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
		log.Printf("Error walking the path %s: %v", moviesPath, err)
		return
	}

	log.Printf("Movie after scanning: %+v", movie)
	_, err = movieRepo.UpsertMovie(movie.Id, movie)
	if err != nil {
		log.Printf("Error upserting movie: %v", err)
	}
	log.Printf("Finished scan for movie ID: %s", movieID)
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
				log.Print(episode.Id)
			}
		}
		if episode.Missing && series.Monitored {
			encodeService.Enqueue(models.Item{Type: "episode", Id: episode.Id})
		}

		episodeRepo.UpsertEpisode(seriesID, seasonNumber, episodeNumber, *episode)

		return nil
	})

	if err != nil {
		log.Printf("Error walking the path %v: %v\n", seriesPath, err)
	}

	// Update series properties before saving
	series.SeasonsCount = len(seasons)
	series.EpisodeCount = 0
	series.Size = 0
	series.SpaceSaved = 0

	for _, season := range seasons {
		log.Print(season.MissingEpisodes, " Missing")
		seasonRepo.UpsertSeason(seriesID, season.SeasonNumber, *season)
		series.EpisodeCount += season.EpisodeCount
		series.Size += season.Size
		series.SpaceSaved += season.SpaceSaved
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
