package models

type Season struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	SeasonNumber    int       `json:"seasonNumber"`
	EpisodeCount    int       `json:"episodeCount"`
	Size            int       `json:"size"`
	SeriesId        string    `json:"seriesId"`
	SpaceSaved      int       `json:"spaceSaved"`
	MissingEpisodes int       `json:"missingEpisodes"`
	Episodes        []Episode `json:"episodes"`
}
