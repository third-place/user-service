package service

import (
	"github.com/third-place/user-service/internal/entity"
	"github.com/third-place/user-service/internal/mapper"
)

func createSessionResponse(user *entity.User, token string) *AuthResponse {
	return &AuthResponse{
		Token: token,
		User:  mapper.MapUserEntityToUser(user),
	}
}
