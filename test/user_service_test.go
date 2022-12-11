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

const dummyPassword = "fOobar12345!"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestMain(m *testing.M) {
	if os.Getenv("CI") == "" {
		_ = godotenv.Load()
	}
	os.Exit(m.Run())
}

func Test_CreateNewUser_SanityCheck(t *testing.T) {
	svc := service.CreateTestUserService()
	code, _ := svc.CreateInvite()
	user, err := svc.CreateUser(&model.NewUser{
		Name:       "foo",
		Username:   util.RandomUsername(),
		Email:      util.RandomEmailAddress(),
		Password:   dummyPassword,
		InviteCode: code.Code,
	})

	if user == nil || err != nil {
		t.Error(err)
	}
}

func Test_RegisteringWith_DuplicateEmails_Error(t *testing.T) {
	svc := service.CreateTestUserService()
	code1, _ := svc.CreateInvite()
	code2, _ := svc.CreateInvite()
	email := util.RandomEmailAddress()
	userModel1 := &model.NewUser{
		Name:       "foo",
		Username:   util.RandomUsername(),
		Email:      email,
		Password:   dummyPassword,
		InviteCode: code1.Code,
	}
	userModel2 := &model.NewUser{
		Name:       "foo",
		Username:   util.RandomUsername(),
		Email:      email,
		Password:   dummyPassword,
		InviteCode: code2.Code,
	}
	_, _ = svc.CreateUser(userModel1)
	_, err := svc.CreateUser(userModel2)

	if err == nil {
		t.Error("expected duplicate email")
	}
}
