package models

type Episode struct {
	Id            string `json:"id"`
	SeriesId      string `json:"seriesId"`
	SeasonId      string `json:"seasonId"`
	EpisodeNumber int    `json:"episodeNumber"`
	SeasonName    string `json:"seasonName"`
	SeasonNumber  int    `json:"seasonNumber"`
	Filename      string `json:"filename"`
	EpisodeName   string `json:"episodeName"`
	VideoCodec    string `json:"videoCodec"`
	AirDate       string `json:"airDate"`
	Size          int    `json:"size"`
	SpaceSaved    int    `json:"spaceSaved"`
	OriginalSize  int    `json:"originalSize"`
	Path          string `json:"path"`
	Missing       bool   `json:"missing"`
}
