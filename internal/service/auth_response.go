package service

import (
	"github.com/third-place/user-service/internal/model"
)

type AuthResponseType int

const (
	Unknown                     AuthResponseType = iota + 1
	ChallengeNewPassword        AuthResponseType = iota
	SessionAuthenticated        AuthResponseType = iota
	SessionFailedAuthentication AuthResponseType = iota
)

type AuthResponse struct {
	AuthResponse AuthResponseType
	Token        string
	User         *model.User
	Message      string
}
