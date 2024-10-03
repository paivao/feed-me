package model

import (
	"gorm.io/gorm"
)

type Feed struct {
	gorm.Model
	Name     string `gorm:"size:128;not null;unique"`
	IsPublic bool   `gorm:"not null"`
}
