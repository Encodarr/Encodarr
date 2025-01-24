package models

type Event struct {
	BaseModel
	Id        int    `gorm:"primary_key" json:"id"`
	Timestamp string `gorm:"type:varchar(255)" json:"timestamp"`
	Level     string `gorm:"type:varchar(255)" json:"level"`
	Service   string `gorm:"type:varchar(255)" json:"service"`
	Message   string `gorm:"type:varchar(255)" json:"message"`
}
