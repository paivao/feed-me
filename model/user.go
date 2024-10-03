package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string       `gorm:"size:128;not null;unique"`
	Email        string       `gorm:"size:128;not null;unique"`
	PasswordHash string       `gorm:"size:255;not null"`
	Permission   []Permission `gorm:"many2many:user_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Groups       []Group      `gorm:"many2many:user_groups;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Group struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"size:128;not null;unique"`
	Members []User `gorm:"many2many:user_groups;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Permission struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255;not null;unique"`
}
