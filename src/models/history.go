package models

type History struct {
	BaseModel
	Id            int    `gorm:"primary_key"`
	MediaId       string `gorm:"type:varchar(255)"`
	Name          string `gorm:"type:varchar(255)"`
	Type          string `gorm:"type:varchar(255)"`
	SeasonNumber  int    `gorm:"type:int"`
	EpisodeNumber int    `gorm:"type:int"`
	ProfileId     string `gorm:"type:varchar(255)"`
	PrevCodec     string `gorm:"type:varchar(255)"`
	NewCodec      string `gorm:"type:varchar(255)"`
	PrevSize      string `gorm:"type:varchar(255)"`
	NewSize       string `gorm:"type:varchar(255)"`
	Date          string `gorm:"type:varchar(255)"`
}
