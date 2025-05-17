package data

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
	"uiren/internal/app/exercises"
	"uiren/internal/app/lessons"
	"uiren/internal/app/modules"
	"uiren/internal/app/progress"
	"uiren/internal/app/users"
	"uiren/pkg/logger"

	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func init() {
	logger.InitLogger("debug")
}

func Test_dataService_GetUser(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.TODO()
		ctrl        = gomock.NewController(t)
		userService = NewMockuserService(ctrl)
		service     = &DataService{userService: userService}
		repoReturn  = users.UserDTO{
			ID:        "user-123",
			Username:  "testuser",
			Email:     "testuser@example.com",
			Password:  "hashed_password_123",
			Firstname: "Test",
			Lastname:  "User",
			Phone:     "+1234567890",
			IsActive:  true,
			IsAdmin:   false,
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
			DeletedAt: time.Time{},
		}
		progress = users.UserProgress{
			Badges:       []string{"starter", "explorer"},
			XP:           1500,
			Achievements: nil,
		}
		errRepo = errors.New("erere")
	)

	t.Run("(with progress) success", func(t *testing.T) {
		userService.EXPECT().GetUserByUsername(ctx, repoReturn.Username).Return(repoReturn, nil)
		userService.EXPECT().GetUserProgress(ctx, repoReturn.ID).Return(progress, nil)
		result, err := service.GetUserWithProgress(ctx, repoReturn.Username)
		assert.NoError(t, err)
		assert.Equal(t, UserInfo{
			ID:        repoReturn.ID,
			Username:  repoReturn.Username,
			Firstname: repoReturn.Firstname,
			Lastname:  repoReturn.Lastname,
			Email:     repoReturn.Email,
			Phone:     repoReturn.Phone,
			Progress:  &progress,
		}, result)
	})
	t.Run("(with progress) repo failed#1", func(t *testing.T) {
		userService.EXPECT().GetUserByUsername(ctx, repoReturn.Username).Return(users.UserDTO{}, errRepo)
		result, err := service.GetUserWithProgress(ctx, repoReturn.Username)
		assert.Equal(t, err, errRepo)
		assert.Equal(t, UserInfo{}, result)
	})
	t.Run("(with progress) repo failed#2", func(t *testing.T) {
		userService.EXPECT().GetUserByUsername(ctx, repoReturn.Username).Return(repoReturn, nil)
		userService.EXPECT().GetUserProgress(ctx, repoReturn.ID).Return(users.UserProgress{}, errRepo)
		result, err := service.GetUserWithProgress(ctx, repoReturn.Username)
		assert.Equal(t, err, errRepo)
		assert.Equal(t, UserInfo{}, result)
	})
	t.Run("(no progress) success", func(t *testing.T) {
		userService.EXPECT().GetUserByUsername(ctx, repoReturn.Username).Return(repoReturn, nil)
		result, err := service.GetUserWithoutProgress(ctx, repoReturn.Username)
		assert.NoError(t, err)
		assert.Equal(t, UserInfo{
			ID:        repoReturn.ID,
			Username:  repoReturn.Username,
			Firstname: repoReturn.Firstname,
			Lastname:  repoReturn.Lastname,
			Email:     repoReturn.Email,
			Phone:     repoReturn.Phone,
		}, result)
	})
	t.Run("(no progress) repo failed", func(t *testing.T) {
		userService.EXPECT().GetUserByUsername(ctx, repoReturn.Username).Return(users.UserDTO{}, errRepo)
		result, err := service.GetUserWithoutProgress(ctx, repoReturn.Username)
		assert.Equal(t, err, errRepo)
		assert.Equal(t, UserInfo{}, result)
	})
}

func Test_dataService_GetPublicModules(t *testing.T) {
	t.Parallel()
	var (
		ctx            = context.TODO()
		ctrl           = gomock.NewController(t)
		modulesService = NewMockmodulesService(ctrl)
		redisCli       = NewMockredisClient(ctrl)
		service        = &DataService{modulesService: modulesService, redisClient: redisCli, dataTTL: time.Microsecond}
		errRepo        = errors.New("ere")
		returnRepo     = []modules.Module{
			{
				Code:        "mod-001",
				Title:       "Introduction to Kazakh",
				Description: "Learn basic Kazakh phrases and greetings.",
				Goal:        "Understand and use common phrases",
				Difficulty:  "Beginner",
				UnlockReq: modules.UnlockRequirements{
					PrevModuleCode: "",
					MinimumXP:      0,
				},
				Reward: modules.Reward{
					XP:    100,
					Badge: "starter",
				},
				Lessons:   []string{"lesson-001", "lesson-002"},
				CreatedAt: time.Unix(999999, 999),
				DeletedAt: nil,
			},
			{
				Code:        "mod-002",
				Title:       "Kazakh Alphabet",
				Description: "Learn how to read and write Kazakh letters.",
				Goal:        "Master the Kazakh alphabet",
				Difficulty:  "Beginner",
				UnlockReq: modules.UnlockRequirements{
					PrevModuleCode: "mod-001",
					MinimumXP:      50,
				},
				Reward: modules.Reward{
					XP:    150,
					Badge: "alphabet_master",
				},
				Lessons:   []string{"lesson-003", "lesson-004"},
				CreatedAt: time.Unix(999999, 999),
				DeletedAt: nil,
			},
			{
				Code:        "mod-003",
				Title:       "Daily Conversations",
				Description: "Practice dialogues used in daily life.",
				Goal:        "Hold basic conversations",
				Difficulty:  "Intermediate",
				UnlockReq: modules.UnlockRequirements{
					PrevModuleCode: "mod-002",
					MinimumXP:      200,
				},
				Reward: modules.Reward{
					XP:    200,
					Badge: "conversationalist",
				},
				Lessons:   []string{"lesson-005", "lesson-006"},
				CreatedAt: time.Unix(999999, 999),
				DeletedAt: nil,
			},
		}
	)

	t.Run("success(redis)", func(t *testing.T) {
		data, _ := json.Marshal(returnRepo)
		redisCli.EXPECT().Get(ctx, getModulesCacheKey).Return(string(data), nil)

		result, err := service.GetPublicModules(ctx)
		assert.NoError(t, err)
		assert.Equal(t, result, ModulesList{
			Modules: returnRepo,
			Total:   len(returnRepo),
		})
	})
	t.Run("success(db) redis-set no error", func(t *testing.T) {
		redisCli.EXPECT().Get(ctx, getModulesCacheKey).Return("", redis.Nil)
		modulesService.EXPECT().GetModulesList(ctx).Return(returnRepo, nil)
		data, _ := json.Marshal(returnRepo)
		redisCli.EXPECT().Set(ctx, getModulesCacheKey, data, &service.dataTTL).Return(nil)

		result, err := service.GetPublicModules(ctx)
		assert.NoError(t, err)
		assert.Equal(t, result, ModulesList{
			Modules: returnRepo,
			Total:   len(returnRepo),
		})
	})

	t.Run("success(db) redis-set error", func(t *testing.T) {
		redisCli.EXPECT().Get(ctx, getModulesCacheKey).Return("", redis.Nil)
		modulesService.EXPECT().GetModulesList(ctx).Return(returnRepo, nil)
		data, _ := json.Marshal(returnRepo)
		redisCli.EXPECT().Set(ctx, getModulesCacheKey, data, &service.dataTTL).Return(redis.Nil)

		result, err := service.GetPublicModules(ctx)
		assert.NoError(t, err)
		assert.Equal(t, result, ModulesList{
			Modules: returnRepo,
			Total:   len(returnRepo),
		})
	})

	t.Run("redis cli error not redis.Nil", func(t *testing.T) {
		redisCli.EXPECT().Get(ctx, getModulesCacheKey).Return("", errRepo)
		result, err := service.GetPublicModules(ctx)
		assert.Equal(t, err, errRepo)
		assert.Equal(t, result, ModulesList{})
	})

	t.Run("modulesService error", func(t *testing.T) {
		redisCli.EXPECT().Get(ctx, getModulesCacheKey).Return("", redis.Nil)
		modulesService.EXPECT().GetModulesList(ctx).Return(nil, errRepo)

		result, err := service.GetPublicModules(ctx)
		assert.Equal(t, err, errRepo)
		assert.Equal(t, result, ModulesList{})
	})
}

func Test_dataService_GetPublicLesson(t *testing.T) {
	t.Parallel()
	var (
		ctx           = context.TODO()
		ctrl          = gomock.NewController(t)
		lessonService = NewMocklessonsService(ctrl)
		redisCli      = NewMockredisClient(ctrl)
		mockTime      = time.Unix(123165, 156)
		service       = &DataService{lessonsService: lessonService, redisClient: redisCli, dataTTL: time.Microsecond}
		errRepo       = errors.New("ere")
		returnRepo    = lessons.LessonDTO{
			Code:        "lesson-001",
			Title:       "Basics of Kazakh",
			Description: "Learn basic greetings and introductions in Kazakh.",
			Exercises: []exercises.Exercise{
				{
					Code:          "ex-001",
					ExerciseType:  "multiple_choice",
					Question:      "How do you say 'Hello' in Kazakh?",
					Hints:         []string{"It's a common greeting", "Starts with 'S'"},
					Explanation:   "The word 'Сәлем' means 'Hello' in Kazakh.",
					Options:       []string{"Сәлем", "Қайырлы күн", "Сау бол"},
					CorrectAnswer: "Сәлем",
					CreatedAt:     mockTime,
					DeletedAt:     nil,
				},
				{
					Code:         "ex-002",
					ExerciseType: "order_words",
					Question:     "Arrange the words to form: 'My name is Ayan'",
					Hints:        []string{"Start with 'Менің'", "Ends with 'Аян'"},
					Explanation:  "In Kazakh, the correct order is: 'Менің атым Аян'",
					CorrectOrder: []string{"Менің", "атым", "Аян"},
					CreatedAt:    mockTime,
					DeletedAt:    nil,
				},
				{
					Code:         "ex-003",
					ExerciseType: "match_pairs",
					Question:     "Match the Kazakh words with their English meanings.",
					Hints:        []string{"Common everyday words"},
					Explanation:  "These are basic vocabulary words.",
					Pairs: []exercises.Pair{
						{Term: "Ит", Match: "Dog"},
						{Term: "Күн", Match: "Sun"},
						{Term: "Кітап", Match: "Book"},
					},
					CreatedAt: mockTime,
					DeletedAt: nil,
				},
			},
			CreatedAt: mockTime,
			DeletedAt: time.Time{},
		}
	)

	t.Run("success(redis)", func(t *testing.T) {
		data, _ := json.Marshal(returnRepo)
		redisCli.EXPECT().Get(ctx, generateLessonKey(returnRepo.Code)).Return(string(data), nil)

		result, err := service.GetPublicLesson(ctx, returnRepo.Code)
		assert.NoError(t, err)
		assert.Equal(t, result, returnRepo)
	})
	t.Run("success(db) redis-set no error", func(t *testing.T) {
		redisCli.EXPECT().Get(ctx, generateLessonKey(returnRepo.Code)).Return("", redis.Nil)
		lessonService.EXPECT().GetLesson(ctx, returnRepo.Code).Return(returnRepo, nil)
		data, _ := json.Marshal(returnRepo)
		redisCli.EXPECT().Set(ctx, generateLessonKey(returnRepo.Code), data, &service.dataTTL).Return(nil)

		result, err := service.GetPublicLesson(ctx, returnRepo.Code)
		assert.NoError(t, err)
		assert.Equal(t, result, returnRepo)
	})

	t.Run("success(db) redis-set error", func(t *testing.T) {
		redisCli.EXPECT().Get(ctx, generateLessonKey(returnRepo.Code)).Return("", redis.Nil)
		lessonService.EXPECT().GetLesson(ctx, returnRepo.Code).Return(returnRepo, nil)
		data, _ := json.Marshal(returnRepo)
		redisCli.EXPECT().Set(ctx, generateLessonKey(returnRepo.Code), data, &service.dataTTL).Return(redis.Nil)

		result, err := service.GetPublicLesson(ctx, returnRepo.Code)
		assert.NoError(t, err)
		assert.Equal(t, result, returnRepo)
	})

	t.Run("redis cli error not redis.Nil", func(t *testing.T) {
		redisCli.EXPECT().Get(ctx, generateLessonKey(returnRepo.Code)).Return("", errRepo)
		result, err := service.GetPublicLesson(ctx, returnRepo.Code)
		assert.Equal(t, err, errRepo)
		assert.Equal(t, result, lessons.LessonDTO{})
	})

	t.Run("lessonService error", func(t *testing.T) {
		redisCli.EXPECT().Get(ctx, generateLessonKey(returnRepo.Code)).Return("", redis.Nil)
		lessonService.EXPECT().GetLesson(ctx, returnRepo.Code).Return(lessons.LessonDTO{}, errRepo)

		result, err := service.GetPublicLesson(ctx, returnRepo.Code)
		assert.Equal(t, err, errRepo)
		assert.Equal(t, result, lessons.LessonDTO{})
	})
}

func Test_dataService_GetPublicExercise(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		exerciseService = NewMockexerciseService(ctrl)
		redisCli        = NewMockredisClient(ctrl)
		mockTime        = time.Unix(123165, 156)
		service         = &DataService{exerciseService: exerciseService, redisClient: redisCli, dataTTL: time.Microsecond}
		errRepo         = errors.New("ere")
		returnRepo      = exercises.Exercise{
			Code:          "ex-001",
			ExerciseType:  "multiple_choice",
			Question:      "How do you say 'Thank you' in Kazakh?",
			Hints:         []string{"It starts with 'Р'", "A common polite phrase"},
			Explanation:   "The word 'Рақмет' means 'Thank you' in Kazakh.",
			Options:       []string{"Сәлем", "Рақмет", "Кешіріңіз"},
			CorrectAnswer: "Рақмет",
			CreatedAt:     mockTime,
			DeletedAt:     nil,
		}
	)

	t.Run("success(redis)", func(t *testing.T) {
		data, _ := json.Marshal(returnRepo)
		redisCli.EXPECT().Get(ctx, generateExerciseKey(returnRepo.Code)).Return(string(data), nil)

		result, err := service.GetPublicExercise(ctx, returnRepo.Code)
		assert.NoError(t, err)
		assert.Equal(t, result, returnRepo)
	})
	t.Run("success(db) redis-set no error", func(t *testing.T) {
		redisCli.EXPECT().Get(ctx, generateExerciseKey(returnRepo.Code)).Return("", redis.Nil)
		exerciseService.EXPECT().GetExercise(ctx, returnRepo.Code).Return(returnRepo, nil)
		data, _ := json.Marshal(returnRepo)
		redisCli.EXPECT().Set(ctx, generateExerciseKey(returnRepo.Code), data, &service.dataTTL).Return(nil)

		result, err := service.GetPublicExercise(ctx, returnRepo.Code)
		assert.NoError(t, err)
		assert.Equal(t, result, returnRepo)
	})

	t.Run("success(db) redis-set error", func(t *testing.T) {
		redisCli.EXPECT().Get(ctx, generateExerciseKey(returnRepo.Code)).Return("", redis.Nil)
		exerciseService.EXPECT().GetExercise(ctx, returnRepo.Code).Return(returnRepo, nil)
		data, _ := json.Marshal(returnRepo)
		redisCli.EXPECT().Set(ctx, generateExerciseKey(returnRepo.Code), data, &service.dataTTL).Return(redis.Nil)

		result, err := service.GetPublicExercise(ctx, returnRepo.Code)
		assert.NoError(t, err)
		assert.Equal(t, result, returnRepo)
	})

	t.Run("redis cli error not redis.Nil", func(t *testing.T) {
		redisCli.EXPECT().Get(ctx, generateExerciseKey(returnRepo.Code)).Return("", errRepo)
		result, err := service.GetPublicExercise(ctx, returnRepo.Code)
		assert.Equal(t, err, errRepo)
		assert.Equal(t, result, exercises.Exercise{})
	})

	t.Run("exerciseService error", func(t *testing.T) {
		redisCli.EXPECT().Get(ctx, generateExerciseKey(returnRepo.Code)).Return("", redis.Nil)
		exerciseService.EXPECT().GetExercise(ctx, returnRepo.Code).Return(exercises.Exercise{}, errRepo)

		result, err := service.GetPublicExercise(ctx, returnRepo.Code)
		assert.Equal(t, err, errRepo)
		assert.Equal(t, result, exercises.Exercise{})
	})
}

func Test_dataService_GetXPLeaderboard(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		progressService = NewMockprogressService(ctrl)
		redisCli        = NewMockredisClient(ctrl)
		service         = &DataService{progressService: progressService, xpLeaderboardLimit: 200, redisClient: redisCli, dataTTL: time.Microsecond}
		errRepo         = errors.New("ere")
		returnRepo      = progress.XPLeaderboard{
			Leaders: []progress.XPLeaderboardEntry{
				{
					Rank:     1,
					UserID:   "user-001",
					Username: "admin_hero",
					XP:       1500,
				},
				{
					Rank:     2,
					UserID:   "user-002",
					Username: "kaz_learn",
					XP:       1200,
				},
				{
					Rank:     3,
					UserID:   "user-003",
					Username: "turk_master",
					XP:       1100,
				},
			},
			Total: 3,
		}
	)

	t.Run("success(redis)", func(t *testing.T) {
		data, _ := json.Marshal(returnRepo)
		redisCli.EXPECT().Get(ctx, generateXpLeaderboardKey(200)).Return(string(data), nil)

		result, err := service.GetXPLeaderboard(ctx)
		assert.NoError(t, err)
		assert.Equal(t, result, XPLeaderboard{
			Board: returnRepo,
		})
	})
	t.Run("success(db) redis-set no error", func(t *testing.T) {
		redisCli.EXPECT().Get(ctx, generateXpLeaderboardKey(200)).Return("", redis.Nil)
		progressService.EXPECT().GetXPLeaderboard(ctx, 200).Return(returnRepo, nil)
		data, _ := json.Marshal(returnRepo)
		redisCli.EXPECT().Set(ctx, generateXpLeaderboardKey(200), data, &service.dataTTL).Return(nil)

		result, err := service.GetXPLeaderboard(ctx)
		assert.NoError(t, err)
		assert.Equal(t, result, XPLeaderboard{Board: returnRepo})
	})

	t.Run("success(db) redis-set error", func(t *testing.T) {
		redisCli.EXPECT().Get(ctx, generateXpLeaderboardKey(200)).Return("", redis.Nil)
		progressService.EXPECT().GetXPLeaderboard(ctx, 200).Return(returnRepo, nil)
		data, _ := json.Marshal(returnRepo)
		redisCli.EXPECT().Set(ctx, generateXpLeaderboardKey(200), data, &service.dataTTL).Return(redis.Nil)

		result, err := service.GetXPLeaderboard(ctx)
		assert.NoError(t, err)
		assert.Equal(t, result, XPLeaderboard{Board: returnRepo})
	})

	t.Run("redis cli error not redis.Nil", func(t *testing.T) {
		redisCli.EXPECT().Get(ctx, generateXpLeaderboardKey(200)).Return("", errRepo)
		result, err := service.GetXPLeaderboard(ctx)
		assert.Equal(t, err, errRepo)
		assert.Equal(t, result, XPLeaderboard{})
	})

	t.Run("progressService error", func(t *testing.T) {
		redisCli.EXPECT().Get(ctx, generateXpLeaderboardKey(200)).Return("", redis.Nil)
		progressService.EXPECT().GetXPLeaderboard(ctx, 200).Return(progress.XPLeaderboard{}, errRepo)

		result, err := service.GetXPLeaderboard(ctx)
		assert.Equal(t, err, errRepo)
		assert.Equal(t, result, XPLeaderboard{})
	})
}
