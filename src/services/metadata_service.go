package services

import (
	"fmt"
	"log"
	"sync"
	"time"
	"transfigurr/interfaces"
	"transfigurr/models"
	"transfigurr/utils"
)

type MetadataService struct {
	metadataQueue chan models.Item
	metadataSet   map[string]struct{}
	mu            sync.Mutex
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

func NewMetadataService(seriesRepo interfaces.SeriesRepositoryInterface, seasonRepo interfaces.SeasonRepositoryInterface, episodeRepo interfaces.EpisodeRepositoryInterface, movieRepo interfaces.MovieRepositoryInterface, settingRepo interfaces.SettingRepositoryInterface, systemRepo interfaces.SystemRepositoryInterface, profileRepo interfaces.ProfileRepositoryInterface, authRepo interfaces.AuthRepositoryInterface, userRepo interfaces.UserRepositoryInterface, historyRepo interfaces.HistoryRepositoryInterface, eventRepo interfaces.EventRepositoryInterface, codecRepo interfaces.CodecRepositoryInterface) interfaces.MetadataServiceInterface {
	return &MetadataService{
		metadataQueue: make(chan models.Item),
		metadataSet:   make(map[string]struct{}),
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
}

func (s *MetadataService) Enqueue(item models.Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	itemID := fmt.Sprintf("%s_%s", item.Type, item.Id)
	if _, ok := s.metadataSet[itemID]; !ok {
		s.metadataSet[itemID] = struct{}{}
		s.metadataQueue <- item
	}
}

func (s *MetadataService) process() {
	for {
		select {
		case item, ok := <-s.metadataQueue:
			if ok {
				s.processItem(item)
			}
		case <-time.After(1 * time.Second):
			continue
		}
	}
}

func (s *MetadataService) processItem(item models.Item) {
	if item.Type == "movie" {
		movie, err := s.movieRepo.GetMovieById(item.Id)
		if err != nil {
			log.Print(err)
		}
		movie, err = utils.GetMovieMetadata(movie)
		if err != nil {
			log.Print(err)
		}
		s.movieRepo.UpsertMovie(movie.Id, movie)
	} else if item.Type == "series" {
		series, err := s.seriesRepo.GetSeriesByID(item.Id)
		if err != nil {
			log.Print(err)
		}
		series, err = utils.GetSeriesMetadata(series)
		if err != nil {
			log.Print(err)
		}
		s.seriesRepo.UpsertSeries(series.Id, series)
	}

	s.mu.Lock()
	delete(s.metadataSet, fmt.Sprintf("%s_%s", item.Type, item.Id))
	s.mu.Unlock()
}

func (s *MetadataService) Startup() {
	go s.process()
}
