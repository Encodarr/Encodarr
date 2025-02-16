package models

type Series struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	ReleaseDate     string   `json:"releaseDate"`
	Genre           string   `json:"genre"`
	Status          string   `json:"status"`
	LastAirDate     string   `json:"lastAirDate"`
	Networks        string   `json:"networks"`
	Overview        string   `json:"overview"`
	ProfileID       int      `json:"profileId"`
	Monitored       bool     `json:"monitored"`
	EpisodeCount    int      `json:"episodeCount"`
	Size            int      `json:"size"`
	SeasonsCount    int      `json:"seasonsCount"`
	SpaceSaved      int      `json:"spaceSaved"`
	MissingEpisodes int      `json:"missingEpisodes"`
	Runtime         int      `json:"runtime"`
	Seasons         []Season `json:"seasons"`
}
