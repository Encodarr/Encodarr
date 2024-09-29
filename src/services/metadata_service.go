package services

import (
	"fmt"
	"sync"
	"transfigurr/interfaces"
	"transfigurr/models"
	"transfigurr/utils"
)

type MetadataService struct {
	metadataQueue []models.Item
	metadataSet   map[string]struct{}
	mu            sync.Mutex
	cond          *sync.Cond
	eventService  interfaces.EventServiceInterface
	seriesRepo    interfaces.SeriesRepositoryInterface
	seasonRepo    interfaces.SeasonRepositoryInterface
	episodeRepo   interfaces.EpisodeRepositoryInterface
	movieRepo     interfaces.MovieRepositoryInterface
	settingRepo   interfaces.SettingRepositoryInterface
	systemRepo    interfaces.SystemRepositoryInterface
	profileRepo   interfaces.ProfileRepositoryInterface
	authRepo      interfaces.AuthRepositoryInterface
	userRepo      interfaces.UserRepositoryInterface
	historyRepo   interfaces.HistoryRepositoryInterface
	eventRepo     interfaces.EventRepositoryInterface
	codecRepo     interfaces.CodecRepositoryInterface
}

func NewMetadataService(eventService interfaces.EventServiceInterface, seriesRepo interfaces.SeriesRepositoryInterface, seasonRepo interfaces.SeasonRepositoryInterface, episodeRepo interfaces.EpisodeRepositoryInterface, movieRepo interfaces.MovieRepositoryInterface, settingRepo interfaces.SettingRepositoryInterface, systemRepo interfaces.SystemRepositoryInterface, profileRepo interfaces.ProfileRepositoryInterface, authRepo interfaces.AuthRepositoryInterface, userRepo interfaces.UserRepositoryInterface, historyRepo interfaces.HistoryRepositoryInterface, eventRepo interfaces.EventRepositoryInterface, codecRepo interfaces.CodecRepositoryInterface) interfaces.MetadataServiceInterface {
	service := &MetadataService{
		metadataQueue: make([]models.Item, 0),
		metadataSet:   make(map[string]struct{}),
		eventService:  eventService,
		seriesRepo:    seriesRepo,
		seasonRepo:    seasonRepo,
		episodeRepo:   episodeRepo,
		movieRepo:     movieRepo,
		settingRepo:   settingRepo,
		systemRepo:    systemRepo,
		profileRepo:   profileRepo,
		authRepo:      authRepo,
		userRepo:      userRepo,
		historyRepo:   historyRepo,
		eventRepo:     eventRepo,
		codecRepo:     codecRepo,
	}
	service.cond = sync.NewCond(&service.mu)
	return service
}

func (s *MetadataService) EnqueueAll() {
	series, err := s.seriesRepo.GetSeries()
	if err != nil {
	}
	for _, series := range series {
		s.Enqueue(models.Item{Type: "series", Id: series.Id})
	}
	movies, err := s.movieRepo.GetMovies()
	if err != nil {
	}
	for _, m := range movies {
		s.Enqueue(models.Item{Type: "movie", Id: m.Id})
	}
}

func (s *MetadataService) Enqueue(item models.Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	itemID := fmt.Sprintf("%s_%s", item.Type, item.Id)
	if _, ok := s.metadataSet[itemID]; !ok {
		s.metadataSet[itemID] = struct{}{}
		s.metadataQueue = append(s.metadataQueue, item)
		s.cond.Signal()
	}
}

func (s *MetadataService) process() {
	for {
		s.mu.Lock()
		for len(s.metadataQueue) == 0 {
			s.cond.Wait()
		}
		item := s.metadataQueue[0]
		s.metadataQueue = s.metadataQueue[1:]
		s.mu.Unlock()
		s.processItem(item)
	}
}

func (s *MetadataService) processItem(item models.Item) {
	if item.Type == "movie" {
		movie, err := s.movieRepo.GetMovieById(item.Id)
		if err != nil {
		}
		movie, err = utils.GetMovieMetadata(movie)
		if err != nil {
		}
		s.movieRepo.UpsertMovie(movie.Id, movie)
	} else if item.Type == "series" {
		series, err := s.seriesRepo.GetSeriesByID(item.Id)
		if err != nil {
		}
		series, err = utils.GetSeriesMetadata(series)

		s.seriesRepo.UpsertSeries(series.Id, series)
		for _, season := range series.Seasons {
			for _, episode := range season.Episodes {
				s.episodeRepo.UpsertEpisode(series.Id, season.SeasonNumber, episode.EpisodeNumber, episode)
			}
		}

	}

	s.mu.Lock()
	delete(s.metadataSet, fmt.Sprintf("%s_%s", item.Type, item.Id))
	s.mu.Unlock()
}

func (s *MetadataService) Startup() {
	go s.process()
}
