package entity

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Password struct {
	gorm.Model
	Uuid uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID uint
	Password string `gorm:"not null"`
}
