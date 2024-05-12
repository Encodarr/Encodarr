package models

type ProfileAudioLanguage struct {
	BaseModel
	Id        int    `gorm:"primary_key"`
	ProfileId string `gorm:"type:varchar(255)"`
	Language  string `gorm:"type:varchar(255)"`
}
