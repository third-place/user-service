package model

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

type Claims struct {
	UserUuid string `json:"userUuid"`
	jwt.RegisteredClaims
}

func NewClaims(userUuid uuid.UUID) *Claims {
	expirationTime := time.Now().Add(24 * 7 * time.Hour)
	return &Claims{
		UserUuid: userUuid.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
}
