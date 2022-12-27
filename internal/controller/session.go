package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/third-place/user-service/internal/model"
	"github.com/third-place/user-service/internal/service"
	"net/http"
)

// CreateNewSessionV1 - Create a new session
func CreateNewSessionV1(c *gin.Context) {
	newSessionModel, err := model.DecodeRequestToNewSession(c.Request)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	result, err := service.CreateUserService().CreateSession(newSessionModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

// GetSessionV1 - validate a session token
func GetSessionV1(c *gin.Context) {
	sessionToken := model.DecodeRequestToSessionToken(c.Request)
	session, err := service.CreateUserService().GetSession(sessionToken)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	c.JSON(http.StatusOK, session)
}

// RefreshSessionV1 - refresh a session token
func RefreshSessionV1(c *gin.Context) {
	sessionToken := model.DecodeRequestToSessionToken(c.Request)
	session, err := service.CreateUserService().RefreshSession(sessionToken)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, session)
}
