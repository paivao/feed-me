package controller

import "gorm.io/gorm"

type baseController struct {
	db *gorm.DB
}
