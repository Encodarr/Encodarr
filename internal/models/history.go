package models

type History struct {
	Id            int    `json:"id"`
	MediaId       string `json:"mediaId"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	SeasonNumber  int    `json:"seasonNumber"`
	EpisodeNumber int    `json:"episodeNumber"`
	ProfileId     int    `json:"profileId"`
	PrevCodec     string `json:"prevCodec"`
	NewCodec      string `json:"newCodec"`
	PrevSize      int    `json:"prevSize"`
	NewSize       int    `json:"newSize"`
	Date          string `json:"date"`
}
