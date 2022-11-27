package entity

import "github.com/jinzhu/gorm"

type Invite struct {
	gorm.Model
	Code    string `gorm:"unique;not null"`
	Claimed bool
}
