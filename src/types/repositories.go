package types

import "transfigurr/interfaces"

type Repositories struct {
	SeriesRepo  interfaces.SeriesRepositoryInterface
	SeasonRepo  interfaces.SeasonRepositoryInterface
	EpisodeRepo interfaces.EpisodeRepositoryInterface
	MovieRepo   interfaces.MovieRepositoryInterface
	SettingRepo interfaces.SettingRepositoryInterface
	SystemRepo  interfaces.SystemRepositoryInterface
	ProfileRepo interfaces.ProfileRepositoryInterface
	AuthRepo    interfaces.AuthRepositoryInterface
	UserRepo    interfaces.UserRepositoryInterface
	HistoryRepo interfaces.HistoryRepositoryInterface
	EventRepo   interfaces.EventRepositoryInterface
	CodecRepo   interfaces.CodecRepositoryInterface
}
