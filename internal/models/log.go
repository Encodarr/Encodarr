package models

type Log struct {
	ID        uint   `gorm:"primary_key"`
	Timestamp string `gorm:"type:datetime"`
	Level     string `gorm:"type:varchar(10)"`
	Service   string `gorm:"type:varchar(50)"`
	Message   string `gorm:"type:text"`
}
