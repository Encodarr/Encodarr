package models

type User struct {
	BaseModel
	Username string `gorm:"type:varchar(255)"`
	Password string `gorm:"type:varchar(255)"`
	Secret   string `gorm:"type:varchar(255)"`
}
