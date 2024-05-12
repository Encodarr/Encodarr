package models

type ProfileCodec struct {
	BaseModel
	Id        int    `gorm:"primary_key;AUTO_INCREMENT"`
	ProfileId int    `gorm:"type:int"`
	CodecId   string `gorm:"type:varchar(255)"`
}
