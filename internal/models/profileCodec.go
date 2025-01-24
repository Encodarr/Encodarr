package models

type ProfileCodec struct {
	BaseModel
	Id        int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	ProfileId int    `gorm:"type:int" json:"profileId"`
	CodecId   string `gorm:"type:varchar(255)" json:"codecId"`
}
