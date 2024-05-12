package models

import "github.com/jinzhu/gorm"

type Seed struct {
	gorm.Model
	Name string `gorm:"unique"`
}
