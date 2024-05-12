package models

type Series struct {
	BaseModel
	Id              string `gorm:"primary_key"`
	Name            string `gorm:"type:varchar(255)"`
	ReleaseDate     string `gorm:"type:varchar(255)"`
	Genre           string `gorm:"type:varchar(255)"`
	Status          string `gorm:"type:varchar(255)"`
	LastAirDate     string `gorm:"type:varchar(255)"`
	Networks        string `gorm:"type:varchar(255)"`
	Overview        string `gorm:"type:text"`
	ProfileID       int    `gorm:"type:int"`
	Monitored       bool   `gorm:"type:boolean"`
	EpisodeCount    int    `gorm:"type:int"`
	Size            int    `gorm:"type:int"`
	SeasonsCount    int    `gorm:"type:int"`
	SpaceSaved      int    `gorm:"type:int"`
	MissingEpisodes int    `gorm:"type:int"`
	Runtime         int    `gorm:"type:int"`
}
