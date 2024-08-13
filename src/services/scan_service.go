package services

import (
	"fmt"
	"log"
	"sync"
	"time"
	"transfigurr/interfaces"
	"transfigurr/models"
	"transfigurr/tasks"
)

type ScanService struct {
	scanQueue       chan models.Item
	scanSet         map[string]struct{}
	metadataService interfaces.MetadataServiceInterface
	encodeService   interfaces.EncodeServiceInterface
	mu              sync.Mutex
	seriesRepo      interfaces.SeriesRepositoryInterface
	seasonRepo      interfaces.SeasonRepositoryInterface
	episodeRepo     interfaces.EpisodeRepositoryInterface
	movieRepo       interfaces.MovieRepositoryInterface
	settingRepo     interfaces.SettingRepositoryInterface
	systemRepo      interfaces.SystemRepositoryInterface
	profileRepo     interfaces.ProfileRepositoryInterface
	authRepo        interfaces.AuthRepositoryInterface
	userRepo        interfaces.UserRepositoryInterface
	historyRepo     interfaces.HistoryRepositoryInterface
	eventRepo       interfaces.EventRepositoryInterface
	codecRepo       interfaces.CodecRepositoryInterface
}

func NewScanService(metadataService interfaces.MetadataServiceInterface, encodeService interfaces.EncodeServiceInterface, seriesRepo interfaces.SeriesRepositoryInterface, seasonRepo interfaces.SeasonRepositoryInterface, episodeRepo interfaces.EpisodeRepositoryInterface, movieRepo interfaces.MovieRepositoryInterface, settingRepo interfaces.SettingRepositoryInterface, systemRepo interfaces.SystemRepositoryInterface, profileRepo interfaces.ProfileRepositoryInterface, authRepo interfaces.AuthRepositoryInterface, userRepo interfaces.UserRepositoryInterface, historyRepo interfaces.HistoryRepositoryInterface, eventRepo interfaces.EventRepositoryInterface, codecRepo interfaces.CodecRepositoryInterface) interfaces.ScanServiceInterface {
	return &ScanService{
		scanQueue:       make(chan models.Item, 100),
		scanSet:         make(map[string]struct{}),
		metadataService: metadataService,
		encodeService:   encodeService,
		seriesRepo:      seriesRepo,
		seasonRepo:      seasonRepo,
		episodeRepo:     episodeRepo,
		movieRepo:       movieRepo,
		settingRepo:     settingRepo,
		systemRepo:      systemRepo,
		profileRepo:     profileRepo,
		authRepo:        authRepo,
		userRepo:        userRepo,
		historyRepo:     historyRepo,
		eventRepo:       eventRepo,
		codecRepo:       codecRepo,
	}
}

func (s *ScanService) Enqueue(item models.Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	itemID := fmt.Sprintf("%s_%s", item.Type, item.Id)
	if _, ok := s.scanSet[itemID]; !ok {
		s.scanSet[itemID] = struct{}{}
		s.scanQueue <- item
	}
}

func (s *ScanService) process() {
	for {
		select {
		case item, ok := <-s.scanQueue:
			if ok {
				s.processItem(item)
			}
		case <-time.After(1 * time.Second):
			continue
		}
	}
}

func (s *ScanService) processItem(item models.Item) {
	if item.Type == "movie" {
		tasks.ScanMovie(item.Id, s.movieRepo, s.settingRepo, s.profileRepo)
		movie, err := s.movieRepo.GetMovieById(item.Id)
		if err != nil {
			log.Print(err)
		}

		if movie.Name == "" {
			s.metadataService.Enqueue(models.Item{Type: "movie", Id: movie.Id})
		}
		if movie.Missing && movie.Monitored {
			s.encodeService.Enqueue(models.Item{Type: "movie", Id: movie.Id})
		}
	} else if item.Type == "series" {
		tasks.ScanSeries(s.encodeService, item.Id, s.seriesRepo, s.seasonRepo, s.episodeRepo, s.settingRepo, s.profileRepo)
		series, err := s.seriesRepo.GetSeriesByID(item.Id)
		if err != nil {
			log.Print(err)
		}

		if series.Name == "" {
			s.metadataService.Enqueue(models.Item{Type: "series", Id: series.Id})
		}
	}

	s.mu.Lock()
	delete(s.scanSet, fmt.Sprintf("%s_%s", item.Type, item.Id))
	s.mu.Unlock()
}

func (s *ScanService) Startup() {
	go s.process()
}
