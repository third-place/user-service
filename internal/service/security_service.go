package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/third-place/user-service/internal/db"
	"github.com/third-place/user-service/internal/entity"
	"github.com/third-place/user-service/internal/model"
	"github.com/third-place/user-service/internal/repository"
)

type SecurityService struct {
	userRepository *repository.UserRepository
}

func CreateSecurityService() *SecurityService {
	conn := db.CreateDefaultConnection()
	return &SecurityService{
		repository.CreateUserRepository(conn),
	}
}

func (s *SecurityService) IsInGoodStanding(session *model.Session) bool {
	user, err := s.getUser(session)
	if err != nil {
		return false
	}
	return !user.IsBanned
}

func (s *SecurityService) IsModerator(session *model.Session) bool {
	user, err := s.getUser(session)
	if err != nil {
		return false
	}
	return !user.IsBanned && (user.Role == string(model.MODERATOR) || user.Role == string(model.ADMIN))
}

func (s *SecurityService) getUser(session *model.Session) (*entity.User, error) {
	if session.User == nil {
		return nil, errors.New("user not defined")
	}
	userUuid, err := uuid.Parse(session.User.Uuid)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepository.GetUserFromUuid(userUuid)
	if err != nil {
		return nil, err
	}
	return user, nil
}
