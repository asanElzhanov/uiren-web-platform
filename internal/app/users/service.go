package users

import (
	"context"
	"errors"
	"uiren/internal/infrastracture/hasher"
	"uiren/pkg/logger"
)

type repository interface {
	createUser(ctx context.Context, params CreateUserDTO) (string, error)
	getUserByUsername(ctx context.Context, username string) (UserDTO, error)
	getUserByEmail(ctx context.Context, email string) (UserDTO, error)
	updateUser(ctx context.Context, dto UpdateUserDTO) (UserDTO, error)
	enableUser(ctx context.Context, username string) error
}

type UserService struct {
	repo repository
}

func NewUserService(repo repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, params CreateUserDTO) (string, error) {
	logger.Info("UserService.CreateUser new request")

	if err := params.Validate(); err != nil {
		logger.Error("UserService.CreateUser Validate error: ", err)
		return "", err
	}

	hashedPassword, err := hasher.BcryptHash(params.Password)
	if err != nil {
		logger.Error("UserService.CreateUser BcryptHash error: ", err)
		return "", err
	}
	params.Password = hashedPassword

	id, err := s.repo.createUser(ctx, params)
	if err != nil {
		logger.Error("UserService.CreateUser repo.CreateUser error: ", err)
		return "", err
	}

	return id, nil
}

func (s *UserService) GetUserForLogin(ctx context.Context, indetifier string) (UserDTO, error) {
	logger.Info("UserService.GetUserForLogin new request")

	getUser := func(identifier string) (UserDTO, error) {
		if user, err := s.repo.getUserByUsername(ctx, identifier); err == nil {
			return user, nil
		} else if !errors.Is(err, ErrUserNotFound) {
			return UserDTO{}, err
		}

		return s.repo.getUserByEmail(ctx, identifier)
	}

	user, err := getUser(indetifier)
	if err != nil {
		logger.Error("UserService.GetUserForLogin getUser: ", err)
		return UserDTO{}, err
	}

	user.normalize()
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, dto UpdateUserDTO) (UserDTO, error) {
	logger.Info("UserService.UpdateUser newRequest")

	if err := dto.Validate(); err != nil {
		logger.Error("UserService.UpdateUser Validate: ", err)
		return UserDTO{}, err
	}

	updatedUser, err := s.repo.updateUser(ctx, dto)
	if err != nil {
		logger.Error("UserService.UpdateUser updateUser: ", err)
		return UserDTO{}, nil
	}

	return updatedUser, nil
}

func (s *UserService) EnableUser(ctx context.Context, username string) error {
	logger.Info("UserService.EnableUser newRequest")

	return s.repo.enableUser(ctx, username)
}
