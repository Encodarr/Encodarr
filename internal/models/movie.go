package models

type Movie struct {
	BaseModel
	Id           string `gorm:"primary_key" json:"id"`
	Name         string `gorm:"type:varchar(255)" json:"name"`
	ReleaseDate  string `gorm:"type:varchar(255)" json:"releaseDate"`
	Genre        string `gorm:"type:varchar(255)" json:"genre"`
	Status       string `gorm:"type:varchar(255)" json:"status"`
	Filename     string `gorm:"type:varchar(255)" json:"filename"`
	VideoCodec   string `gorm:"type:varchar(255)" json:"videoCodec"`
	Overview     string `gorm:"type:text" json:"overview"`
	Size         int    `gorm:"type:int" json:"size"`
	SpaceSaved   int    `gorm:"type:int" json:"spaceSaved"`
	ProfileID    int    `gorm:"type:int" json:"profileId"`
	Monitored    bool   `gorm:"type:boolean;default:false" json:"monitored"`
	Missing      bool   `gorm:"type:boolean" json:"missing"`
	Studio       string `gorm:"type:varchar(255)" json:"studio"`
	OriginalSize int    `gorm:"type:int" json:"originalSize"`
	Path         string `gorm:"type:varchar(255)" json:"path"`
	Runtime      int    `gorm:"type:int" json:"runtime"`
}
