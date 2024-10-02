package model

import "time"

type Entry struct {
	Id         uint64
	Enabled    bool
	ValidUntil time.Time
}
