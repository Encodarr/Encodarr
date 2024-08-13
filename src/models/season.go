package models

type Season struct {
	BaseModel
	Id              string    `gorm:"primary_key" json:"id"`
	Name            string    `gorm:"type:varchar(255)" json:"name"`
	SeasonNumber    int       `gorm:"type:int" json:"seasonNumber"`
	EpisodeCount    int       `gorm:"type:int" json:"episodeCount"`
	Size            int       `gorm:"type:int" json:"size"`
	SeriesId        string    `gorm:"type:varchar(255)" json:"seriesId"`
	SpaceSaved      int       `gorm:"type:varchar(255)" json:"spaceSaved"`
	MissingEpisodes int       `gorm:"type:int" json:"missingEpisodes"`
	Episodes        []Episode `gorm:"foreignkey:SeriesId" json:"episodes"`
}
