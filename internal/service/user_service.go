package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/third-place/user-service/internal/db"
	"github.com/third-place/user-service/internal/entity"
	kafka2 "github.com/third-place/user-service/internal/kafka"
	"github.com/third-place/user-service/internal/mapper"
	"github.com/third-place/user-service/internal/model"
	"github.com/third-place/user-service/internal/repository"
	"github.com/third-place/user-service/internal/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"log"
	"os"
	"strings"
)

type UserService struct {
	cognitoUserPool     string
	cognitoClientID     string
	cognitoClientSecret string
	cognito             *cognitoidentityprovider.CognitoIdentityProvider
	awsRegion           string
	userRepository      *repository.UserRepository
	inviteRepository    *repository.InviteRepository
	kafkaWriter         *kafka.Producer
}

const UserPasswordAuth = "USER_PASSWORD_AUTH"
const AuthFlowRefreshToken = "REFRESH_TOKEN_AUTH"
const AuthResponseChallenge = "NEW_PASSWORD_REQUIRED"
const JwkTokenUrl = "https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json"

func CreateDefaultUserService() *UserService {
	conn := db.CreateDefaultConnection()
	return CreateUserService(
		repository.CreateUserRepository(conn),
		repository.CreateInviteRepository(conn),
		kafka2.CreateWriter(),
	)
}

func CreateUserService(
	userRepository *repository.UserRepository,
	inviteRepository *repository.InviteRepository,
	kafkaWriter *kafka.Producer,
) *UserService {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &UserService{
		cognito:             cognitoidentityprovider.New(sess),
		cognitoUserPool:     os.Getenv("USER_POOL_ID"),
		cognitoClientID:     os.Getenv("COGNITO_CLIENT_ID"),
		cognitoClientSecret: os.Getenv("COGNITO_CLIENT_SECRET"),
		awsRegion:           os.Getenv("AWS_REGION"),
		userRepository:      userRepository,
		inviteRepository:    inviteRepository,
		kafkaWriter:         kafkaWriter,
	}
}

func (s *UserService) GetUserFromUsername(username string) (*model.User, error) {
	userEntity, err := s.userRepository.GetUserFromUsername(username)
	if err != nil {
		return nil, err
	}
	return mapper.MapUserEntityToUser(userEntity), nil
}

func (s *UserService) GetUserFromUuid(userUuid uuid.UUID) (*model.User, error) {
	userEntity, err := s.userRepository.GetUserFromUuid(userUuid)
	if err != nil {
		return nil, err
	}
	return mapper.MapUserEntityToUser(userEntity), nil
}

func (s *UserService) CreateUser(newUser *model.NewUser) (*model.User, error) {
	minSize, digit, special, lower, upper := util.ValidatePassword(newUser.Password)
	if !minSize || !digit || !special || !lower || !upper {
		log.Print("cannot create user, invalid password")
		msg := "passwords: "
		var errs []string
		if !minSize {
			errs = append(errs, "must be at least 8 characters")
		}
		if !digit {
			errs = append(errs, "need at least one digit")
		}
		if !special {
			errs = append(errs, "need at least one special character")
		}
		if !lower {
			errs = append(errs, "need at least one lower case letter")
		}
		if !upper {
			errs = append(errs, "need at least one upper case letter")
		}
		return nil, util.NewInputFieldError(
			"password",
			msg+strings.Join(errs, ", "),
		)
	}
	invite, err := s.inviteRepository.FindOneByCode(newUser.InviteCode)
	if err != nil {
		log.Print("error finding invite :: ", err)
		return nil, util.NewInputFieldError(
			"inviteCode",
			"invite code not found",
		)
	}
	if invite.Claimed {
		log.Print("attempting to use a claimed invite :: ", newUser.Email, newUser.InviteCode)
		return nil, util.NewInputFieldError(
			"inviteCode",
			"there was a problem with your invite code",
		)
	}
	user := mapper.MapNewUserModelToEntity(newUser)
	user.InviteID = invite.ID
	result := s.userRepository.Create(user)
	if result.Error != nil {
		search, _ := s.userRepository.GetUserFromUsername(newUser.Username)
		if search != nil {
			return nil, util.NewInputFieldError(
				"username",
				"username already in use",
			)
		}
		search, _ = s.userRepository.GetUserFromEmail(newUser.Email)
		if search != nil {
			return nil, util.NewInputFieldError(
				"email",
				"email already registered, try logging in",
			)
		}
		return nil, errors.New("error creating user")
	}
	userEntity, err := s.userRepository.GetUserFromUsername(newUser.Username)
	if err != nil {
		log.Print("user not found :: ", err)
		return nil, err
	}
	response, err := s.PublishToCognito(userEntity, newUser.Password)
	if err != nil {
		log.Print("error creating cognito user :: ", err)
		s.userRepository.Delete(user)
		return nil, errors.New("error creating user")
	}
	invite.Claimed = true
	s.inviteRepository.Save(invite)
	user.CognitoId = uuid.MustParse(*response.UserSub)
	result = s.userRepository.Save(user)
	if result.Error != nil {
		log.Print("error updating user with cognito ID :: ", result.Error)
	}
	userModel := mapper.MapUserEntityToModel(user)
	err = s.publishUserToKafka(user)
	if err != nil {
		log.Print("error publishing to kafka :: ", err)
	}
	return userModel, nil
}

func (s *UserService) PublishToCognito(user *entity.User, password string) (*cognitoidentityprovider.SignUpOutput, error) {
	return s.cognito.SignUp(&cognitoidentityprovider.SignUpInput{
		Username: aws.String(user.CurrentEmail),
		Password: aws.String(password),
		ClientId: aws.String(s.cognitoClientID),
	})
}

func (s *UserService) UpdateUser(userModel *model.User) error {
	userEntity, err := s.userRepository.GetUserFromUuid(uuid.MustParse(userModel.Uuid))
	if err != nil {
		return err
	}
	userEntity.UpdateUserProfileFromModel(userModel)
	s.userRepository.Save(userEntity)
	_ = s.publishUserToKafka(userEntity)
	return nil
}

func (s *UserService) CreateSession(newSession *model.NewSession) (*AuthResponse, error) {
	if newSession.Email == "" {
		return nil, util.NewInputFieldError(
			"email",
			"email address is required",
		)
	}
	if newSession.Password == "" {
		return nil, util.NewInputFieldError(
			"password",
			"password is required",
		)
	}
	search, _ := s.userRepository.GetUserFromEmail(newSession.Email)
	if search == nil {
		return nil, util.NewInputFieldError(
			"email",
			"email not found, do you need to sign up?",
		)
	}
	response, err := s.cognito.InitiateAuth(&cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String(UserPasswordAuth),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(newSession.Email),
			"PASSWORD": aws.String(newSession.Password),
		},
		ClientId: aws.String(s.cognitoClientID),
	})

	if err != nil {
		log.Print("login failed", err.Error())
		return nil, util.NewInputFieldError(
			"password",
			"login failed, do you need a password reset?",
		)
	}

	if response.AuthenticationResult != nil {
		log.Print("updating user tokens with response from AWS for user ID: ", search.ID, ", response: ", response.String())
		s.updateUserTokens(search, response.AuthenticationResult)
		return createSessionResponse(search, response), nil
	}

	s.updateUserWithCreateSessionResult(search, response)
	log.Print("created session from AWS: ", response.String())
	return createChallengeSessionResponse(search, response), nil
}

func (s *UserService) ProvideChallengeResponse(passwordReset *model.PasswordReset) *AuthResponse {
	log.Print("provide challenge response :: ", passwordReset)
	user, err := s.userRepository.GetUserFromEmail(passwordReset.Email)

	if err != nil {
		log.Print("user not found")
		return createAuthFailedSessionResponse("user not found")
	}

	log.Print("requesting reset with: ", passwordReset.Email, ", session: ", user.LastSessionToken)

	data := &cognitoidentityprovider.RespondToAuthChallengeInput{
		ChallengeName: aws.String(AuthResponseChallenge),
		ChallengeResponses: map[string]*string{
			"USERNAME":     aws.String(passwordReset.Email),
			"NEW_PASSWORD": aws.String(passwordReset.Password),
		},
		ClientId: aws.String(s.cognitoClientID),
		Session:  aws.String(user.LastSessionToken),
	}

	response, err := s.cognito.RespondToAuthChallenge(data)

	if err != nil {
		log.Print("error responding to auth challenge: ", err)
		return createAuthFailedSessionResponse("auth failed")
	}

	log.Print("response from provide challenge: ", response.String())

	if response.AuthenticationResult != nil {
		s.updateUserTokens(user, response.AuthenticationResult)
	}

	return createChallengeResponse(response)
}

func (s *UserService) GetSession(sessionToken *model.SessionToken) (*model.Session, error) {
	keySet, jwkErr := jwk.Fetch(fmt.Sprintf(JwkTokenUrl, s.awsRegion, s.cognitoUserPool))
	if jwkErr != nil {
		log.Print("error fetching jwk: ", jwkErr)
		return nil, errors.New("jwk fetch error")
	}

	token, parseErr := jwt.Parse(sessionToken.Token, func(token *jwt.Token) (interface{}, error) {
		kid, _ := token.Header["kid"].(string)
		keys := keySet.LookupKeyID(kid)
		if len(keys) > 0 {
			return keys[0].Materialize()
		}
		log.Print("error finding user session")
		return nil, errors.New("no session found")
	})
	if parseErr != nil {
		log.Print("jwt parse error", parseErr)
		return nil, parseErr
	}

	claims := token.Claims.(jwt.MapClaims)
	if err := claims.Valid(); err != nil || claims.VerifyAudience(s.cognitoClientID, false) == false {
		log.Print("token verification failed with: ", err)
		return nil, errors.New("verification failed")
	}

	response, err := s.cognito.GetUser(&cognitoidentityprovider.GetUserInput{AccessToken: aws.String(sessionToken.Token)})
	if err != nil {
		log.Print("error retrieving user: ", err)
		return nil, err
	}
	user, err := s.userRepository.GetUserFromSessionToken(sessionToken.Token)
	if err != nil {
		log.Print("user does not match jwt: ", response.String(), " and user: ", user)
		return nil, errors.New("user does not match jwt")
	}
	return model.CreateSession(mapper.MapUserEntityToUser(user), sessionToken.Token), nil
}

func (s *UserService) RefreshSession(sessionRefresh *model.SessionRefresh) *AuthResponse {
	log.Print("request refresh session :: ", sessionRefresh.Token)
	user, err := s.userRepository.GetUserFromSessionToken(sessionRefresh.Token)

	if err != nil {
		log.Print("error finding user :: ", err)
		return createAuthFailedSessionResponse("auth failed")
	}

	if user.LastRefreshToken == "" {
		log.Print("no available refresh tokens")
		return createAuthFailedSessionResponse("no available refresh tokens")
	}

	result, err := s.cognito.InitiateAuth(&cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String(AuthFlowRefreshToken),
		AuthParameters: map[string]*string{
			"REFRESH_TOKEN": aws.String(user.LastRefreshToken),
			"DEVICE_KEY":    aws.String(user.DeviceKey),
		},
		ClientId: aws.String(s.cognitoClientID),
	})

	if err != nil {
		log.Print("error refreshing user session :: ", err)
		return createAuthFailedSessionResponse("auth failed")
	}

	s.updateUserTokens(user, result.AuthenticationResult)
	return createSuccessfulRefreshResponse(result)
}

func (s *UserService) DeleteSession(sessionToken *model.SessionToken) error {
	_, err := s.cognito.GlobalSignOut(&cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: &sessionToken.Token,
	})
	if err != nil {
		return errors.New("something failed")
	}
	return nil
}

func (s *UserService) BanUser(sessionUser *entity.User, userEntity *entity.User) error {
	if !s.canAdminister(sessionUser, userEntity) {
		return errors.New("cannot ban user")
	}
	userEntity.IsBanned = true
	s.userRepository.Save(userEntity)
	_ = s.publishUserToKafka(userEntity)
	return nil
}

func (s *UserService) UnbanUser(sessionUser *entity.User, userEntity *entity.User) error {
	if !s.canAdminister(sessionUser, userEntity) {
		return errors.New("cannot ban user")
	}
	userEntity.IsBanned = false
	s.userRepository.Save(userEntity)
	_ = s.publishUserToKafka(userEntity)
	return nil
}

func (s *UserService) SubmitOTP(otp *model.Otp) error {
	_, err := s.cognito.ConfirmSignUp(&cognitoidentityprovider.ConfirmSignUpInput{
		ConfirmationCode: aws.String(otp.Code),
		Username:         aws.String(otp.User.Username),
		ClientId:         aws.String(s.cognitoClientID),
	})
	if err != nil {
		log.Print("err with OTP :: ", err.Error())
	}
	return err
}

func (s *UserService) ForgotPassword(user *model.User) error {
	_, err := s.cognito.ForgotPassword(&cognitoidentityprovider.ForgotPasswordInput{
		Username: aws.String(user.Username),
		ClientId: aws.String(s.cognitoClientID),
	})
	if err != nil {
		log.Print("err with forgot password :: ", err.Error())
	}
	return err
}

func (s *UserService) ConfirmForgotPassword(otp *model.Otp) error {
	_, err := s.cognito.ConfirmForgotPassword(&cognitoidentityprovider.ConfirmForgotPasswordInput{
		Username:         aws.String(otp.User.Username),
		Password:         aws.String(otp.User.Password),
		ClientId:         aws.String(s.cognitoClientID),
		ConfirmationCode: aws.String(otp.Code),
	})
	if err != nil {
		log.Print("err with forgot password :: ", err.Error())
	}
	return err
}

func (s *UserService) GetInvites(offset int) []*model.Invite {
	invites := s.inviteRepository.FindInvites(offset)
	return mapper.MapInviteEntitiesToModels(invites)
}

func (s *UserService) GetInvite(code string) (*model.Invite, error) {
	invite, err := s.inviteRepository.FindOneByCode(code)
	if err != nil {
		return nil, err
	}
	return mapper.MapInviteEntityToModel(invite), nil
}

func (s *UserService) CreateInviteFromCode(code string) (*model.Invite, error) {
	invite := &entity.Invite{
		Code: code,
	}
	result := s.inviteRepository.Create(invite)
	if result.Error != nil {
		return nil, result.Error
	}
	return mapper.MapInviteEntityToModel(invite), nil
}

func (s *UserService) publishUserToKafka(userEntity *entity.User) error {
	topic := "users"
	userModel := mapper.MapUserEntityToModel(userEntity)
	userData, _ := json.Marshal(userModel)
	return s.kafkaWriter.Produce(
		&kafka.Message{
			Value: userData,
			TopicPartition: kafka.TopicPartition{Topic: &topic,
				Partition: kafka.PartitionAny},
		},
		nil)
}

func (s *UserService) canAdminister(sessionUser *entity.User, user *entity.User) bool {
	if sessionUser.IsBanned || sessionUser.Role == "user" {
		return false
	}
	if user.Role == "moderator" {
		return sessionUser.Role == "admin"
	}
	if user.Role == "admin" {
		return false
	}
	return true
}

func (s *UserService) updateUserWithCreateSessionResult(user *entity.User, result *cognitoidentityprovider.InitiateAuthOutput) {
	user.SRP = *result.ChallengeParameters["USER_ID_FOR_SRP"]
	user.LastSessionToken = *result.Session
	s.userRepository.Save(user)
}

func (s *UserService) updateUserTokens(user *entity.User, result *cognitoidentityprovider.AuthenticationResultType) {
	if result.NewDeviceMetadata != nil {
		user.DeviceGroupKey = *result.NewDeviceMetadata.DeviceGroupKey
		user.DeviceKey = *result.NewDeviceMetadata.DeviceKey
	}
	user.LastAccessToken = *result.AccessToken
	user.LastIdToken = *result.IdToken
	if result.RefreshToken != nil {
		user.LastRefreshToken = *result.RefreshToken
	}
	s.userRepository.Save(user)
}
