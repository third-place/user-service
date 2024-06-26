package service

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/third-place/user-service/internal/db"
	"github.com/third-place/user-service/internal/entity"
	"github.com/third-place/user-service/internal/kafka"
	"github.com/third-place/user-service/internal/mapper"
	"github.com/third-place/user-service/internal/model"
	"github.com/third-place/user-service/internal/repository"
	"github.com/third-place/user-service/internal/util"
	"log"
	"time"
)

type UserService struct {
	userRepository   *repository.UserRepository
	inviteRepository *repository.InviteRepository
	mailService      *MailService
	kafkaWriter      kafka.Producer
	securityService  *SecurityService
}

func CreateTestUserService() *UserService {
	conn := util.SetupTestDatabase()
	writer, err := util.CreateTestProducer()
	if err != nil {
		log.Fatal("error creating test kafka writer :: ", err)
	}
	return &UserService{
		repository.CreateUserRepository(conn),
		repository.CreateInviteRepository(conn),
		CreateTestMailService(),
		writer,
		CreateSecurityService(),
	}
}

func CreateUserService() *UserService {
	conn := db.CreateDefaultConnection()
	writer, err := kafka.CreateProducer()
	if err != nil {
		log.Fatal("error creating kafka writer :: ", err)
	}
	return &UserService{
		repository.CreateUserRepository(conn),
		repository.CreateInviteRepository(conn),
		CreateMailService(),
		writer,
		CreateSecurityService(),
	}
}

func isPrivileged(role model.Role) bool {
	return role == model.MODERATOR || role == model.ADMIN
}

func filterUserModelForAccess(viewerUser *model.User, user *model.User) *model.User {
	if viewerUser == nil || (user.Uuid != viewerUser.Uuid && !isPrivileged(viewerUser.Role)) {
		user.Birthday = ""
		user.Email = ""
	}
	return user
}

func filterUserModelsForAccess(viewerUser *model.User, users []*model.User) []*model.User {
	for i, user := range users {
		users[i] = filterUserModelForAccess(viewerUser, user)
	}
	return users
}

func (s *UserService) GetUserFromUsername(viewerUser *model.User, username string) (*model.User, error) {
	userEntity, err := s.userRepository.GetUserFromUsername(username)
	if err != nil {
		return nil, err
	}
	return filterUserModelForAccess(viewerUser, mapper.MapUserEntityToModel(userEntity)), nil
}

func (s *UserService) GetUsers(viewerUser *model.User, offset int) []*model.User {
	userEntities := s.userRepository.GetUsers(offset)
	return filterUserModelsForAccess(viewerUser, mapper.MapUserEntitiesToModels(userEntities))
}

func (s *UserService) GetUserFromUuid(viewerUser *model.User, userUuid uuid.UUID) (*model.User, error) {
	userEntity, err := s.userRepository.GetUserFromUuid(userUuid)
	if err != nil {
		return nil, err
	}
	return filterUserModelForAccess(viewerUser, mapper.MapUserEntityToModel(userEntity)), nil
}

func (s *UserService) CreateUser(newUser *model.NewUser) (*model.User, error) {
	minSize := util.ValidatePassword(newUser.Password)
	if !minSize {
		log.Print("cannot create user, invalid password")
		return nil, util.NewInputFieldError(
			"password",
			"passwords must be at least 8 characters",
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
	invite.Claimed = true
	s.inviteRepository.Save(invite)
	user.OTP = util.GenerateCode()
	user.Password, _ = util.HashPassword(newUser.Password)
	s.userRepository.Save(user)
	_, err = s.mailService.SendVerificationEmail(user)
	if err != nil {
		log.Print(err)
	}
	userModel := mapper.MapUserEntityToModel(user)
	err = s.publishUserToKafka(user)
	if err != nil {
		log.Print("error publishing to kafka :: ", err)
	}
	return userModel, nil
}

func (s *UserService) UpdateUser(session *model.Session, userModel *model.User) error {
	if session.User.Uuid != userModel.Uuid {
		return errors.New("unauthorized")
	}
	userEntity, err := s.userRepository.GetUserFromUuid(uuid.MustParse(userModel.Uuid))
	if err != nil {
		return err
	}
	userEntity.UpdateUserFromModel(userModel)
	s.userRepository.Save(userEntity)
	_ = s.publishUserToKafka(userEntity)
	return nil
}

func (s *UserService) CreateSession(newSession *model.NewSession) (*model.Session, error) {
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
	if !util.CheckPasswordHash(newSession.Password, search.Password) {
		return nil, errors.New("authentication failed")
	}
	token, err := s.getJWT(search)
	if err != nil {
		return nil, err
	}
	return &model.Session{
		Token: token,
		User:  mapper.MapUserEntityToModel(search),
	}, nil
}

func (s *UserService) GetSession(sessionToken *model.SessionToken) (*model.Session, error) {
	claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(sessionToken.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return util.JwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token not valid")
	}
	userUuid, err := uuid.Parse(claims.UserUuid)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepository.GetUserFromUuid(userUuid)
	if err != nil {
		return nil, err
	}
	return model.CreateSession(mapper.MapUserEntityToModel(user), sessionToken.Token), nil
}

func (s *UserService) RefreshSession(sessionToken *model.SessionToken) (*model.SessionToken, error) {
	claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(sessionToken.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return util.JwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	if time.Until(claims.ExpiresAt.Time) > 24*4*time.Hour {
		return nil, errors.New("token not ready for refresh")
	}
	expirationTime := time.Now().Add(1 * time.Hour)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(util.JwtKey)
	return &model.SessionToken{
		Token: tokenString,
	}, nil
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
		return errors.New("code mismatch")
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
	_, err = s.mailService.SendPasswordResetEmail(userEntity)
	if err != nil {
		log.Print(err)
	}
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
	minSize := util.ValidatePassword(otp.User.Password)
	if !minSize {
		return errors.New("password too short")
	}
	userEntity.Password, err = util.HashPassword(otp.User.Password)
	if err != nil {
		log.Print(err.Error())
		return errors.New("error with password")
	}
	userEntity.Verified = true
	s.userRepository.Save(userEntity)
	return nil
}

func (s *UserService) GetInvites(session *model.Session, offset int) ([]*model.Invite, error) {
	if !s.securityService.IsModerator(session) {
		return nil, errors.New("not allowed")
	}
	invites := s.inviteRepository.FindInvites(offset)
	return mapper.MapInviteEntitiesToModels(invites), nil
}

func (s *UserService) GetInvite(code string) (*model.Invite, error) {
	invite, err := s.inviteRepository.FindOneByCode(code)
	if err != nil {
		return nil, err
	}
	return mapper.MapInviteEntityToModel(invite), nil
}

func (s *UserService) CreateInviteFromCode(session *model.Session, code string) (*model.Invite, error) {
	if !s.securityService.IsModerator(session) {
		return nil, errors.New("not allowed")
	}
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
	return s.kafkaWriter.Produce(kafka.CreateMessage(userData, topic), nil)
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

func (s *UserService) getJWT(user *entity.User) (string, error) {
	claims := model.NewClaims(user.Uuid)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(util.JwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
