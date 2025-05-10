package data

import (
	"context"
	"encoding/json"
	"time"
	"uiren/internal/app/modules"
	"uiren/internal/app/users"
	"uiren/pkg/logger"

	"github.com/redis/go-redis/v9"
)

const (
	getModulesCacheKey = "all_modules_list"
)

type userService interface {
	GetUserByUsername(ctx context.Context, username string) (users.UserDTO, error)
	GetUserProgress(ctx context.Context, id string) (users.UserProgress, error)
}

type modulesService interface {
	GetModulesList(ctx context.Context) ([]modules.Module, error)
}

type redisClient interface {
	Set(ctx context.Context, key string, value interface{}, ttl *time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type DataService struct {
	redisClient    redisClient
	userService    userService
	modulesService modulesService
	dataTTL        time.Duration
}

func NewDataService(
	redisClient redisClient,
	userService userService,
	modulesService modulesService,
	dataTTL time.Duration,
) *DataService {
	return &DataService{
		redisClient:    redisClient,
		userService:    userService,
		modulesService: modulesService,
		dataTTL:        dataTTL,
	}
}

// todo: write tests
func (s *DataService) GetUserWithProgress(ctx context.Context, username string) (UserInfo, error) {
	logger.Info("DataService.GetUser new request")

	userDTO, err := s.userService.GetUserByUsername(ctx, username)
	if err != nil {
		logger.Error("DataService.GetUser userService.GetUserByUsername: ", err)
		return UserInfo{}, err
	}

	userProgress, err := s.userService.GetUserProgress(ctx, userDTO.ID)
	if err != nil {
		logger.Error("DataService.GetUser userService.GetUserProgress: ", err)
		return UserInfo{}, err
	}

	return UserInfo{
		ID:        userDTO.ID,
		Username:  userDTO.Username,
		Firstname: userDTO.Firstname,
		Lastname:  userDTO.Lastname,
		Email:     userDTO.Email,
		Phone:     userDTO.Phone,
		Progress:  &userProgress,
	}, nil
}

// todo: write tests
func (s *DataService) GetUserWithoutProgress(ctx context.Context, username string) (UserInfo, error) {
	logger.Info("DataService.GetUserWithoutProgress new request")

	userDTO, err := s.userService.GetUserByUsername(ctx, username)
	if err != nil {
		logger.Error("DataService.GetUserWithoutProgress userService.GetUserByUsername: ", err)
		return UserInfo{}, err
	}

	return UserInfo{
		ID:        userDTO.ID,
		Username:  userDTO.Username,
		Firstname: userDTO.Firstname,
		Lastname:  userDTO.Lastname,
		Email:     userDTO.Email,
		Phone:     userDTO.Phone,
	}, nil
}

// todo: write tests
func (s *DataService) GetModules(ctx context.Context) (ModulesList, error) {
	logger.Info("DataService.GetModules new request")
	var modulesList []modules.Module
	data, err := s.redisClient.Get(ctx, getModulesCacheKey)

	if err != nil {
		if err != redis.Nil {
			logger.Error("DataService.GetModules redisClient.Get: ", err)
			return ModulesList{}, err
		}

		modulesList, err = s.modulesService.GetModulesList(ctx)
		if err != nil {
			logger.Error("DataService.GetModules modulesService.GetModules: ", err)
			return ModulesList{}, err
		}
		newData, err := json.Marshal(modulesList)
		if err != nil {
			logger.Error("DataService.GetModules json.Marshal: ", err)
			return ModulesList{}, err
		}

		err = s.redisClient.Set(ctx, getModulesCacheKey, newData, &s.dataTTL)
		if err != nil {
			logger.Error("DataService.GetModules redisClient.Set: ", err)
			return ModulesList{}, err
		}
		return ModulesList{
			Modules: modulesList,
			Total:   len(modulesList),
		}, nil
	}

	err = json.Unmarshal([]byte(data), &modulesList)
	if err != nil {
		logger.Error("DataService.GetModules json.Unmarshal: ", err)
		return ModulesList{}, err
	}
	return ModulesList{
		Modules: modulesList,
		Total:   len(modulesList),
	}, nil
}
