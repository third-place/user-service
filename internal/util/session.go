package util

import (
	"github.com/gin-gonic/gin"
	"github.com/third-place/user-service/internal/model"
)

func GetSessionTokenModel(c *gin.Context) *model.SessionToken {
	sessionToken := c.GetHeader("x-session-token")
	return &model.SessionToken{
		Token: sessionToken,
	}
}
