package model

type IPFeed struct {
	Feed
	Entries []IPEntry `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
