package models

type User struct {
	BaseModel
	Username string `gorm:"type:varchar(255)" json:"username"`
	Password string `gorm:"type:varchar(255)" json:"password"`
	Secret   string `gorm:"type:varchar(255)" json:"secret"`
}
