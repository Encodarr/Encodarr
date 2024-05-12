package models

type Season struct {
	BaseModel
	Id              string `gorm:"primary_key"`
	Name            string `gorm:"type:varchar(255)"`
	SeasonNumber    int    `gorm:"type:int"`
	EpisodeCount    int    `gorm:"type:int"`
	Size            int    `gorm:"type:int"`
	SeriesId        string `gorm:"type:varchar(255)"`
	SpaceSaved      string `gorm:"type:varchar(255)"`
	MissingEpisodes int    `gorm:"type:int"`
}
