package users

import (
	"context"
	"errors"
	"testing"
	"uiren/internal/infrastracture/hasher"
	"uiren/pkg/logger"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func init() {
	logger.InitLogger("info")
}

func Test_userService_CreateUser_success(t *testing.T) {
	t.Parallel()
	var (
		ctx    = context.TODO()
		ctrl   = gomock.NewController(t)
		repo   = NewMockrepository(ctrl)
		srv    = NewUserService(repo, NewMockuserProgressRepo(ctrl))
		params = CreateUserDTO{
			Username: "user9871",
			Email:    "asan@gmail.com    ",
			Password: "Pass@123456",
		}
	)
	normalizedParams := params
	normalizedParams.normalize()
	pass, _ := hasher.BcryptHash(normalizedParams.Password)
	normalizedParams.Password = pass
	repo.EXPECT().createUser(ctx, createUserDTOMatcher{normalizedParams}).Return("uuid", nil)

	id, err := srv.CreateUser(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, id, "uuid")
}

func Test_userService_CreateUser_fail(t *testing.T) {
	t.Parallel()
	var (
		ctx  = context.TODO()
		ctrl = gomock.NewController(t)
		repo = NewMockrepository(ctrl)
		srv  = NewUserService(repo, NewMockuserProgressRepo(ctrl))
	)
	/*
		normalizedParams := params
		normalizedParams.normalize()
		pass, _ := hasher.BcryptHash(normalizedParams.Password)
		normalizedParams.Password = pass
	*/
	t.Run("invalid email", func(t *testing.T) {
		params := CreateUserDTO{
			Username: "user9871",
			Email:    "invalidemail",
			Password: "Pass@123456",
		}

		_, err := srv.CreateUser(ctx, params)
		assert.Error(t, err)
		assert.Equal(t, err, ErrIncorrectEmail)
	})

	t.Run("invalid password(no symbols)", func(t *testing.T) {
		params := CreateUserDTO{
			Username: "user9871",
			Email:    "a@a.ru",
			Password: "Pass123456",
		}

		_, err := srv.CreateUser(ctx, params)
		assert.Error(t, err)
		assert.Equal(t, err, ErrIncorrectPassword)
	})

	t.Run("invalid password(no numbers)", func(t *testing.T) {
		params := CreateUserDTO{
			Username: "user9871",
			Email:    "a@a.ru",
			Password: "Pass@afdsafasdf",
		}

		_, err := srv.CreateUser(ctx, params)
		assert.Error(t, err)
		assert.Equal(t, err, ErrIncorrectPassword)
	})

	t.Run("invalid password(no lower case letters)", func(t *testing.T) {
		params := CreateUserDTO{
			Username: "user9871",
			Email:    "a@a.ru",
			Password: "PASS@FDALFMDAK321312",
		}

		_, err := srv.CreateUser(ctx, params)
		assert.Error(t, err)
		assert.Equal(t, err, ErrIncorrectPassword)
	})

	t.Run("invalid password(no upper case letters)", func(t *testing.T) {
		params := CreateUserDTO{
			Username: "user9871",
			Email:    "a@a.ru",
			Password: "123@dsadas",
		}

		_, err := srv.CreateUser(ctx, params)
		assert.Error(t, err)
		assert.Equal(t, err, ErrIncorrectPassword)
	})

	t.Run("invalid password(less than 8 symbols)", func(t *testing.T) {
		params := CreateUserDTO{
			Username: "user9871",
			Email:    "a@a.ru",
			Password: "123@dsF",
		}

		_, err := srv.CreateUser(ctx, params)
		assert.Error(t, err)
		assert.Equal(t, err, ErrIncorrectPassword)
	})

	t.Run("invalid password(no letters)", func(t *testing.T) {
		params := CreateUserDTO{
			Username: "user9871",
			Email:    "a@a.ru",
			Password: "123@@@@432412",
		}

		_, err := srv.CreateUser(ctx, params)
		assert.Error(t, err)
		assert.Equal(t, err, ErrIncorrectPassword)
	})

	t.Run("createUser repo error", func(t *testing.T) {
		someErr := errors.New("some")
		params := CreateUserDTO{
			Username: "user9871",
			Email:    "a@a.ru",
			Password: "123@@@@432412fadfadsFDA",
		}
		normalizedParams := params
		normalizedParams.normalize()
		pass, _ := hasher.BcryptHash(normalizedParams.Password)
		normalizedParams.Password = pass
		repo.EXPECT().createUser(ctx, createUserDTOMatcher{normalizedParams}).Return("", someErr)

		id, err := srv.CreateUser(ctx, params)
		assert.Equal(t, "", id)
		assert.Equal(t, err, someErr)
	})
}

func Test_userService_UpdateUser_success(t *testing.T) {
	t.Parallel()
	var (
		ctx    = context.TODO()
		ctrl   = gomock.NewController(t)
		repo   = NewMockrepository(ctrl)
		srv    = NewUserService(repo, NewMockuserProgressRepo(ctrl))
		params = UpdateUserDTO{
			Firstname:   "first",
			Lastname:    "last",
			Phone:       "+77474434210",
			PhoneRegion: "KZ",
		}
	)

	repo.EXPECT().updateUser(ctx, params).Return(UserDTO{ID: params.ID}, nil)
	user, err := srv.UpdateUser(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, user, UserDTO{ID: params.ID})
}

func Test_userService_UpdateUser_fail(t *testing.T) {
	t.Parallel()
	var (
		ctx  = context.TODO()
		ctrl = gomock.NewController(t)
		repo = NewMockrepository(ctrl)
		srv  = NewUserService(repo, NewMockuserProgressRepo(ctrl))
	)

	t.Run("wrong phone number", func(t *testing.T) {
		params := UpdateUserDTO{
			Firstname:   "first",
			Lastname:    "last",
			Phone:       "+77474434210fds",
			PhoneRegion: "KZ",
		}
		_, err := srv.UpdateUser(ctx, params)
		assert.Equal(t, err, ErrIncorrectPhone)
	})

	t.Run("number does not match the region", func(t *testing.T) {
		params := UpdateUserDTO{
			Firstname:   "first",
			Lastname:    "last",
			Phone:       "+77474434210",
			PhoneRegion: "RU",
		}
		_, err := srv.UpdateUser(ctx, params)
		assert.Equal(t, err, ErrIncorrectPhone)
	})

	t.Run("repo failed", func(t *testing.T) {
		someErr := errors.New("some err")
		params := UpdateUserDTO{
			Firstname:   "first",
			Lastname:    "last",
			Phone:       "+77474434210",
			PhoneRegion: "KZ",
		}

		params.normalize()
		repo.EXPECT().updateUser(ctx, params).Return(UserDTO{}, someErr)
		_, err := srv.UpdateUser(ctx, params)
		assert.Equal(t, err, someErr)
	})
}

func Test_userService_GetUserForLogin_success(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.TODO()
		ctrl        = gomock.NewController(t)
		repo        = NewMockrepository(ctrl)
		srv         = NewUserService(repo, NewMockuserProgressRepo(ctrl))
		reqUsername = "username1"
		reqEmail    = "email@email.ru"
		repoResp    = UserDTO{
			ID:       "uuid",
			Username: "username1",
			Email:    "email@email.ru",
		}
	)

	t.Run("with username", func(t *testing.T) {
		repo.EXPECT().getUserByUsername(ctx, reqUsername).Return(repoResp, nil)
		repoResp.normalize()
		user, err := srv.GetUserForLogin(ctx, reqUsername)
		assert.NoError(t, err)
		assert.Equal(t, user, repoResp)
	})
	t.Run("with email", func(t *testing.T) {
		repo.EXPECT().getUserByUsername(ctx, reqEmail).Return(UserDTO{}, ErrUserNotFound)
		repo.EXPECT().getUserByEmail(ctx, reqEmail).Return(repoResp, nil)
		repoResp.normalize()
		user, err := srv.GetUserForLogin(ctx, reqEmail)
		assert.NoError(t, err)
		assert.Equal(t, user, repoResp)
	})
}

func Test_userService_GetUserForLogin_repoFailed(t *testing.T) {
	t.Parallel()
	var (
		ctx        = context.TODO()
		ctrl       = gomock.NewController(t)
		repo       = NewMockrepository(ctrl)
		srv        = NewUserService(repo, NewMockuserProgressRepo(ctrl))
		identifier = "user1"
		repoError  = errors.New("repo error")
	)

	t.Run("getUserByUsername", func(t *testing.T) {
		repo.EXPECT().getUserByUsername(ctx, identifier).Return(UserDTO{}, repoError)
		_, err := srv.GetUserForLogin(ctx, identifier)
		assert.Equal(t, err, repoError)
	})
	t.Run("getUserByEmail", func(t *testing.T) {
		repo.EXPECT().getUserByUsername(ctx, identifier).Return(UserDTO{}, ErrUserNotFound)
		repo.EXPECT().getUserByEmail(ctx, identifier).Return(UserDTO{}, repoError)
		_, err := srv.GetUserForLogin(ctx, identifier)
		assert.Equal(t, err, repoError)
	})
}

func Test_UserService_EnableUser_Success(t *testing.T) {
	t.Parallel()
	var (
		ctx      = context.TODO()
		ctrl     = gomock.NewController(t)
		repo     = NewMockrepository(ctrl)
		srv      = NewUserService(repo, NewMockuserProgressRepo(ctrl))
		username = "test_user"
	)

	repo.EXPECT().enableUser(ctx, username).Return(nil)

	err := srv.EnableUser(ctx, username)
	assert.NoError(t, err)
}

func Test_UserService_EnableUser_Fail(t *testing.T) {
	t.Parallel()
	var (
		ctx      = context.TODO()
		ctrl     = gomock.NewController(t)
		repo     = NewMockrepository(ctrl)
		srv      = NewUserService(repo, NewMockuserProgressRepo(ctrl))
		username = "test_user"
	)

	repo.EXPECT().enableUser(ctx, username).Return(ErrUserNotFound)

	err := srv.EnableUser(ctx, username)
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
}

func Test_UserService_CheckUserExists_Success(t *testing.T) {
	t.Parallel()
	var (
		ctx      = context.TODO()
		ctrl     = gomock.NewController(t)
		repo     = NewMockrepository(ctrl)
		srv      = NewUserService(repo, NewMockuserProgressRepo(ctrl))
		username = "test_user"
	)

	repo.EXPECT().checkUserExists(ctx, username).Return(nil)

	err := srv.CheckUserExists(ctx, username)
	assert.NoError(t, err)
}

func Test_UserService_CheckUserExists_Fail(t *testing.T) {
	t.Parallel()
	var (
		ctx      = context.TODO()
		ctrl     = gomock.NewController(t)
		repo     = NewMockrepository(ctrl)
		srv      = NewUserService(repo, NewMockuserProgressRepo(ctrl))
		username = "test_user"
	)

	repo.EXPECT().checkUserExists(ctx, username).Return(ErrUserNotFound)

	err := srv.CheckUserExists(ctx, username)
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
}
