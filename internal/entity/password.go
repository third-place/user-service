package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Password struct {
	gorm.Model
	Uuid     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID   uint
	Password string `gorm:"not null"`
}
