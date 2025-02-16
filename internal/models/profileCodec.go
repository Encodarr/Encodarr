package models

type ProfileCodec struct {
	Id        int    `json:"id"`
	ProfileId int    `json:"profileId"`
	CodecId   string `json:"codecId"`
}
