package service

import (
	"github.com/third-place/user-service/internal/model"
)

type TestService struct {
	userService *UserService
}

func CreateTestService() *TestService {
	return &TestService{
		userService: CreateTestUserService(),
	}
}

func (t *TestService) CreateInvitedUser(user *model.NewUser) (*model.User, error) {
	inviteCode, _ := t.userService.CreateInvite()
	user.InviteCode = inviteCode.Code
	user.Name = "foo"
	return t.userService.CreateUser(user)
}

func (t *TestService) UpdateUser(session *model.Session, user *model.User) error {
	return t.userService.UpdateUser(session, user)
}
