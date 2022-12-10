package controller

import (
	"encoding/json"
	"github.com/third-place/user-service/internal/model"
	"github.com/third-place/user-service/internal/service"
	"github.com/third-place/user-service/internal/util"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// CreateInviteV1 -- create new invites for new users
func CreateInviteV1(w http.ResponseWriter, r *http.Request) {
	userService := service.CreateDefaultUserService()
	sessionToken := getSessionToken(r)
	sessionModel := &model.SessionToken{
		Token: sessionToken,
	}
	session, err := userService.GetSession(sessionModel)
	if err != nil || session.User.Role == model.USER {
		w.WriteHeader(http.StatusForbidden)
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	invite, err := userService.CreateInviteFromCode(code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(invite)
	_, _ = w.Write(data)
}

// GetInvitesV1 -- get a list of invites
func GetInvitesV1(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	offset := 0
	if value := query.Get("offset"); value != "" {
		offset, err := strconv.Atoi(value)
		if err != nil || offset < 0 || offset > 100 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	userService := service.CreateDefaultUserService()
	invites := userService.GetInvites(offset)
	data, _ := json.Marshal(invites)
	_, _ = w.Write(data)
}
