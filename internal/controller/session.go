package controller

import (
	"encoding/json"
	"github.com/third-place/user-service/internal/model"
	"github.com/third-place/user-service/internal/service"
	"log"
	"net/http"
)

// CreateSessionV1 - Create a new session
func CreateSessionV1(w http.ResponseWriter, r *http.Request) {
	newSessionModel := model.DecodeRequestToNewSession(r)
	result, err := service.CreateDefaultUserService().CreateSession(newSessionModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		data, _ := json.Marshal(err)
		_, _ = w.Write(data)
		return
	}
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(result)
	_, _ = w.Write(data)
}

// RespondToChallengeV1 - Respond to an authentication challenge with a password reset
func RespondToChallengeV1(w http.ResponseWriter, r *http.Request) {
	passwordResetModel, err := model.DecodeRequestToPasswordReset(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result := service.CreateDefaultUserService().ProvideChallengeResponse(passwordResetModel)
	data, _ := json.Marshal(result)
	_, _ = w.Write(data)
}

// GetSessionV1 - validate a session token
func GetSessionV1(w http.ResponseWriter, r *http.Request) {
	sessionToken := model.DecodeRequestToSessionToken(r)
	session, err := service.CreateDefaultUserService().GetSession(sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Print("sanity: ", sessionToken.Token)
		log.Print(err)
		return
	}
	data, _ := json.Marshal(session)
	_, _ = w.Write(data)
}

// RefreshSessionV1 - refresh a session token
func RefreshSessionV1(w http.ResponseWriter, r *http.Request) {
	sessionToken := model.DecodeRequestToSessionRefresh(r)
	response := service.CreateDefaultUserService().RefreshSession(sessionToken)
	_, _ = w.Write(response.ToJson())
}

// DeleteSessionV1 - Delete a user's session (log out)
func DeleteSessionV1(w http.ResponseWriter, r *http.Request) {
	sessionToken := model.DecodeRequestToSessionToken(r)
	err := service.CreateDefaultUserService().DeleteSession(sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
