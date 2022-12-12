package service

import (
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/third-place/user-service/internal/model"
	"github.com/third-place/user-service/internal/util"
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const dummyPassword = "fOobar12345!"

func TestMain(m *testing.M) {
	if os.Getenv("CI") == "" {
		_ = godotenv.Load()
	}
	os.Exit(m.Run())
}

func Test_CreateNewUser_SanityCheck(t *testing.T) {
	// setup
	svc := CreateTestService()

	// when
	user, err := svc.CreateInvitedUser(&model.NewUser{
		Username: util.RandomUsername(),
		Email:    util.RandomEmailAddress(),
		Password: dummyPassword,
	})

	// then
	if user == nil || err != nil {
		t.Error(err)
	}
}

func Test_GetUserByUuid(t *testing.T) {
	// setup
	svc := CreateTestService()

	// given
	user, err := svc.CreateInvitedUser(&model.NewUser{
		Username: util.RandomUsername(),
		Email:    util.RandomEmailAddress(),
		Password: dummyPassword,
	})

	// when
	getUser, err := svc.GetUserFromUuid(uuid.MustParse(user.Uuid))

	// then
	if getUser == nil || err != nil {
		t.Error(err)
	}
}

func Test_GetUserByUsername(t *testing.T) {
	// setup
	svc := CreateTestService()

	// given
	user, err := svc.CreateInvitedUser(&model.NewUser{
		Username: util.RandomUsername(),
		Email:    util.RandomEmailAddress(),
		Password: dummyPassword,
	})

	// when
	getUser, err := svc.GetUserFromUsername(user.Username)

	// then
	if getUser == nil || err != nil {
		t.Error(err)
	}
}

func Test_Can_Login(t *testing.T) {
	// setup
	svc := CreateTestService()

	// given
	user, err := svc.CreateInvitedUser(&model.NewUser{
		Username: util.RandomUsername(),
		Email:    util.RandomEmailAddress(),
		Password: dummyPassword,
	})

	// when
	session, err := svc.CreateSession(&model.NewSession{
		Email:    user.Email,
		Password: dummyPassword,
	})

	// then
	if session == nil || err != nil {
		t.Error(err)
	}
}

func Test_UserCan_UpdateSelf(t *testing.T) {
	// setup
	svc := CreateTestService()

	// given
	user, _ := svc.CreateInvitedUser(&model.NewUser{
		Username: util.RandomUsername(),
		Email:    util.RandomEmailAddress(),
		Password: dummyPassword,
	})
	_ = svc.UpdateUser(
		&model.Session{
			User: user,
		},
		&model.User{
			Uuid:       user.Uuid,
			Name:       "MyName",
			ProfilePic: "MyProfilePic",
			BioMessage: "Hello World",
			Birthday:   "2000-01-01",
		},
	)

	// when
	user, _ = svc.GetUserFromUuid(uuid.MustParse(user.Uuid))

	// then
	if user.Name != "MyName" {
		t.Error("expected to update user")
	}
	if user.ProfilePic != "MyProfilePic" {
		t.Error("expected to update profile pic")
	}
	if user.BioMessage != "Hello World" {
		t.Error("expected to update bio message")
	}
	if user.Birthday != "2000-01-01" {
		t.Error("expected to update birthday")
	}
}

func Test_CannotUpdate_OtherUsers(t *testing.T) {
	// setup
	svc := CreateTestService()

	// given
	user1, err := svc.CreateInvitedUser(&model.NewUser{
		Username: util.RandomUsername(),
		Email:    util.RandomEmailAddress(),
		Password: dummyPassword,
	})
	user2, err := svc.CreateInvitedUser(&model.NewUser{
		Username: util.RandomUsername(),
		Email:    util.RandomEmailAddress(),
		Password: dummyPassword,
	})

	// when
	err = svc.UpdateUser(
		&model.Session{
			User: user1,
		},
		user2,
	)

	// then
	if err == nil || err.Error() != "unauthorized" {
		t.Error("expected error when one user updates another user")
	}
}

func Test_Needs_Valid_Invite(t *testing.T) {
	// setup
	svc := CreateTestService()

	// when
	_, err := svc.CreateUser(
		&model.Invite{
			Code: "this-does-not-exist",
		},
		&model.NewUser{
			Username: util.RandomUsername(),
			Email:    util.RandomEmailAddress(),
			Password: dummyPassword,
		},
	)

	// then
	if err == nil {
		t.Error("expected valid invite code")
	}
}

func Test_Cannot_Reuse_Invite(t *testing.T) {
	// setup
	svc := CreateTestService()

	// given
	invite, _ := svc.CreateInvite()

	// when
	_, _ = svc.CreateUser(
		invite,
		&model.NewUser{
			Username: util.RandomUsername(),
			Email:    util.RandomEmailAddress(),
			Password: dummyPassword,
		},
	)
	_, err := svc.CreateUser(
		invite,
		&model.NewUser{
			Username: util.RandomUsername(),
			Email:    util.RandomEmailAddress(),
			Password: dummyPassword,
		},
	)

	// then
	if err == nil {
		t.Error("expected failure when re-using invite")
	}
}

func Test_Email_Uniqueness(t *testing.T) {
	// setup
	svc := CreateTestService()

	// given
	email := util.RandomEmailAddress()

	// when
	_, _ = svc.CreateInvitedUser(&model.NewUser{
		Username: util.RandomUsername(),
		Email:    email,
		Password: dummyPassword,
	})
	_, err := svc.CreateInvitedUser(&model.NewUser{
		Username: util.RandomUsername(),
		Email:    email,
		Password: dummyPassword,
	})

	// then
	if err == nil {
		t.Error("expected duplicate email")
	}
}

func Test_Username_Uniqueness(t *testing.T) {
	// setup
	svc := CreateTestService()

	// given
	username := util.RandomUsername()

	// when
	_, _ = svc.CreateInvitedUser(&model.NewUser{
		Username: username,
		Email:    util.RandomEmailAddress(),
		Password: dummyPassword,
	})
	_, err := svc.CreateInvitedUser(&model.NewUser{
		Username: username,
		Email:    util.RandomEmailAddress(),
		Password: dummyPassword,
	})

	// then
	if err == nil {
		t.Error("expected duplicate email")
	}
}

func Test_Password_Length(t *testing.T) {
	// setup
	svc := CreateTestService()

	// when
	_, err := svc.CreateInvitedUser(&model.NewUser{
		Username: util.RandomUsername(),
		Email:    util.RandomEmailAddress(),
		Password: "foo",
	})

	// then
	if err == nil {
		t.Error("expected error")
	}
	inputErr := err.(*util.InputFieldError)
	if inputErr.Input != "password" {
		t.Error("input error expected to be password")
	}
}

func Test_Password_Complexity(t *testing.T) {
	// setup
	svc := CreateTestService()

	// when
	_, err := svc.CreateInvitedUser(&model.NewUser{
		Username: util.RandomUsername(),
		Email:    util.RandomEmailAddress(),
		Password: "fooooooo",
	})

	// then
	if err == nil {
		t.Error("expected error")
	}
	inputErr := err.(*util.InputFieldError)
	if inputErr.Input != "password" {
		t.Error("input error expected to be password")
	}
}
