package model

import "time"

type FeedType uint8

const (
	IPFeed FeedType = iota
	DomainFeed
	URLFeed
)

type Feed struct {
	Id        uint64
	Name      string
	Type      FeedType
	Entries   []Entry
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
