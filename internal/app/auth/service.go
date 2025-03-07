package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"uiren/internal/app/users"
	"uiren/internal/infrastracture/hasher"
	jwt_maker "uiren/internal/infrastracture/jwt"
	yandex_sender "uiren/internal/infrastracture/mail/yandex"
	"uiren/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type userService interface {
	GetUserForLogin(ctx context.Context, indetifier string) (users.UserDTO, error)
	CreateUser(ctx context.Context, params users.CreateUserDTO) (string, error)
	EnableUser(ctx context.Context, username string) error
}

type jwtMaker interface {
	NewToken(payload jwt_maker.PayloadDTO) (string, error)
}

type redisClient interface {
	Set(ctx context.Context, key string, value interface{}, ttl *time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type verificationCodeRepository interface {
	createVerificationCode(ctx context.Context, req CreateVerificationCodeRequest) error
	getVerificationCode(ctx context.Context, username string) (Verification, error)
}

type AuthService struct {
	userService     userService
	jwtMaker        jwtMaker
	verifRepo       verificationCodeRepository
	verifCodeTTL    time.Duration
	refreshTokenTTL time.Duration
	redisClient     redisClient
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

func (s *AuthService) SetRefreshTokenTTL(refreshTokenTTL time.Duration) {
	s.refreshTokenTTL = refreshTokenTTL
}

func (s *AuthService) WithRedisClient(redisClient redisClient) {
	s.redisClient = redisClient
}

func (s *AuthService) SignIn(ctx context.Context, params LoginParams) (string, string, error) {
	logger.Info("AuthService.SignIn new request")

	user, err := s.userService.GetUserForLogin(ctx, params.Identificator)
	if err != nil {
		logger.Error("AuthService.SignIn getUser: ", err)
		if errors.Is(err, users.ErrUserNotFound) {
			return "", "", ErrInvalidCredentials
		}
		return "", "", err
	}

	if !hasher.BcryptComparePasswordAndHash(params.Password, user.Password) {
		logger.Error("AuthService.SignIn BcryptComparePasswordAndHash: ", err)
		return "", "", ErrInvalidCredentials
	}

	payload := jwt_maker.PayloadDTO{
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		IsAdmin:   user.IsAdmin,
	}
	return s.generateTokens(ctx, payload)
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

func (s *AuthService) RefreshToken(ctx context.Context, token string) (string, string, error) {
	logger.Info("AuthService.RefreshToken new request")
	if s.redisClient == nil {
		return "", "", fmt.Errorf("redis client is not initialized")
	}
	key := refreshTokenKey(token)

	payloadJSON, err := s.redisClient.Get(ctx, key)
	if err != nil {
		logger.Error("AuthService.RefreshToken redisClient.Get: ", err)
		if errors.Is(err, redis.Nil) {
			return "", "", ErrRefreshTokenNotFound
		}
		return "", "", err
	}

	var payload jwt_maker.PayloadDTO
	if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
		logger.Error("AuthService.RefreshToken json.Unmarshal: ", err)
		return "", "", ErrInvalidToken
	}

	if err := s.redisClient.Delete(ctx, key); err != nil {
		logger.Error("AuthService.RefreshToken redisClient.Delete: ", err)
	}

	return s.generateTokens(ctx, payload)
}

func (s *AuthService) generateTokens(ctx context.Context, payload jwt_maker.PayloadDTO) (string, string, error) {
	token, err := s.jwtMaker.NewToken(payload)
	if err != nil {
		logger.Error("AuthService.generateTokens NewToken: ", err)
		return "", "", err
	}

	refreshToken := generateAlphanumericCode(50)

	if s.redisClient != nil {
		key := refreshTokenKey(refreshToken)

		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			logger.Error("AuthService.generateTokens json.Marshal: ", err)
			return "", "", err
		}

		if err := s.redisClient.Set(ctx, key, string(payloadJSON), &s.refreshTokenTTL); err != nil {
			logger.Error("AuthService.generateTokens redisClient.Set: ", err)
			return "", "", err
		}
	}

	return token, refreshToken, nil
}
