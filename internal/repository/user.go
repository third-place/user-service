package repository

import (
	"errors"
	"github.com/danielmunro/otto-user-service/internal/entity"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	conn *gorm.DB
}

func CreateUserRepository(conn *gorm.DB) *UserRepository {
	return &UserRepository{conn}
}

func (r *UserRepository) GetUserFromUsername(username string) (*entity.User, error) {
	user := &entity.User{}
	r.conn.Where("username = ?", username).Find(&user)
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) GetUserFromUuid(uuid uuid.UUID) (*entity.User, error) {
	user := &entity.User{}
	r.conn.Where("uuid = ?", uuid.String()).Find(&user)
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) GetUserFromEmail(email string) (*entity.User, error) {
	user := &entity.User{}
	r.conn.Where("current_email = ?", email).Find(&user)
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) GetUserFromSessionToken(token string) (*entity.User, error) {
	user := &entity.User{}
	r.conn.Where("last_access_token = ?", token).Find(&user)
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) Create(user *entity.User) *gorm.DB {
	return r.conn.Create(user)
}

func (r *UserRepository) Delete(user *entity.User) *gorm.DB {
	return r.conn.Unscoped().Delete(user)
}

func (r *UserRepository) Save(user *entity.User) *gorm.DB {
	return r.conn.Save(user)
}
