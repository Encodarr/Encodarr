package models

type Item struct {
	Id            string `json:"id"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	EpisodeName   string `json:"episodeName"`
	ProfileId     int    `json:"profileId"`
	SeriesId      string `json:"seriesId"`
	SeasonNumber  int    `json:"seasonNumber"`
	EpisodeNumber int    `json:"episodeNumber"`
	Codec         string `json:"codec"`
	Size          int    `json:"size"`
}
