package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/third-place/user-service/internal/db"
	"github.com/third-place/user-service/internal/model"
	"github.com/third-place/user-service/internal/repository"
	"github.com/third-place/user-service/internal/service"
	"github.com/third-place/user-service/internal/util"
	"log"
	"net/http"
)

// CreateNewUserV1 - Create a new user
func CreateNewUserV1(c *gin.Context) {
	newUserModel, err := model.DecodeRequestToNewUser(c.Request)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	user, err := service.CreateUserService().CreateUser(newUserModel)
	if err != nil {
		if _, ok := err.(*util.InputFieldError); ok {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusCreated, user)
}

// GetUserByUsernameV1 - Get a user by username
func GetUserByUsernameV1(c *gin.Context) {
	c.Header("Cache-Control", "max-age=30")
	username := c.Param("username")
	userService := service.CreateUserService()
	session, err := userService.GetSession(util.GetSessionTokenModel(c))
	if err != nil {
		log.Print(err)
		c.Status(http.StatusForbidden)
		return
	}
	user, err := service.CreateUserService().GetUserFromUsername(session.User, username)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUsersV1 - Get a list of users
func GetUsersV1(c *gin.Context) {
	userService := service.CreateUserService()
	_, err := userService.GetSession(util.GetSessionTokenModel(c))
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	offset, err := util.GetOffsetParam(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	session, err := userService.GetSession(util.GetSessionTokenModel(c))
	if err != nil || session.User.Role == model.USER {
		c.Status(http.StatusForbidden)
		return
	}
	c.Header("Cache-Control", "max-age=30")
	users := service.CreateUserService().GetUsers(session.User, offset)
	c.JSON(http.StatusOK, users)
}

// UpdateUserV1 - Update a user
func UpdateUserV1(c *gin.Context) {
	userModel, err := model.DecodeRequestToUser(c.Request)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	userService := service.CreateUserService()
	sessionModel := util.GetSessionTokenModel(c)
	session, err := userService.GetSession(sessionModel)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	err = userService.UpdateUser(session, userModel)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, userModel)
}

// BanUserV1 - ban a user
func BanUserV1(c *gin.Context) {
	usernameParam, success := c.Params.Get("username")
	if !success {
		return
	}
	userService := service.CreateUserService()
	userRepository := repository.CreateUserRepository(db.CreateDefaultConnection())
	sessionModel := util.GetSessionTokenModel(c)
	session, err := userService.GetSession(sessionModel)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	sessionUser, err := userRepository.GetUserFromUsername(session.User.Username)
	if err != nil || sessionUser.IsBanned {
		c.Status(http.StatusBadRequest)
		return
	}
	userEntity, err := userRepository.GetUserFromUsername(usernameParam)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err = userService.BanUser(sessionUser, userEntity)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
}

// UnbanUserV1 - ban a user
func UnbanUserV1(c *gin.Context) {
	usernameParam := c.Param("username")
	userService := service.CreateUserService()
	userRepository := repository.CreateUserRepository(db.CreateDefaultConnection())
	sessionModel := util.GetSessionTokenModel(c)
	session, err := userService.GetSession(sessionModel)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	sessionUser, err := userRepository.GetUserFromUsername(session.User.Username)
	if err != nil || sessionUser.IsBanned {
		c.Status(http.StatusBadRequest)
		return
	}
	userEntity, err := userRepository.GetUserFromUsername(usernameParam)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err = userService.UnbanUser(sessionUser, userEntity)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
}

// SubmitOTPV1 - Submit a new OTP
func SubmitOTPV1(c *gin.Context) {
	otpModel, err := model.DecodeRequestToOtp(c.Request)
	if err != nil {
		log.Print(err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	err = service.CreateUserService().SubmitOTP(otpModel)
	if err != nil {
		log.Print(err.Error())
		c.Status(http.StatusBadRequest)
	}
}

// SubmitForgotPasswordV1 - Submit a forgot password request
func SubmitForgotPasswordV1(c *gin.Context) {
	userModel, err := model.DecodeRequestToUser(c.Request)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err = service.CreateUserService().ForgotPassword(userModel)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
}

// ConfirmForgotPasswordV1 - Submit a forgot password request
func ConfirmForgotPasswordV1(c *gin.Context) {
	otpModel, err := model.DecodeRequestToOtp(c.Request)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err = service.CreateUserService().ConfirmForgotPassword(otpModel)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
}
