package models

type ProfileSubtitleLanguage struct {
	Id        int    `json:"id"`
	ProfileId int    `json:"profileId"`
	Language  string `json:"language"`
}
