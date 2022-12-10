package service

import (
	"encoding/json"
	"github.com/third-place/user-service/internal/entity"
	"github.com/third-place/user-service/internal/mapper"
)

const challengeNewPasswordString = "ChallengeNewPassword"

func createSessionResponse(user *entity.User) *AuthResponse {
	return &AuthResponse{
		User: mapper.MapUserEntityToUser(user),
	}
}

func createAuthFailedSessionResponse(message string) *AuthResponse {
	token := ""
	return &AuthResponse{
		AuthResponse: SessionFailedAuthentication,
		Message:      message,
		Token:        &token,
	}
}

func getChallengeString(authResponse AuthResponseType) string {
	if authResponse == ChallengeNewPassword {
		return challengeNewPasswordString
	}
	return ""
}

func (c *AuthResponse) ToJson() []byte {
	data, _ := json.Marshal(map[string]string{
		"authResponse": getChallengeString(c.AuthResponse),
		"token":        *c.Token,
	})
	return data
}
