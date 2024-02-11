package entity

import (
	"github.com/google/uuid"
	"github.com/third-place/user-service/internal/enum"
	"gorm.io/gorm"
)

type Email struct {
	gorm.Model
	Uuid           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID         uint
	Email          string `gorm:"unique;not null"`
	VerifiedStatus string `gorm:"not null"`
}

func CreateEmail(email string) *Email {
	return &Email{
		Email:          email,
		VerifiedStatus: string(enum.EmailStatusUnverified),
	}
}
