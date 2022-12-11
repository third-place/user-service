package test

import (
	"github.com/joho/godotenv"
	"github.com/third-place/user-service/internal/model"
	"github.com/third-place/user-service/internal/service"
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
	svc := service.CreateTestService()

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

func Test_CannotUpdate_OtherUsers(t *testing.T) {
	// setup
	svc := service.CreateTestService()

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

func Test_Email_Uniqueness(t *testing.T) {
	// setup
	svc := service.CreateTestService()

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
	svc := service.CreateTestService()

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
	svc := service.CreateTestService()

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
	svc := service.CreateTestService()

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
