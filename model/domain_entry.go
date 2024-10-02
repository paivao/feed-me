package model

type DomainEntry struct {
	Entry
	Domain string `gorm:""`
}
