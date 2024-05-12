package models

type System struct {
	BaseModel
	Id    string `gorm:"primary_key"`
	Value string `gorm:"type:varchar(255)"`
}
