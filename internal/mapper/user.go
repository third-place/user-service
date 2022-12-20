package mapper

import (
	"github.com/third-place/user-service/internal/entity"
	"github.com/third-place/user-service/internal/model"
	"github.com/third-place/user-service/internal/util"
)

func MapUserEntityToModel(user *entity.User) *model.User {
	return &model.User{
		Uuid:          user.Uuid.String(),
		Name:          user.Name,
		Username:      user.Username,
		Email:         user.Email,
		ProfilePic:    user.ProfilePic,
		Role:          model.Role(user.Role),
		IsBanned:      user.IsBanned,
		AddressCity:   user.AddressCity,
		AddressStreet: user.AddressStreet,
		AddressZip:    user.AddressZip,
		BioMessage:    user.BioMessage,
		Birthday:      user.Birthday,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}

func MapUserEntityToUser(user *entity.User) *model.User {
	return &model.User{
		Uuid:          user.Uuid.String(),
		Name:          user.Name,
		Username:      user.Username,
		ProfilePic:    user.ProfilePic,
		IsBanned:      user.IsBanned,
		Role:          model.Role(user.Role),
		AddressCity:   user.AddressCity,
		AddressStreet: user.AddressStreet,
		AddressZip:    user.AddressZip,
		BioMessage:    user.BioMessage,
		Birthday:      user.Birthday,
		CreatedAt:     user.CreatedAt,
	}
}

func MapNewUserModelToEntity(user *model.NewUser) *entity.User {
	pass, _ := util.HashPassword(user.Password)
	return &entity.User{
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Password: pass,
	}
}
