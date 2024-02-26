package util

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/third-place/user-service/internal/model"
)

func GetSessionTokenModel(c *gin.Context) *model.SessionToken {
	sessionToken := c.GetHeader("x-session-token")
	if sessionToken == "" {
		return nil
	}
	return &model.SessionToken{
		Token: sessionToken,
	}
}

func GetSession(c *gin.Context) (*model.Session, error) {
	sessionToken := c.GetHeader("x-session-token")
	claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(sessionToken, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token not valid")
	}
	_, err = uuid.Parse(claims.UserUuid)
	if err != nil {
		return nil, err
	}
	return &model.Session{
		User: &model.User{
			Uuid: claims.UserUuid,
		},
		Token: sessionToken,
	}, nil
}
