package service

import (
	"github.com/google/uuid"
	"github.com/third-place/user-service/internal/entity"
	"github.com/third-place/user-service/internal/mapper"
	"github.com/third-place/user-service/internal/model"
	"github.com/third-place/user-service/internal/repository"
	"github.com/third-place/user-service/internal/util"
)

type TestService struct {
	userService      *UserService
	inviteRepository *repository.InviteRepository
}

func CreateTestService() *TestService {
	conn := util.SetupTestDatabase()
	return &TestService{
		userService:      CreateTestUserService(),
		inviteRepository: repository.CreateInviteRepository(conn),
	}
}

func (t *TestService) CreateInvitedUser(user *model.NewUser) (*model.User, error) {
	inviteCode, _ := t.createInvite()
	user.InviteCode = inviteCode.Code
	user.Name = "foo"
	return t.userService.CreateUser(user)
}

func (t *TestService) CreateInvite() (*model.Invite, error) {
	return t.createInvite()
}

func (t *TestService) CreateUser(inviteCode *model.Invite, user *model.NewUser) (*model.User, error) {
	user.InviteCode = inviteCode.Code
	user.Name = "foo"
	return t.userService.CreateUser(user)
}

func (t *TestService) UpdateUser(session *model.Session, user *model.User) error {
	return t.userService.UpdateUser(session, user)
}

func (t *TestService) GetUserFromUuid(uuid uuid.UUID) (*model.User, error) {
	return t.userService.GetUserFromUuid(uuid)
}

func (t *TestService) GetUserFromUsername(username string) (*model.User, error) {
	return t.userService.GetUserFromUsername(username)
}

func (t *TestService) CreateSession(newSession *model.NewSession) (*model.Session, error) {
	return t.userService.CreateSession(newSession)
}

func (t *TestService) GetSession(sessionToken *model.SessionToken) (*model.Session, error) {
	return t.userService.GetSession(sessionToken)
}

func (t *TestService) createInvite() (*model.Invite, error) {
	invite := &entity.Invite{
		Code: util.GenerateCode(),
	}
	result := t.inviteRepository.Create(invite)
	if result.Error != nil {
		return nil, result.Error
	}
	return mapper.MapInviteEntityToModel(invite), nil
}
