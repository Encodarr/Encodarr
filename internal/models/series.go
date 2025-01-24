package models

type Series struct {
	BaseModel
	Id              string   `gorm:"primary_key" json:"id"`
	Name            string   `gorm:"type:varchar(255)" json:"name"`
	ReleaseDate     string   `gorm:"type:varchar(255)" json:"releaseDate"`
	Genre           string   `gorm:"type:varchar(255)" json:"genre"`
	Status          string   `gorm:"type:varchar(255)" json:"status"`
	LastAirDate     string   `gorm:"type:varchar(255)" json:"lastAirDate"`
	Networks        string   `gorm:"type:varchar(255)" json:"networks"`
	Overview        string   `gorm:"type:text" json:"overview"`
	ProfileID       int      `gorm:"type:int" json:"profileId"`
	Monitored       bool     `gorm:"type:boolean" json:"monitored"`
	EpisodeCount    int      `gorm:"type:int" json:"episodeCount"`
	Size            int      `gorm:"type:int" json:"size"`
	SeasonsCount    int      `gorm:"type:int" json:"seasonsCount"`
	SpaceSaved      int      `gorm:"type:int" json:"spaceSaved"`
	MissingEpisodes int      `gorm:"type:int" json:"missingEpisodes"`
	Runtime         int      `gorm:"type:int" json:"runtime"`
	Seasons         []Season `gorm:"foreignkey:SeriesId" json:"seasons"`
}
