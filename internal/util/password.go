package util

import "unicode"

func ValidatePassword(password string) (minSize, digit, special, lowercase, uppercase bool) {
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			digit = true
		case unicode.IsUpper(c):
			uppercase = true
		case unicode.IsLower(c):
			lowercase = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		}
	}
	minSize = len(password) >= 8
	return
}
