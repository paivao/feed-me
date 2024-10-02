package model

import "time"

type Entry struct {
	Id         uint64
	Enabled    bool      `gorm:"not null"`
	ValidUntil time.Time `gorm:"not null"`
}
