package services

import (
	"fmt"
	"io/ioutil"
	"sync"
	"transfigurr/constants"
	"transfigurr/interfaces"
	"transfigurr/models"
	"transfigurr/tasks"
	"transfigurr/types"
)

type ScanService struct {
	scanQueue       []models.Item
	scanSet         map[string]struct{}
	mu              sync.Mutex
	cond            *sync.Cond
	metadataService interfaces.MetadataServiceInterface
	encodeService   interfaces.EncodeServiceInterface
	eventService    interfaces.EventServiceInterface
	repositories    *types.Repositories
}

func NewScanService(eventService interfaces.EventServiceInterface, metadataService interfaces.MetadataServiceInterface, encodeService interfaces.EncodeServiceInterface, repositories *types.Repositories) interfaces.ScanServiceInterface {
	service := &ScanService{
		scanQueue:       make([]models.Item, 0),
		scanSet:         make(map[string]struct{}),
		metadataService: metadataService,
		encodeService:   encodeService,
		eventService:    eventService,
		repositories:    repositories,
	}
	service.cond = sync.NewCond(&service.mu)
	return service
}

func (s *ScanService) EnqueueAll() {
	s.EnqueueAllMovies()
	s.EnqueueAllSeries()
}

func (s *ScanService) EnqueueAllMovies() {
	movies, err := s.repositories.MovieRepo.GetMovies()
	if err != nil {
		return
	}
	movieFiles, err := ioutil.ReadDir(constants.MoviesPath)
	if err != nil {
		return
	}
	for _, file := range movieFiles {
		s.Enqueue(models.Item{Id: file.Name(), Type: "movie"})
	}
	for _, movieItem := range movies {
		s.Enqueue(models.Item{Id: movieItem.Id, Type: "movie"})
	}
}

func (s *ScanService) EnqueueAllSeries() {
	series, err := s.repositories.SeriesRepo.GetSeries()
	if err != nil {
		return
	}
	seriesFiles, err := ioutil.ReadDir(constants.SeriesPath)
	if err != nil {
		return
	}
	for _, file := range seriesFiles {
		s.Enqueue(models.Item{Id: file.Name(), Type: "series"})
	}

	for _, seriesItem := range series {
		s.Enqueue(models.Item{Id: seriesItem.Id, Type: "series"})
	}
}

func (s *ScanService) Enqueue(item models.Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	itemID := fmt.Sprintf("%s_%s", item.Type, item.Id)
	if _, ok := s.scanSet[itemID]; !ok {
		s.scanSet[itemID] = struct{}{}
		s.scanQueue = append(s.scanQueue, item)
		s.cond.Signal()
	}
}

func (s *ScanService) process() {
	for {
		s.mu.Lock()
		for len(s.scanQueue) == 0 {
			s.cond.Wait()
		}
		item := s.scanQueue[0]
		s.scanQueue = s.scanQueue[1:]
		s.mu.Unlock()
		s.processItem(item)
		tasks.ScanSystem(s.repositories.SeriesRepo, s.repositories.SystemRepo)
	}
}

func (s *ScanService) processItem(item models.Item) {
	if item.Type == "movie" {
		tasks.ScanMovie(item.Id, s.repositories.MovieRepo, s.repositories.SettingRepo, s.repositories.ProfileRepo)
		tasks.ValidateMovie(item.Id, s.repositories.MovieRepo)
		movie, err := s.repositories.MovieRepo.GetMovieById(item.Id)
		if err != nil {
			return
		}

		if movie.Name == "" {
			s.eventService.Log("INFO", "scan", "Scanning movie: "+item.Id)
			s.metadataService.Enqueue(models.Item{Type: "movie", Id: movie.Id})
		}
		if movie.Missing && movie.Monitored {
			s.encodeService.Enqueue(models.Item{Type: "movie", Id: movie.Id, ProfileId: movie.ProfileID, Codec: movie.VideoCodec, Size: movie.Size})
		}
	} else if item.Type == "series" {
		s.eventService.Log("INFO", "scan", "Scanning series: "+item.Id)
		tasks.ScanSeries(s.encodeService, item.Id, s.repositories.SeriesRepo, s.repositories.SeasonRepo, s.repositories.EpisodeRepo, s.repositories.SettingRepo, s.repositories.ProfileRepo)
		tasks.ValidateSeries(item.Id, s.repositories.SeriesRepo, s.repositories.SeasonRepo, s.repositories.EpisodeRepo)
		series, _ := s.repositories.SeriesRepo.GetSeriesByID(item.Id)

		if series.Name == "" {
			s.metadataService.Enqueue(models.Item{Type: "series", Id: series.Id})
		} else {
			for _, season := range series.Seasons {
				for _, episode := range season.Episodes {
					if episode.EpisodeName == "" {
						s.metadataService.Enqueue(models.Item{Type: "series", Id: series.Id})
						break
					}
				}
			}
		}
	}

	s.mu.Lock()
	delete(s.scanSet, fmt.Sprintf("%s_%s", item.Type, item.Id))
	s.mu.Unlock()
}

func (s *ScanService) Startup() {
	go s.process()
}
