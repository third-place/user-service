package mapper

import (
	"github.com/third-place/user-service/internal/entity"
	"github.com/third-place/user-service/internal/model"
)

func MapUserEntityToModel(user *entity.User) *model.User {
	return &model.User{
		Uuid:       user.Uuid.String(),
		Name:       user.Name,
		Username:   user.Username,
		ProfilePic: user.ProfilePic,
		Role:       model.Role(user.Role),
		IsBanned:   user.IsBanned,
		BioMessage: user.BioMessage,
		Birthday:   user.Birthday,
		CreatedAt:  user.CreatedAt,
	}
}

func MapUserEntitiesToModels(users []*entity.User) []*model.User {
	userModels := make([]*model.User, len(users))
	for i, v := range users {
		userModels[i] = MapUserEntityToModel(v)
	}
	return userModels
}

func MapNewUserModelToEntity(user *model.NewUser) *entity.User {
	return &entity.User{
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
	}
}
