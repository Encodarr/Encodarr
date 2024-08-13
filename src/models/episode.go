package models

type Episode struct {
	BaseModel
	Id            string `gorm:"type:varchar(255)" json:"id"`
	SeriesId      string `gorm:"type:varchar(255)" json:"seriesId"`
	SeasonId      string `gorm:"type:varchar(255)" json:"seasonId"`
	EpisodeNumber int    `gorm:"type:int" json:"episodeNumber"`
	SeasonName    string `gorm:"type:varchar(255)" json:"seasonName"`
	SeasonNumber  int    `gorm:"type:int" json:"seasonNumber"`
	Filename      string `gorm:"type:varchar(255)" json:"filename"`
	EpisodeName   string `gorm:"type:varchar(255)" json:"episodeName"`
	VideoCodec    string `gorm:"type:varchar(255)" json:"videoCodec"`
	AirDate       string `gorm:"type:varchar(255)" json:"airDate"`
	Size          int    `gorm:"type:int" json:"size"`
	SpaceSaved    int    `gorm:"type:int" json:"spaceSaved"`
	OriginalSize  int    `gorm:"type:int" json:"originalSize"`
	Path          string `gorm:"type:varchar(255)" json:"path"`
	Missing       bool   `gorm:"type:boolean" json:"missing"`
}
