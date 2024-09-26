package model

import "time"

type Entry struct {
	Id         uint64
	Value      string
	Enabled    bool
	ValidUntil time.Time
}
