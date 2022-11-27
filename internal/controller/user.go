package controller

import (
	"encoding/json"
	"github.com/danielmunro/otto-user-service/internal/db"
	"github.com/danielmunro/otto-user-service/internal/model"
	"github.com/danielmunro/otto-user-service/internal/repository"
	"github.com/danielmunro/otto-user-service/internal/service"
	"github.com/danielmunro/otto-user-service/internal/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// CreateNewUserV1 - Create a new user
func CreateNewUserV1(w http.ResponseWriter, r *http.Request) {
	newUserModel := model.DecodeRequestToNewUser(r)
	user, err := service.CreateDefaultUserService().CreateUser(newUserModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, ok := err.(*util.InputFieldError); ok {
			data, _ := json.Marshal(err)
			_, _ = w.Write(data)
			return
		}
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(user)
	_, _ = w.Write(data)
}

// GetUserByUsernameV1 - Get a user by username
func GetUserByUsernameV1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "max-age=30")
	params := mux.Vars(r)
	username := params["username"]

	user, err := service.CreateDefaultUserService().GetUserFromUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	_, _ = w.Write(data)

}

// UpdateUserV1 - Update a user
func UpdateUserV1(w http.ResponseWriter, r *http.Request) {
	userModel, err := model.DecodeRequestToUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userService := service.CreateDefaultUserService()
	sessionToken := getSessionToken(r)
	sessionModel := &model.SessionToken{
		Token: sessionToken,
	}
	session, err := userService.GetSession(sessionModel)
	if err != nil || session.User.Uuid != userModel.Uuid {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	err = userService.UpdateUser(userModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(userModel)
	_, _ = w.Write(data)
}

// BanUserV1 - ban a user
func BanUserV1(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	usernameParam := params["username"]
	userService := service.CreateDefaultUserService()
	userRepository := repository.CreateUserRepository(db.CreateDefaultConnection())
	sessionToken := getSessionToken(r)
	sessionModel := &model.SessionToken{
		Token: sessionToken,
	}
	session, err := userService.GetSession(sessionModel)
	if err != nil {
		log.Print("error 0 :: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionUser, err := userRepository.GetUserFromUsername(session.User.Username)
	if err != nil || sessionUser.IsBanned {
		log.Print("error 1 :: ", err.Error())
		log.Print("sessionUser isBanned :: ", sessionUser.IsBanned)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userEntity, err := userRepository.GetUserFromUsername(usernameParam)
	if err != nil {
		log.Print("error 2 :: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = userService.BanUser(sessionUser, userEntity)
	if err != nil {
		log.Print("error 3 :: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}
}

// UnbanUserV1 - ban a user
func UnbanUserV1(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	usernameParam := params["username"]
	userService := service.CreateDefaultUserService()
	userRepository := repository.CreateUserRepository(db.CreateDefaultConnection())
	sessionToken := getSessionToken(r)
	sessionModel := &model.SessionToken{
		Token: sessionToken,
	}
	session, err := userService.GetSession(sessionModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionUser, err := userRepository.GetUserFromUsername(session.User.Username)
	if err != nil || sessionUser.IsBanned {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userEntity, err := userRepository.GetUserFromUsername(usernameParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = userService.UnbanUser(sessionUser, userEntity)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// SubmitOTPV1 - Submit a new OTP
func SubmitOTPV1(w http.ResponseWriter, r *http.Request) {
	otpModel := model.DecodeRequestToOtp(r)
	userService := service.CreateDefaultUserService()
	err := userService.SubmitOTP(otpModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// SubmitForgotPasswordV1 - Submit a forgot password request
func SubmitForgotPasswordV1(w http.ResponseWriter, r *http.Request) {
	userModel, err := model.DecodeRequestToUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userService := service.CreateDefaultUserService()
	err = userService.ForgotPassword(userModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// ConfirmForgotPasswordV1 - Submit a forgot password request
func ConfirmForgotPasswordV1(w http.ResponseWriter, r *http.Request) {
	otpModel := model.DecodeRequestToOtp(r)
	userService := service.CreateDefaultUserService()
	err := userService.ConfirmForgotPassword(otpModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func getSessionToken(r *http.Request) string {
	return r.Header.Get("x-session-token")
}
