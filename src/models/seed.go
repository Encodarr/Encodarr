package models

import "gorm.io/gorm"

type Seed struct {
	gorm.Model
	Name string `gorm:"unique"`
}
