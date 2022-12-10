package service

import (
	"github.com/joho/godotenv"
	"github.com/third-place/user-service/internal/model"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

const dummyPassword = "fOobar12345!"

func TestMain(m *testing.M) {
	if os.Getenv("CI") == "" {
		_ = godotenv.Load()
	}
	os.Exit(m.Run())
}

func GetEmailAddress() string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	num := r.Intn(100000)
	return os.Getenv("EMAIL_PREFIX") + "+" + strconv.Itoa(num) + "@" + os.Getenv("EMAIL_DOMAIN")
}

func Test_CreateNewUser_SanityCheck(t *testing.T) {
	user, err := CreateUserService().CreateUser(&model.NewUser{
		Name:     "foo",
		Email:    GetEmailAddress(),
		Password: dummyPassword,
	})

	if user == nil || err != nil {
		t.Error(err)
	}
}

func Test_RegisteringWith_DuplicateEmails_Error(t *testing.T) {
	svc := CreateUserService()
	userModel := &model.NewUser{
		Name:     "foo",
		Email:    GetEmailAddress(),
		Password: dummyPassword,
	}
	_, _ = svc.CreateUser(userModel)
	_, err := svc.CreateUser(userModel)

	if err == nil {
		t.Error("expected duplicate email")
	}
}

func Test_CreateFirstSession_WillReceiveChallenge(t *testing.T) {
	svc := CreateUserService()
	email := GetEmailAddress()
	_, _ = svc.CreateUser(&model.NewUser{
		Email:    email,
		Password: dummyPassword,
	})
	response, _ := svc.CreateSession(&model.NewSession{
		Email:    email,
		Password: dummyPassword,
	})
	if response == nil || response.AuthResponse != ChallengeNewPassword {
		t.Error("expected challenge")
	}
}

func Test_AuthFlow_FromStart_ToVerifiedUser(t *testing.T) {
	svc := CreateUserService()
	email := GetEmailAddress()
	_, _ = svc.CreateUser(&model.NewUser{
		Email:    email,
		Password: dummyPassword,
	})
	svc.CreateSession(&model.NewSession{
		Email:    email,
		Password: dummyPassword,
	})
	response := svc.ProvideChallengeResponse(&model.PasswordReset{
		Email:    email,
		Password: "my-awesome-new-pAssword-123!",
	})
	if response == nil || response.AuthResponse != SessionAuthenticated {
		t.Error("authflow")
	}
}
