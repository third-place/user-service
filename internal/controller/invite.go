package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/third-place/user-service/internal/model"
	"github.com/third-place/user-service/internal/service"
	"github.com/third-place/user-service/internal/util"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// CreateInviteV1 -- create new invites for new users
func CreateInviteV1(c *gin.Context) {
	userService := service.CreateUserService()
	session, err := userService.GetSession(util.GetSessionTokenModel(c))
	if err != nil || session.User.Role == model.USER {
		c.Status(http.StatusForbidden)
		return
	}
	code := util.GenerateCode()
	attempt := 0
	for {
		_, err = userService.GetInvite(code)
		if err.Error() == "no invite found" {
			break
		}
		code = util.GenerateCode()
		attempt += 1
		if attempt > 5 {
			c.Status(http.StatusInternalServerError)
			return
		}
	}
	invite, err := userService.CreateInviteFromCode(code)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, invite)
}

// GetInvitesV1 -- get a list of invites
func GetInvitesV1(c *gin.Context) {
	session, err := util.GetSession(c)
	offset, err := util.GetOffsetParam(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	invites, err := service.CreateUserService().GetInvites(session, offset)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	c.JSON(http.StatusOK, invites)
}
