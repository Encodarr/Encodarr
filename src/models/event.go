package models

type Event struct {
	BaseModel
	Id        int    `gorm:"primary_key"`
	Timestamp string `gorm:"type:varchar(255)"`
	Level     string `gorm:"type:varchar(255)"`
	Service   string `gorm:"type:varchar(255)"`
	Message   string `gorm:"type:varchar(255)"`
}
