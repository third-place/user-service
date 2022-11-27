package service

import "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

func createChallengeResponse(response *cognitoidentityprovider.RespondToAuthChallengeOutput) *AuthResponse {
	if response.AuthenticationResult != nil {
		return &AuthResponse{
			AuthResponse: SessionAuthenticated,
			Token:        response.AuthenticationResult.AccessToken,
		}
	}

	return &AuthResponse{
		AuthResponse: SessionFailedAuthentication,
	}
}
