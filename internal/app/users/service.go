package users

import (
	"context"
	"errors"
	"uiren/internal/app/progress"
	"uiren/internal/infrastracture/hasher"
	"uiren/pkg/logger"
)

//go:generate mockgen -source service.go -destination service_mock.go -package users

type repository interface {
	createUser(ctx context.Context, params CreateUserDTO) (string, error)
	getUserByUsername(ctx context.Context, username string) (UserDTO, error)
	getUserByEmail(ctx context.Context, email string) (UserDTO, error)
	updateUser(ctx context.Context, dto UpdateUserDTO) (UserDTO, error)
	enableUser(ctx context.Context, username string) error
	checkUserExists(ctx context.Context, username string) error
	getAllUsers(ctx context.Context) ([]UserDTO, error)
	getUserByID(ctx context.Context, id string) (UserDTO, error)
}

// temp
type ProgressService interface {
	GetBadges(ctx context.Context, id string) ([]string, error)
	GetXP(ctx context.Context, id string) (int, error)
	GetAchievements(ctx context.Context, id string) ([]progress.UserAchievement, error)
}

type UserService struct {
	repo       repository
	prgService ProgressService
}

func NewUserService(repo repository, prgService ProgressService) *UserService {
	return &UserService{
		repo:       repo,
		prgService: prgService,
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
		return UserDTO{}, err
	}

	return updatedUser, nil
}

func (s *UserService) EnableUser(ctx context.Context, username string) error {
	logger.Info("UserService.EnableUser newRequest")

	return s.repo.enableUser(ctx, username)
}

func (s *UserService) CheckUserExists(ctx context.Context, username string) error {
	logger.Info("UserService.CheckUserExists new request")

	return s.repo.checkUserExists(ctx, username)
}

// todo: write tests
func (s *UserService) GetUserProgress(ctx context.Context, id string) (UserProgress, error) {
	logger.Info("UserService.GetUserProgress new request")

	badges, err := s.prgService.GetBadges(ctx, id)
	if err != nil {
		logger.Error("UserService.GetUserProgress GetBadges: ", err)
		return UserProgress{}, err
	}

	xp, err := s.prgService.GetXP(ctx, id)
	if err != nil {
		logger.Error("UserService.GetUserProgress GetXP: ", err)
		return UserProgress{}, err
	}

	achievements, err := s.prgService.GetAchievements(ctx, id)
	if err != nil {
		logger.Error("UserService.GetUserProgress GetAchievements: ", err)
		return UserProgress{}, err
	}

	return UserProgress{
		Badges:       badges,
		XP:           int32(xp),
		Achievements: achievements,
	}, nil
}

// todo: write tests
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (UserDTO, error) {
	logger.Info("UserService.GetUserByUsername new request")

	user, err := s.repo.getUserByUsername(ctx, username)
	if err != nil {
		logger.Error("UserService.GetUserByUsername getUserByUsername: ", err)
		return UserDTO{}, err
	}

	user.normalize()
	return user, nil
}

// todo: write tests
func (s *UserService) GetAllUsers(ctx context.Context) ([]UserDTO, error) {
	logger.Info("UserService.GetUsers new request")

	users, err := s.repo.getAllUsers(ctx)
	if err != nil {
		logger.Error("UserService.GetUsers getUsers: ", err)
		return nil, err
	}

	for i := range users {
		users[i].normalize()
		users[i].Password = "..."
	}

	return users, nil
}

//todo: write tests

func (s *UserService) GetUserByID(ctx context.Context, id string) (UserDTO, error) {
	logger.Info("UserService.GetUserByID new request")
	user, err := s.repo.getUserByID(ctx, id)
	if err != nil {
		logger.Error("UserService.GetUserByID getUserByID: ", err)
		return UserDTO{}, err
	}
	user.normalize()
	return user, nil
}
