package models

type ProfileAudioLanguage struct {
	BaseModel
	Id        int    `gorm:"primary_key" json:"id"`
	ProfileId int    `gorm:"type:int" json:"profileId"`
	Language  string `gorm:"type:varchar(255)" json:"language"`
}
