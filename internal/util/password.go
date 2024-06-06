package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
)

var StaticSalt = os.Getenv("STATIC_SALT")

func formatPassword(userID uint, password string) []byte {
	return []byte(fmt.Sprintf("%d-%s-%s", userID, StaticSalt, password))
}

func ValidatePassword(password string) (minSize bool) {
	minSize = len(password) >= 8
	return
}

func HashPassword(userID uint, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(formatPassword(userID, password), 14)
	return string(bytes), err
}

func CheckPasswordHash(userID uint, password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), formatPassword(userID, password))
	return err == nil
}
