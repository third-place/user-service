package service

import (
	"github.com/third-place/user-service/internal/entity"
	"github.com/third-place/user-service/internal/mapper"
)

func createSessionResponse(user *entity.User) *AuthResponse {
	return &AuthResponse{
		Token: &user.JWT,
		User:  mapper.MapUserEntityToUser(user),
	}
}
