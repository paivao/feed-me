package model

import (
	"gorm.io/gorm"
)

type Feed struct {
	gorm.Model
	Name string
}
