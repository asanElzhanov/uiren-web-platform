package data

import (
	"context"
	"encoding/json"
	"time"
	"uiren/internal/app/exercises"
	"uiren/internal/app/lessons"
	"uiren/internal/app/modules"
	"uiren/internal/app/progress"
	"uiren/internal/app/users"
	"uiren/pkg/logger"

	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source service.go -destination service_mock.go -package data

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

type lessonsService interface {
	GetLesson(ctx context.Context, code string) (lessons.LessonDTO, error)
}

type exerciseService interface {
	GetExercise(ctx context.Context, code string) (exercises.Exercise, error)
}

type redisClient interface {
	Set(ctx context.Context, key string, value interface{}, ttl *time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type progressService interface {
	GetXPLeaderboard(ctx context.Context, limit int) (progress.XPLeaderboard, error)
}

type DataService struct {
	redisClient        redisClient
	userService        userService
	modulesService     modulesService
	lessonsService     lessonsService
	exerciseService    exerciseService
	progressService    progressService
	dataTTL            time.Duration
	xpLeaderboardLimit int
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

func (s *DataService) WithProgressService(progressService progressService, xpLeaderboardLimit int) {
	s.progressService = progressService

	if xpLeaderboardLimit <= 0 {
		xpLeaderboardLimit = 20
	}
	s.xpLeaderboardLimit = xpLeaderboardLimit
}

func (s *DataService) WithLessonService(lessonsService lessonsService) {
	s.lessonsService = lessonsService
}

func (s *DataService) WithExerciseService(exerciseService exerciseService) {
	s.exerciseService = exerciseService
}

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
		CreatedAt: userDTO.CreatedAt,
	}, nil
}

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
		CreatedAt: userDTO.CreatedAt,
	}, nil
}

func (s *DataService) GetPublicModules(ctx context.Context) (ModulesList, error) {
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

func (s *DataService) GetXPLeaderboard(ctx context.Context) (XPLeaderboard, error) {
	logger.Info("DataService.GetXPLeaderboard new request")
	var leaderboard progress.XPLeaderboard

	key := generateXpLeaderboardKey(s.xpLeaderboardLimit)
	data, err := s.redisClient.Get(ctx, key)
	if err != nil {
		if err != redis.Nil {
			logger.Error("DataService.GetXPLeaderboard redisClient.Get: ", err)
			return XPLeaderboard{}, err
		}

		leaderboard, err = s.progressService.GetXPLeaderboard(ctx, s.xpLeaderboardLimit)
		if err != nil {
			logger.Error("DataService.GetXPLeaderboard progressService.GetXPLeaderboard: ", err)
			return XPLeaderboard{}, err
		}

		newData, err := json.Marshal(leaderboard)
		if err != nil {
			logger.Error("DataService.GetXPLeaderboard json.Marshal: ", err)
			return XPLeaderboard{}, err
		}

		err = s.redisClient.Set(ctx, key, newData, &s.dataTTL)
		if err != nil {
			logger.Error("DataService.GetXPLeaderboard redisClient.Set: ", err)
		}

		return XPLeaderboard{Board: leaderboard}, nil
	}

	err = json.Unmarshal([]byte(data), &leaderboard)
	if err != nil {
		logger.Error("DataService.GetXPLeaderboard json.Unmarshal: ", err)
		return XPLeaderboard{}, err
	}

	return XPLeaderboard{Board: leaderboard}, nil
}

func (s *DataService) GetPublicLesson(ctx context.Context, code string) (lessons.LessonDTO, error) {
	logger.Info("DataService.GetPublicLesson new request")
	var lesson lessons.LessonDTO

	key := generateLessonKey(code)
	data, err := s.redisClient.Get(ctx, key)

	if err != nil {
		if err != redis.Nil {
			logger.Error("DataService.GetPublicLesson redisClient.Get: ", err)
			return lessons.LessonDTO{}, err
		}

		lesson, err = s.lessonsService.GetLesson(ctx, code)
		if err != nil {
			logger.Error("DataService.GetPublicLesson lessonsService.GetLesson: ", err)
			return lessons.LessonDTO{}, err
		}

		newData, err := json.Marshal(lesson)
		if err != nil {
			logger.Error("DataService.GetPublicLesson json.Marshal: ", err)
			return lessons.LessonDTO{}, err
		}

		err = s.redisClient.Set(ctx, key, newData, &s.dataTTL)
		if err != nil {
			logger.Error("DataService.GetPublicLesson redisClient.Set: ", err)
		}

		return lesson, nil
	}

	err = json.Unmarshal([]byte(data), &lesson)
	if err != nil {
		logger.Error("DataService.GetPublicLesson json.Unmarshal: ", err)
		return lessons.LessonDTO{}, err
	}

	return lesson, nil
}

func (s *DataService) GetPublicExercise(ctx context.Context, code string) (exercises.Exercise, error) {
	logger.Info("DataService.GetPublicExercise new requst")
	var exercise exercises.Exercise

	key := generateExerciseKey(code)
	data, err := s.redisClient.Get(ctx, key)

	if err != nil {
		if err != redis.Nil {
			logger.Error("DataService.GetPublicExercise redisClient.Get: ", err)
			return exercises.Exercise{}, err
		}

		exercise, err = s.exerciseService.GetExercise(ctx, code)
		if err != nil {
			logger.Error("DataService.GetPublicExercise lessonsService.GetLesson: ", err)
			return exercises.Exercise{}, err
		}

		newData, err := json.Marshal(exercise)
		if err != nil {
			logger.Error("DataService.GetPublicExercise json.Marshal: ", err)
			return exercises.Exercise{}, err
		}

		err = s.redisClient.Set(ctx, key, newData, &s.dataTTL)
		if err != nil {
			logger.Error("DataService.GetPublicExercise redisClient.Set: ", err)
		}

		return exercise, nil
	}

	err = json.Unmarshal([]byte(data), &exercise)
	if err != nil {
		logger.Error("DataService.GetPublicExercise json.Unmarshal: ", err)
		return exercises.Exercise{}, err
	}

	return exercise, nil
}
