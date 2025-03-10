package auth

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
	"uiren/internal/app/users"
	"uiren/internal/infrastracture/hasher"
	jwt_maker "uiren/internal/infrastracture/jwt"
	yandex_sender "uiren/internal/infrastracture/mail/yandex"
	"uiren/pkg/logger"

	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func init() {
	logger.InitLogger("info")
	yandex_sender.Init("", "", "")
}

func Test_authService_SignIn_success(t *testing.T) {
	t.Parallel()
	var (
		ctrl            = gomock.NewController(t)
		ctx             = context.TODO()
		paramsWithLogin = LoginParams{
			Identificator: "user",
			Password:      "123456Aa@",
		}
		paramsWithEmail = LoginParams{
			Identificator: "user@user.com",
			Password:      "123456@Aa",
		}
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		redisClient = NewMockredisClient(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)

		keyMock = refreshTokenMathcer{prefix: "auth:refresh:", randomLength: 50}
	)

	authService.WithRedisClient(redisClient)
	authService.SetRefreshTokenTTL(24 * time.Hour)

	t.Run("success with login", func(t *testing.T) {
		hashedPass, _ := hasher.BcryptHash(paramsWithLogin.Password)
		userService.EXPECT().GetUserForLogin(ctx, paramsWithLogin.Identificator).Return(users.UserDTO{Username: "user", Password: hashedPass}, nil)

		token1 := generateAlphanumericCode(100)
		payload := jwt_maker.PayloadDTO{
			Username: "user",
		}
		jwtMaker.EXPECT().NewToken(payload).Return(token1, nil)

		refreshToken1 := generateAlphanumericCode(50)

		payloadJSON, _ := json.Marshal(payload)
		redisClient.EXPECT().Set(ctx, keyMock, string(payloadJSON), &authService.refreshTokenTTL).Return(nil)

		token2, refreshToken2, err := authService.SignIn(ctx, paramsWithLogin)
		assert.NoError(t, err)
		assert.Equal(t, len(token1), len(token2))
		assert.Equal(t, len(refreshToken1), len(refreshToken2))
	})

	t.Run("success with email", func(t *testing.T) {
		hashedPass, _ := hasher.BcryptHash(paramsWithEmail.Password)
		userService.EXPECT().GetUserForLogin(ctx, paramsWithEmail.Identificator).Return(users.UserDTO{Username: "user", Password: hashedPass}, nil)

		token := generateAlphanumericCode(100)
		payload := jwt_maker.PayloadDTO{
			Username: "user",
		}
		jwtMaker.EXPECT().NewToken(payload).Return(token, nil)

		refreshToken := generateAlphanumericCode(50)

		payloadJSON, _ := json.Marshal(payload)
		redisClient.EXPECT().Set(ctx, keyMock, string(payloadJSON), &authService.refreshTokenTTL).Return(nil)

		token1, refreshToken1, err := authService.SignIn(ctx, paramsWithEmail)
		assert.NoError(t, err)
		assert.Equal(t, len(token), len(token1))
		assert.Equal(t, len(refreshToken), len(refreshToken1))
	})

}

func Test_authService_SignIn_GetUserForLogin_UserNotFound_repofailed(t *testing.T) {
	t.Parallel()
	var (
		ctrl            = gomock.NewController(t)
		ctx             = context.TODO()
		paramsWithLogin = LoginParams{
			Identificator: "user",
			Password:      "123456Aa@",
		}
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)
	)

	userService.EXPECT().GetUserForLogin(ctx, paramsWithLogin.Identificator).Return(users.UserDTO{}, users.ErrUserNotFound)

	_, _, err := authService.SignIn(ctx, paramsWithLogin)
	assert.Error(t, err)
	assert.Equal(t, err, ErrInvalidCredentials)
}

func Test_authService_SignIn_GetUserForLogin_repofailed(t *testing.T) {
	t.Parallel()
	var (
		ctrl            = gomock.NewController(t)
		ctx             = context.TODO()
		paramsWithLogin = LoginParams{
			Identificator: "user",
			Password:      "123456Aa@",
		}
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)
		errorRepo   = errors.New("database exploded")
	)

	userService.EXPECT().GetUserForLogin(ctx, paramsWithLogin.Identificator).Return(users.UserDTO{}, errorRepo)

	_, _, err := authService.SignIn(ctx, paramsWithLogin)
	assert.Error(t, err)
	assert.Equal(t, err, errorRepo)
}

func Test_authService_SignIn_InvalidPassword(t *testing.T) {
	t.Parallel()
	var (
		ctrl            = gomock.NewController(t)
		ctx             = context.TODO()
		paramsWithLogin = LoginParams{
			Identificator: "user",
			Password:      "123456Aa@",
		}
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)

		hashedRightPassword, _ = hasher.BcryptHash("@@11AAaa")
	)

	userService.EXPECT().GetUserForLogin(ctx, paramsWithLogin.Identificator).Return(users.UserDTO{
		Username: "user",
		Password: hashedRightPassword}, nil)

	_, _, err := authService.SignIn(ctx, paramsWithLogin)
	assert.Error(t, err)
	assert.Equal(t, err, ErrInvalidCredentials)
}

func Test_authService_SignIn_NewToken_failed(t *testing.T) {
	t.Parallel()
	var (
		ctrl            = gomock.NewController(t)
		ctx             = context.TODO()
		paramsWithLogin = LoginParams{
			Identificator: "user",
			Password:      "123456Aa@",
		}
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)
		errorToken  = errors.New("jwt failed")

		hashedRightPassword, _ = hasher.BcryptHash("123456Aa@")
	)

	userService.EXPECT().GetUserForLogin(ctx, paramsWithLogin.Identificator).Return(users.UserDTO{
		Username: "user",
		Password: hashedRightPassword}, nil)

	jwtMaker.EXPECT().NewToken(jwt_maker.PayloadDTO{
		Username: "user",
	}).Return("", errorToken)

	_, _, err := authService.SignIn(ctx, paramsWithLogin)
	assert.Error(t, err)
	assert.Equal(t, err, errorToken)
}

func Test_authService_SignIn_redisSet_failed(t *testing.T) {
	t.Parallel()
	var (
		ctrl            = gomock.NewController(t)
		ctx             = context.TODO()
		paramsWithLogin = LoginParams{
			Identificator: "user",
			Password:      "123456Aa@",
		}
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)
		redisClient = NewMockredisClient(ctrl)
		errorRedis  = errors.New("redis connection failed")

		hashedRightPassword, _ = hasher.BcryptHash("123456Aa@")
	)
	authService.WithRedisClient(redisClient)
	authService.SetRefreshTokenTTL(24 * time.Hour)

	userService.EXPECT().GetUserForLogin(ctx, paramsWithLogin.Identificator).Return(users.UserDTO{
		Username: "user",
		Password: hashedRightPassword}, nil)

	jwtMaker.EXPECT().NewToken(jwt_maker.PayloadDTO{
		Username: "user",
	}).Return("", nil)

	redisClient.EXPECT().Set(ctx, gomock.Any(), gomock.Any(), &authService.refreshTokenTTL).Return(errorRedis)

	_, _, err := authService.SignIn(ctx, paramsWithLogin)
	assert.Error(t, err)
	assert.Equal(t, err, errorRedis)
}

func Test_authService_SignIn_redisSet_noRefreshTokenTTL(t *testing.T) {
	t.Parallel()
	var (
		ctrl            = gomock.NewController(t)
		ctx             = context.TODO()
		paramsWithLogin = LoginParams{
			Identificator: "user",
			Password:      "123456Aa@",
		}
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)
		redisClient = NewMockredisClient(ctrl)

		hashedRightPassword, _ = hasher.BcryptHash("123456Aa@")
	)
	authService.WithRedisClient(redisClient)

	userService.EXPECT().GetUserForLogin(ctx, paramsWithLogin.Identificator).Return(users.UserDTO{
		Username: "user",
		Password: hashedRightPassword}, nil)

	jwtMaker.EXPECT().NewToken(jwt_maker.PayloadDTO{
		Username: "user",
	}).Return("", nil)

	redisClient.EXPECT().Set(ctx, gomock.Any(), gomock.Any(), &authService.refreshTokenTTL).Return(nil)

	_, _, err := authService.SignIn(ctx, paramsWithLogin)
	assert.NoError(t, err)
	assert.Equal(t, err, nil)
}

func Test_authService_SignIn_redisSet_noRedisClient(t *testing.T) {
	t.Parallel()
	var (
		ctrl            = gomock.NewController(t)
		ctx             = context.TODO()
		paramsWithLogin = LoginParams{
			Identificator: "user",
			Password:      "123456Aa@",
		}
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)

		hashedRightPassword, _ = hasher.BcryptHash("123456Aa@")
	)

	userService.EXPECT().GetUserForLogin(ctx, paramsWithLogin.Identificator).Return(users.UserDTO{
		Username: "user",
		Password: hashedRightPassword}, nil)

	jwtMaker.EXPECT().NewToken(jwt_maker.PayloadDTO{
		Username: "user",
	}).Return("", nil)

	_, r, err := authService.SignIn(ctx, paramsWithLogin)
	assert.NoError(t, err)
	assert.Equal(t, r, "")
}

func Test_authService_Register_success(t *testing.T) {
	t.Parallel()
	var (
		ctrl            = gomock.NewController(t)
		ctx             = context.TODO()
		paramsWithLogin = RegisterParams{
			DTO: users.CreateUserDTO{
				Username: "user",
				Email:    "user@user.ru",
				Password: "Pass@123456",
			},
		}
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)

		newID = "ruweioruweioruweororeurioewurweio"
	)

	userService.EXPECT().CreateUser(ctx, paramsWithLogin.DTO).Return(newID, nil)
	verifRepo.EXPECT().createVerificationCode(ctx, gomock.Any()).Return(nil)

	id, err := authService.Register(ctx, paramsWithLogin)
	assert.NoError(t, err)
	assert.Equal(t, id, newID)
}

func Test_authService_Register_createUser_failed(t *testing.T) {
	t.Parallel()
	var (
		ctrl            = gomock.NewController(t)
		ctx             = context.TODO()
		paramsWithLogin = RegisterParams{
			DTO: users.CreateUserDTO{
				Username: "user",
				Email:    "user@user.ru",
				Password: "Pass@123456",
			},
		}
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)

		errCreateUser = users.ErrUsernameExists
	)

	userService.EXPECT().CreateUser(ctx, paramsWithLogin.DTO).Return("", errCreateUser)

	res, err := authService.Register(ctx, paramsWithLogin)

	assert.Error(t, err)
	assert.Equal(t, err, errCreateUser)
	assert.Equal(t, res, "")
}

func Test_authService_Register_createVerificationCode_failed(t *testing.T) {
	t.Parallel()
	var (
		ctrl            = gomock.NewController(t)
		ctx             = context.TODO()
		paramsWithLogin = RegisterParams{
			DTO: users.CreateUserDTO{
				Username: "user",
				Email:    "user@user.ru",
				Password: "Pass@123456",
			},
		}
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)

		errCode = errors.New("verification database insert error")
	)

	userService.EXPECT().CreateUser(ctx, paramsWithLogin.DTO).Return(paramsWithLogin.DTO.Username, nil)
	verifRepo.EXPECT().createVerificationCode(ctx, gomock.Any()).Return(errCode)

	res, err := authService.Register(ctx, paramsWithLogin)

	assert.Error(t, err)
	assert.Equal(t, err, errCode)
	assert.Equal(t, res, "")
}

func Test_authService_VerifyUser_success(t *testing.T) {
	t.Parallel()
	var (
		ctx  = context.TODO()
		ctrl = gomock.NewController(t)
		req  = struct {
			username string
			code     string
		}{
			"user",
			"test_code",
		}
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)
	)

	verifRepo.EXPECT().getVerificationCode(ctx, req.username).Return(Verification{
		Username:  req.username,
		Email:     "user@user.com",
		Code:      "test_code",
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil)
	userService.EXPECT().EnableUser(ctx, req.username).Return(nil)

	err := authService.VerifyUser(ctx, req.username, req.code)

	assert.NoError(t, err)
}

func Test_authService_VerifyUser_getVerificationCode_repoFailed(t *testing.T) {
	t.Parallel()
	var (
		ctx  = context.TODO()
		ctrl = gomock.NewController(t)
		req  = struct {
			username string
			code     string
		}{
			"user",
			"test_code",
		}
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)

		errVerifRepo = ErrVerificationNotFound
	)

	verifRepo.EXPECT().getVerificationCode(ctx, req.username).Return(Verification{}, errVerifRepo)

	err := authService.VerifyUser(ctx, req.username, req.code)

	assert.Error(t, err)
	assert.Equal(t, err, errVerifRepo)
}

func Test_authService_VerifyUser_verificationExpired(t *testing.T) {
	t.Parallel()
	var (
		ctx  = context.TODO()
		ctrl = gomock.NewController(t)
		req  = struct {
			username string
			code     string
		}{
			"user",
			"test_code",
		}
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)

		errExpired = ErrVerificationExpired
	)

	verifRepo.EXPECT().getVerificationCode(ctx, req.username).Return(Verification{
		Username:  req.username,
		Email:     "user@user.com",
		Code:      "test_code",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}, nil)

	err := authService.VerifyUser(ctx, req.username, req.code)

	assert.Error(t, err)
	assert.Equal(t, err, errExpired)
}

func Test_authService_VerifyUser_invalidCode(t *testing.T) {
	t.Parallel()
	var (
		ctx  = context.TODO()
		ctrl = gomock.NewController(t)
		req  = struct {
			username string
			code     string
		}{
			"user",
			"test_code",
		}
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)

		errVerifInvalid = ErrVerificationInvalid
	)

	verifRepo.EXPECT().getVerificationCode(ctx, req.username).Return(Verification{
		Username:  req.username,
		Email:     "user@user.com",
		Code:      "test_code_not_like_in_request",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}, nil)

	err := authService.VerifyUser(ctx, req.username, req.code)

	assert.Error(t, err)
	assert.Equal(t, err, errVerifInvalid)
}

func Test_authService_VerifyUser_EnableUser_error(t *testing.T) {
	t.Parallel()
	var (
		ctx  = context.TODO()
		ctrl = gomock.NewController(t)
		req  = struct {
			username string
			code     string
		}{
			"user",
			"test_code",
		}
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)

		errEnable = errors.New("database exploded while query")
	)

	verifRepo.EXPECT().getVerificationCode(ctx, req.username).Return(Verification{
		Username:  req.username,
		Email:     "user@user.com",
		Code:      "test_code",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}, nil)

	userService.EXPECT().EnableUser(ctx, req.username).Return(errEnable)

	err := authService.VerifyUser(ctx, req.username, req.code)

	assert.Error(t, err)
	assert.Equal(t, err, errEnable)
}

// for RefreshToken
var (
	payloadUnmarshaled = jwt_maker.PayloadDTO{
		Username:  "user",
		Firstname: "first",
		Lastname:  "last",
		IsAdmin:   true,
	}
	payloadMarshaled, _ = json.Marshal(payloadUnmarshaled)
)

func Test_authService_RefreshToken_success(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.TODO()
		ctrl        = gomock.NewController(t)
		req         = strings.Repeat("token", 10)
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		redisClient = NewMockredisClient(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)

		keyMock = refreshTokenMathcer{prefix: "auth:refresh:", randomLength: 50}
	)
	authService.WithRedisClient(redisClient)
	authService.SetRefreshTokenTTL(time.Hour)

	redisClient.EXPECT().Get(ctx, keyMock).Return(string(payloadMarshaled), nil)
	redisClient.EXPECT().Delete(ctx, keyMock).Return(nil)
	jwtMaker.EXPECT().NewToken(payloadUnmarshaled).Return("access", nil)
	redisClient.EXPECT().Set(ctx, keyMock, string(payloadMarshaled), &authService.refreshTokenTTL).Return(nil)

	token, refreshToken, err := authService.RefreshToken(ctx, req)

	assert.NoError(t, err)
	assert.NotEqual(t, refreshToken, "")
	assert.Equal(t, token, "access")
}

func Test_authService_RefreshToken_noRedis(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.TODO()
		ctrl        = gomock.NewController(t)
		req         = strings.Repeat("token", 10)
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)
	)

	_, _, err := authService.RefreshToken(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "redis client is not initialized")
}

func Test_authService_RefreshToken_redis_Get_tokenNotFound(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.TODO()
		ctrl        = gomock.NewController(t)
		req         = strings.Repeat("token", 10)
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		redisClient = NewMockredisClient(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)

		keyMock = refreshTokenMathcer{prefix: "auth:refresh:", randomLength: 50}
	)
	authService.WithRedisClient(redisClient)
	authService.SetRefreshTokenTTL(time.Hour)

	redisClient.EXPECT().Get(ctx, keyMock).Return("", redis.Nil)

	_, _, err := authService.RefreshToken(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, err, ErrRefreshTokenNotFound)
}

func Test_authService_RefreshToken_redis_Get_someError(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.TODO()
		ctrl        = gomock.NewController(t)
		req         = strings.Repeat("token", 10)
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		redisClient = NewMockredisClient(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)
		someError   = errors.New("some")

		keyMock = refreshTokenMathcer{prefix: "auth:refresh:", randomLength: 50}
	)
	authService.WithRedisClient(redisClient)
	authService.SetRefreshTokenTTL(time.Hour)

	redisClient.EXPECT().Get(ctx, keyMock).Return("", someError)

	_, _, err := authService.RefreshToken(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, err, someError)
}

func Test_authService_RefreshToken_Unmarshal_error(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.TODO()
		ctrl        = gomock.NewController(t)
		req         = strings.Repeat("token", 10)
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		redisClient = NewMockredisClient(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)

		keyMock = refreshTokenMathcer{prefix: "auth:refresh:", randomLength: 50}
	)
	authService.WithRedisClient(redisClient)
	authService.SetRefreshTokenTTL(time.Hour)

	redisClient.EXPECT().Get(ctx, keyMock).Return(string(payloadMarshaled)+"for err", nil)

	_, _, err := authService.RefreshToken(ctx, req)
	assert.Error(t, err)
}

func Test_authService_RefreshToken_Delete_error(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.TODO()
		ctrl        = gomock.NewController(t)
		req         = strings.Repeat("token", 10)
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		redisClient = NewMockredisClient(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)
		someError   = errors.New("some error")

		keyMock = refreshTokenMathcer{prefix: "auth:refresh:", randomLength: 50}
	)
	authService.WithRedisClient(redisClient)
	authService.SetRefreshTokenTTL(time.Hour)

	redisClient.EXPECT().Get(ctx, keyMock).Return(string(payloadMarshaled), nil)
	redisClient.EXPECT().Delete(ctx, keyMock).Return(someError)
	jwtMaker.EXPECT().NewToken(payloadUnmarshaled).Return("access", nil)
	redisClient.EXPECT().Set(ctx, keyMock, string(payloadMarshaled), &authService.refreshTokenTTL).Return(nil)

	_, _, err := authService.RefreshToken(ctx, req)
	assert.NoError(t, err)
}

func Test_authService_RefreshToken_NewToken_error(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.TODO()
		ctrl        = gomock.NewController(t)
		req         = strings.Repeat("token", 10)
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		redisClient = NewMockredisClient(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)
		someError   = errors.New("some error")

		keyMock = refreshTokenMathcer{prefix: "auth:refresh:", randomLength: 50}
	)
	authService.WithRedisClient(redisClient)
	authService.SetRefreshTokenTTL(time.Hour)

	redisClient.EXPECT().Get(ctx, keyMock).Return(string(payloadMarshaled), nil)
	redisClient.EXPECT().Delete(ctx, keyMock).Return(someError)
	jwtMaker.EXPECT().NewToken(payloadUnmarshaled).Return("", someError)
	redisClient.EXPECT().Set(ctx, keyMock, string(payloadMarshaled), &authService.refreshTokenTTL).Return(nil).Times(0)

	_, _, err := authService.RefreshToken(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, err, someError)
}

func Test_authService_RefreshToken_redis_Set_error(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.TODO()
		ctrl        = gomock.NewController(t)
		req         = strings.Repeat("token", 10)
		userService = NewMockuserService(ctrl)
		jwtMaker    = NewMockjwtMaker(ctrl)
		verifRepo   = NewMockverificationCodeRepository(ctrl)
		redisClient = NewMockredisClient(ctrl)
		authService = NewAuthService(userService, jwtMaker, verifRepo)
		someError   = errors.New("some error")

		keyMock = refreshTokenMathcer{prefix: "auth:refresh:", randomLength: 50}
	)
	authService.WithRedisClient(redisClient)
	authService.SetRefreshTokenTTL(time.Hour)

	redisClient.EXPECT().Get(ctx, keyMock).Return(string(payloadMarshaled), nil)
	redisClient.EXPECT().Delete(ctx, keyMock).Return(someError)
	jwtMaker.EXPECT().NewToken(payloadUnmarshaled).Return("access", nil)
	redisClient.EXPECT().Set(ctx, keyMock, string(payloadMarshaled), &authService.refreshTokenTTL).Return(someError)

	_, _, err := authService.RefreshToken(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, err, someError)
}
