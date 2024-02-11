package entity

import "gorm.io/gorm"

type Invite struct {
	gorm.Model
	Code    string `gorm:"unique;not null"`
	Claimed bool
}
