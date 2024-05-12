package models

type Episode struct {
	BaseModel
	Id            string `gorm:"type:varchar(255)"`
	SeriesId      string `gorm:"type:varchar(255)"`
	SeasonId      string `gorm:"type:varchar(255)"`
	EpisodeNumber int    `gorm:"type:int"`
	SeasonName    string `gorm:"type:varchar(255)"`
	SeasonNumber  int    `gorm:"type:int"`
	Filename      string `gorm:"type:varchar(255)"`
	EpisodeName   string `gorm:"type:varchar(255)"`
	VideoCodec    string `gorm:"type:varchar(255)"`
	AirDate       string `gorm:"type:varchar(255)"`
	Size          int    `gorm:"type:int"`
	SpaceSaved    int    `gorm:"type:int"`
	OriginalSize  int    `gorm:"type:int"`
	Path          string `gorm:"type:varchar(255)"`
}
