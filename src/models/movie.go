package models

type Movie struct {
	BaseModel
	Id           string `gorm:"primary_key"`
	Name         string `gorm:"type:varchar(255)"`
	ReleaseDate  string `gorm:"type:varchar(255)"`
	Genre        string `gorm:"type:varchar(255)"`
	Status       string `gorm:"type:varchar(255)"`
	Filename     string `gorm:"type:varchar(255)"`
	VideoCodec   string `gorm:"type:varchar(255)"`
	Overview     string `gorm:"type:text"`
	Size         int    `gorm:"type:int"`
	SpaceSaved   int    `gorm:"type:int"`
	ProfileID    int    `gorm:"type:int"`
	Monitored    bool   `gorm:"type:boolean"`
	Missing      bool   `gorm:"type:boolean"`
	Studio       string `gorm:"type:varchar(255)"`
	OriginalSize int    `gorm:"type:int"`
	Path         string `gorm:"type:varchar(255)"`
	Runtime      int    `gorm:"type:int"`
}
