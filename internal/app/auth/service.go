package auth

import (
	"context"
	"errors"
	"fmt"
	"time"
	"uiren/internal/app/users"
	"uiren/internal/infrastracture/hasher"
	jwt_maker "uiren/internal/infrastracture/jwt"
	yandex_sender "uiren/internal/infrastracture/mail/yandex"
	"uiren/pkg/logger"
)

type userService interface {
	GetUserForLogin(ctx context.Context, indetifier string) (users.UserDTO, error)
	CreateUser(ctx context.Context, params users.CreateUserDTO) (string, error)
	EnableUser(ctx context.Context, username string) error
}

type jwtMaker interface {
	NewToken(payload jwt_maker.PayloadDTO) (string, error)
}

type verificationCodeRepository interface {
	createVerificationCode(ctx context.Context, req CreateVerificationCodeRequest) error
	getVerificationCode(ctx context.Context, username string) (Verification, error)
}

type AuthService struct {
	userService  userService
	jwtMaker     jwtMaker
	verifRepo    verificationCodeRepository
	verifCodeTTL time.Duration
}

func NewAuthService(userService userService, jwtMaker jwtMaker, verifRepo verificationCodeRepository) *AuthService {
	return &AuthService{
		userService: userService,
		jwtMaker:    jwtMaker,
		verifRepo:   verifRepo,
	}
}

func (s *AuthService) SetVerificationCodeTTL(verifCodeTTL time.Duration) {
	s.verifCodeTTL = verifCodeTTL
}

func (s *AuthService) SignIn(ctx context.Context, params LoginParams) (string, error) {
	logger.Info("AuthService.SignIn new request")

	user, err := s.userService.GetUserForLogin(ctx, params.Identificator)
	if err != nil {
		logger.Error("AuthService.SignIn getUser: ", err)
		if errors.Is(err, users.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if !hasher.BcryptComparePasswordAndHash(params.Password, user.Password) {
		logger.Error("AuthService.SignIn BcryptComparePasswordAndHash: ", err)
		return "", ErrInvalidCredentials
	}

	token, err := s.jwtMaker.NewToken(jwt_maker.PayloadDTO{
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		IsAdmin:   user.IsAdmin,
	})
	if err != nil {
		logger.Error("AuthService.SignIn NewToken: ", err)
		return "", err
	}

	return token, nil
}

func (s *AuthService) Register(ctx context.Context, params RegisterParams) (string, error) {
	logger.Info("AuthService.Register new request")

	userID, err := s.userService.CreateUser(ctx, params.DTO)
	if err != nil {
		return "", err
	}

	verifReq := CreateVerificationCodeRequest{
		Username: params.DTO.Username,
		Email:    params.DTO.Email,
		Code:     generateAlphanumericCode(10),
		Duration: s.verifCodeTTL,
	}

	if err := s.verifRepo.createVerificationCode(ctx, verifReq); err != nil {
		logger.Error("UserService.CreateUser verifRepo.createVerificationCode error: ", err)
		return "", err
	}

	errChan := make(chan error, 1)
	go func() {
		defer close(errChan)
		err := yandex_sender.SendEmail(
			"Uiren. Email Verification",
			fmt.Sprintf("localhost:8080/api/verify/%s/%s", verifReq.Username, verifReq.Code),
			[]string{verifReq.Email},
			[]string{}, []string{}, []string{})

		errChan <- err
	}()
	go func() {
		if err := <-errChan; err != nil {
			logger.Error("UserService.CreateUser emailSender.SendEmail error: ", err)
		}
	}()

	return userID, nil
}

func (s *AuthService) VerifyUser(ctx context.Context, username, code string) error {
	logger.Info("UserService.VerifyUser new request")

	verification, err := s.verifRepo.getVerificationCode(ctx, username)
	if err != nil {
		logger.Error("UserService.VerifyUser verifRepo.getVerificationCode: ", err)
		return err
	}

	if time.Now().Unix() > verification.ExpiresAt.Unix() {
		logger.Error("UserService.VerifyUser ", ErrVerificationExpired)
		return ErrVerificationExpired
	}

	if verification.Code != code {
		logger.Error("UserService.VerifyUser ", ErrVerificationInvalid)
		return ErrVerificationInvalid
	}

	return s.userService.EnableUser(ctx, username)
}
