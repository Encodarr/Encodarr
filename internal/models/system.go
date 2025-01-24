package models

type System struct {
	BaseModel
	Id    string `gorm:"primary_key" json:"id"`
	Value string `gorm:"type:varchar(255)" json:"value"`
}
