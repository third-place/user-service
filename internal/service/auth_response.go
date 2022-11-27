package service

import (
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/danielmunro/otto-user-service/internal/model"
)

type AuthResponseType int

const (
	Unknown                     AuthResponseType = iota + 1
	ChallengeNewPassword        AuthResponseType = iota
	SessionAuthenticated        AuthResponseType = iota
	SessionFailedAuthentication AuthResponseType = iota
)

type AuthResponse struct {
	AuthResponse AuthResponseType
	Token        *string
	User         *model.User
	Message      string
}

func createSuccessfulRefreshResponse(response *cognitoidentityprovider.InitiateAuthOutput) *AuthResponse {
	return &AuthResponse{
		AuthResponse: SessionAuthenticated,
		Token:        response.AuthenticationResult.AccessToken,
	}
}
