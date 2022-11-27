package service

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/danielmunro/otto-user-service/internal/entity"
	"github.com/danielmunro/otto-user-service/internal/mapper"
)

const challengeNewPasswordString = "ChallengeNewPassword"

func getAuthResponseFromChallenge(response string) AuthResponseType {
	if response == AuthResponseChallenge {
		return ChallengeNewPassword
	}
	return Unknown
}

func createSessionResponse(user *entity.User, response *cognitoidentityprovider.InitiateAuthOutput) *AuthResponse {
	return &AuthResponse{
		Token: response.AuthenticationResult.AccessToken,
		User:  mapper.MapUserEntityToUser(user),
	}
}

func createChallengeSessionResponse(user *entity.User, response *cognitoidentityprovider.InitiateAuthOutput) *AuthResponse {
	return &AuthResponse{
		AuthResponse: getAuthResponseFromChallenge(*response.ChallengeName),
		Token:        response.Session,
		User:         mapper.MapUserEntityToUser(user),
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
