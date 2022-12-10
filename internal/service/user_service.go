package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/third-place/user-service/internal/db"
	"github.com/third-place/user-service/internal/entity"
	kafka2 "github.com/third-place/user-service/internal/kafka"
	"github.com/third-place/user-service/internal/mapper"
	"github.com/third-place/user-service/internal/model"
	"github.com/third-place/user-service/internal/repository"
	"github.com/third-place/user-service/internal/util"
	"log"
	"os"
	"strings"
)

type UserService struct {
	cognitoUserPool     string
	cognitoClientID     string
	cognitoClientSecret string
	awsRegion           string
	userRepository      *repository.UserRepository
	inviteRepository    *repository.InviteRepository
	kafkaWriter         *kafka.Producer
	context             context.Context
}

func CreateUserService(context context.Context) *UserService {
	conn := db.CreateDefaultConnection()
	return &UserService{
		cognitoUserPool:     os.Getenv("USER_POOL_ID"),
		cognitoClientID:     os.Getenv("COGNITO_CLIENT_ID"),
		cognitoClientSecret: os.Getenv("COGNITO_CLIENT_SECRET"),
		awsRegion:           os.Getenv("AWS_REGION"),
		userRepository:      repository.CreateUserRepository(conn),
		inviteRepository:    repository.CreateInviteRepository(conn),
		kafkaWriter:         kafka2.CreateWriter(),
		context:             context,
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
	user.OTP = util.GenerateCode()
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
	if err != nil {
		log.Print("error creating cognito user :: ", err)
		s.userRepository.Delete(user)
		return nil, errors.New("error creating user")
	}
	invite.Claimed = true
	s.inviteRepository.Save(invite)
	userModel := mapper.MapUserEntityToModel(user)
	err = s.publishUserToKafka(user)
	if err != nil {
		log.Print("error publishing to kafka :: ", err)
	}
	return userModel, nil
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
	match := util.CheckPasswordHash(newSession.Password, search.Password)
	if !match {
		return nil, util.NewInputFieldError(
			"password",
			"login failed, do you need a password reset?",
		)
	}

	util.SessionManager.Put(s.context, "userUUID", search.Uuid.String())
	return createSessionResponse(search), nil
}

func (s *UserService) GetSession() (*model.Session, error) {
	userUuidStr := util.SessionManager.GetString(s.context, "userUUID")
	if userUuidStr == "" {
		return nil, errors.New("no session exists")
	}
	userUuid, err := uuid.Parse(userUuidStr)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepository.GetUserFromUuid(userUuid)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return model.CreateSession(mapper.MapUserEntityToUser(user)), nil
}

func (s *UserService) DeleteSession() error {
	return util.SessionManager.Destroy(s.context)
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
	userEntity, err := s.userRepository.GetUserFromEmail(otp.User.Email)
	if err != nil {
		return err
	}
	if userEntity.OTP != otp.Code {
		return errors.New("validation failed")
	}
	userEntity.Verified = true
	s.userRepository.Save(userEntity)
	return nil
}

func (s *UserService) ForgotPassword(user *model.User) error {
	userEntity, err := s.userRepository.GetUserFromEmail(user.Email)
	if err != nil {
		return err
	}
	userEntity.OTP = util.GenerateCode()
	s.userRepository.Save(userEntity)
	return nil
}

func (s *UserService) ConfirmForgotPassword(otp *model.Otp) error {
	userEntity, err := s.userRepository.GetUserFromEmail(otp.User.Email)
	if err != nil {
		return err
	}
	if userEntity.OTP != otp.Code {
		return errors.New("validation failed")
	}
	minSize, digit, special, lowercase, uppercase := util.ValidatePassword(otp.User.Password)
	if !minSize {
		return errors.New("password too short")
	}
	if !digit {
		return errors.New("password needs a number")
	}
	if !special {
		return errors.New("password needs a special character")
	}
	if !lowercase {
		return errors.New("password needs a lowercase letter")
	}
	if !uppercase {
		return errors.New("password needs an uppercase letter")
	}
	userEntity.Password, _ = util.HashPassword(otp.User.Password)
	userEntity.Verified = true
	s.userRepository.Save(userEntity)
	return nil
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
