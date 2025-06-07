package lessons

import (
	"context"
	"errors"
	"testing"
	"time"
	"uiren/internal/app/exercises"
	"uiren/pkg/logger"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	logger.InitLogger("info")
}

func Test_LessonsService_CreateLesson_success(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		dto              = CreateLessonDTO{
			Code:        "code",
			Title:       "title",
			Description: "desc",
		}
	)

	repo.EXPECT().createLesson(ctx, MatchCreateLessonDTO(dto)).Return(primitive.ObjectID{}, nil)

	id, err := srv.CreateLesson(ctx, dto)
	assert.NoError(t, err)
	assert.Equal(t, id, primitive.ObjectID{})
}

func Test_LessonsService_CreateLesson_fail(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		dto              = CreateLessonDTO{
			Code:        "code",
			Title:       "title",
			Description: "desc",
		}
	)

	repo.EXPECT().createLesson(ctx, MatchCreateLessonDTO(dto)).Return(primitive.ObjectID{}, ErrCodeAlreadyExists)

	id, err := srv.CreateLesson(ctx, dto)
	assert.Error(t, err)
	assert.Equal(t, err, ErrCodeAlreadyExists)
	assert.Equal(t, id, primitive.NilObjectID)
}

func Test_LessonsService_UpdateLesson_success(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		code             = "code"
		newdesc          = "desc"
		dto              = UpdateLessonDTO{
			Description: &newdesc,
		}
	)

	repo.EXPECT().updateLesson(ctx, code, dto).Return(nil)

	err := srv.UpdateLesson(ctx, code, dto)
	assert.NoError(t, err)
}

func Test_LessonsService_UpdateLesson_fail(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		code             = "code"
		dto              = UpdateLessonDTO{}
	)

	repo.EXPECT().updateLesson(ctx, code, dto).Return(ErrNoFieldsToUpdate)

	err := srv.UpdateLesson(ctx, code, dto)
	assert.Error(t, err)
	assert.Equal(t, err, ErrNoFieldsToUpdate)
}

func Test_LessonsService_DeleteLesson_success(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		code             = "lesson_code"
	)

	repo.EXPECT().deleteLesson(ctx, code).Return(nil)

	err := srv.DeleteLesson(ctx, code)
	assert.NoError(t, err)
}

func Test_LessonsService_DeleteLesson_fail(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		code             = "lesson_code"
	)

	repo.EXPECT().deleteLesson(ctx, code).Return(ErrNotFound)

	err := srv.DeleteLesson(ctx, code)
	assert.Error(t, err)
	assert.Equal(t, err, ErrNotFound)
}

func Test_LessonsService_AddExerciseToList_success(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		lessonCode       = "lesson_code"
		exerciseCode     = "exercise_code"
	)

	exercisesService.EXPECT().ExerciseExists(ctx, exerciseCode).Return(true, nil)
	repo.EXPECT().addExerciseToList(ctx, lessonCode, exerciseCode).Return(nil)

	err := srv.AddExerciseToList(ctx, lessonCode, exerciseCode)
	assert.NoError(t, err)
}

func Test_LessonsService_AddExerciseToList_fail(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		lessonCode       = "lesson_code"
		exerciseCode     = "exercise_code"
	)

	exercisesService.EXPECT().ExerciseExists(ctx, exerciseCode).Return(true, nil)
	repo.EXPECT().addExerciseToList(ctx, lessonCode, exerciseCode).Return(ErrExerciseAlreadyInSet)

	err := srv.AddExerciseToList(ctx, lessonCode, exerciseCode)
	assert.Error(t, err)
	assert.Equal(t, ErrExerciseAlreadyInSet, err)
}

func Test_LessonsService_AddExerciseToList_ExerciseExists(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		lessonCode       = "lesson_code"
		exerciseCode     = "exercise_code"
		errRepo          = errors.New("error")
	)

	t.Run("repo error", func(t *testing.T) {
		exercisesService.EXPECT().ExerciseExists(ctx, exerciseCode).Return(false, errRepo)
		err := srv.AddExerciseToList(ctx, lessonCode, exerciseCode)
		assert.Error(t, err)
		assert.Equal(t, errRepo, err)
	})

	t.Run("false", func(t *testing.T) {
		exercisesService.EXPECT().ExerciseExists(ctx, exerciseCode).Return(false, nil)
		err := srv.AddExerciseToList(ctx, lessonCode, exerciseCode)
		assert.Error(t, err)
		assert.Equal(t, exercises.ErrNotFound, err)
	})

}

func Test_LessonsService_DeleteExerciseFromList_success(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		lessonCode       = "lesson_code"
		exerciseCode     = "exercise_code"
	)

	repo.EXPECT().deleteExerciseFromList(ctx, lessonCode, exerciseCode).Return(nil)

	err := srv.DeleteExerciseFromList(ctx, lessonCode, exerciseCode)
	assert.NoError(t, err)
}

func Test_LessonsService_DeleteExerciseFromList_fail(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		lessonCode       = "lesson_code"
		exerciseCode     = "exercise_code"
	)

	repo.EXPECT().deleteExerciseFromList(ctx, lessonCode, exerciseCode).Return(ErrExerciseNotInList)

	err := srv.DeleteExerciseFromList(ctx, lessonCode, exerciseCode)
	assert.Error(t, err)
	assert.Equal(t, ErrExerciseNotInList, err)
}

func Test_LessonsService_GetLesson_success(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		code             = "lesson_code"
		lesson           = lesson{
			Code:        code,
			Title:       "Lesson Title",
			Description: "Lesson Description",
			Exercises:   []string{"ex1", "ex2"},
		}
		exercises = []exercises.Exercise{{Code: "ex1"}, {Code: "ex2"}}
	)

	repo.EXPECT().getLesson(ctx, code).Return(lesson, nil)
	exercisesService.EXPECT().GetExercisesByCodes(ctx, lesson.Exercises).Return(exercises, nil)

	result, err := srv.GetLesson(ctx, code)

	assert.NoError(t, err)
	assert.Equal(t, lesson.toDTO(exercises), result)
}

func Test_LessonsService_GetLesson_fail_repo(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		code             = "lesson_code"
	)

	repo.EXPECT().getLesson(ctx, code).Return(lesson{}, ErrNotFound)

	_, err := srv.GetLesson(ctx, code)

	assert.Error(t, err)
	assert.Equal(t, ErrNotFound.Error(), err.Error())
}

func Test_LessonsService_GetLesson_fail_exerciseService(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		code             = "lesson_code"
		lesson           = lesson{
			Code:        code,
			Title:       "Lesson Title",
			Description: "Lesson Description",
			Exercises:   []string{"ex1", "ex2"},
		}
	)

	repo.EXPECT().getLesson(ctx, code).Return(lesson, nil)
	exercisesService.EXPECT().GetExercisesByCodes(ctx, lesson.Exercises).Return(nil, errors.New("database exploded"))

	_, err := srv.GetLesson(ctx, code)

	assert.Error(t, err)
	assert.Equal(t, "database exploded", err.Error())
}

func Test_LessonsService_GetLessonsByCodes_success(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		codes            = []string{"lesson1", "lesson2"}
		lessons          = []lesson{
			{Code: "lesson1", Title: "Lesson 1", Exercises: []string{"ex1", "ex2"}},
			{Code: "lesson2", Title: "Lesson 2", Exercises: []string{"ex3"}},
		}
		exercisesMap = map[string][]exercises.Exercise{
			"lesson1": {{Code: "ex1"}, {Code: "ex2"}},
			"lesson2": {{Code: "ex3"}},
		}
		expectedDTOs = []LessonDTO{
			{Code: "lesson1", Title: "Lesson 1", Exercises: exercisesMap["lesson1"]},
			{Code: "lesson2", Title: "Lesson 2", Exercises: exercisesMap["lesson2"]},
		}
	)

	repo.EXPECT().getLessonsByCodes(ctx, codes).Return(lessons, nil)
	for _, lesson := range lessons {
		exercisesService.EXPECT().GetExercisesByCodes(ctx, lesson.Exercises).Return(exercisesMap[lesson.Code], nil)
	}

	result, err := srv.GetLessonsByCodes(ctx, codes)
	assert.NoError(t, err)
	assert.Equal(t, expectedDTOs, result)
}

func Test_LessonsService_GetLessonsByCodes_fail_repo(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		codes            = []string{"lesson1", "lesson2"}
	)

	repo.EXPECT().getLessonsByCodes(ctx, codes).Return(nil, errors.New("repo error"))

	result, err := srv.GetLessonsByCodes(ctx, codes)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func Test_LessonsService_GetLessonsByCodes_fail_exerciseService(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		codes            = []string{"lesson1"}
		lessons          = []lesson{
			{Code: "lesson1", Title: "Lesson 1", Exercises: []string{"ex1", "ex2"}},
		}
	)

	repo.EXPECT().getLessonsByCodes(ctx, codes).Return(lessons, nil)
	exercisesService.EXPECT().GetExercisesByCodes(ctx, lessons[0].Exercises).Return(nil, errors.New("exercise service error"))

	result, err := srv.GetLessonsByCodes(ctx, codes)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func Test_LessonsService_LessonExists(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		req              = "code"
		repoErr          = errors.New("error")
	)

	t.Run("success-true", func(t *testing.T) {
		repo.EXPECT().lessonExists(ctx, req).Return(true, nil)
		exists, err := srv.LessonExists(ctx, req)
		assert.NoError(t, err)
		assert.True(t, exists)
	})
	t.Run("success-false", func(t *testing.T) {
		repo.EXPECT().lessonExists(ctx, req).Return(false, nil)
		exists, err := srv.LessonExists(ctx, req)
		assert.NoError(t, err)
		assert.False(t, exists)
	})
	t.Run("repo error", func(t *testing.T) {
		repo.EXPECT().lessonExists(ctx, req).Return(true, repoErr)
		exists, err := srv.LessonExists(ctx, req)
		assert.Error(t, err)
		assert.False(t, exists)
	})
}

func Test_LessonsService_GetAllLessonsWithExercises(t *testing.T) {
	t.Parallel()
	var (
		ctx              = context.TODO()
		ctrl             = gomock.NewController(t)
		exercisesService = NewMockexerciseService(ctrl)
		repo             = NewMockrepository(ctrl)
		srv              = NewLessonsService(repo, exercisesService)
		repoErr          = errors.New("error")
		lessons          = []lesson{
			{
				Code:        "lesson_1",
				Title:       "Basics of Kazakh",
				Description: "Learn simple words and greetings",
				Exercises:   []string{"ex1", "ex2"},
				CreatedAt:   time.Now(),
				DeletedAt:   nil,
			},
			{
				Code:        "lesson_2",
				Title:       "Food Vocabulary",
				Description: "Names of common food items",
				Exercises:   []string{"ex3"},
				CreatedAt:   time.Now(),
				DeletedAt:   nil,
			},
		}

		exerciseMap = map[string][]exercises.Exercise{
			"lesson_1": {
				{
					Code:          "ex1",
					ExerciseType:  "multiple_choice",
					Question:      "How do you say 'hello' in Kazakh?",
					Options:       []string{"Сәлем", "Пока", "Как дела?"},
					CorrectAnswer: "Сәлем",
					Explanation:   "‘Сәлем’ means hello.",
					CreatedAt:     time.Now(),
				},
				{
					Code:          "ex2",
					ExerciseType:  "manual_typing",
					Question:      "Type the word for 'thanks' in Kazakh.",
					CorrectAnswer: "Рахмет",
					Explanation:   "‘Рахмет’ is the Kazakh word for thanks.",
					CreatedAt:     time.Now(),
				},
			},
			"lesson_2": {
				{
					Code:         "ex3",
					ExerciseType: "match_pairs",
					Question:     "Match the Kazakh words with English ones.",
					Pairs:        []exercises.Pair{{Term: "Сүт", Match: "Milk"}, {Term: "Нан", Match: "Bread"}},
					Explanation:  "Match based on meaning.",
					CreatedAt:    time.Now(),
				},
			},
		}
	)

	t.Run("success", func(t *testing.T) {
		var resultExpected []LessonDTO
		repo.EXPECT().getAllLessons(ctx).Return(lessons, nil)
		for _, lesson := range lessons {
			exercisesService.EXPECT().GetExercisesByCodes(ctx, lesson.Exercises).Return(exerciseMap[lesson.Code], nil)
			resultExpected = append(resultExpected, lesson.toDTO(exerciseMap[lesson.Code]))
		}
		result, err := srv.GetAllLessonsWithExercises(ctx)
		assert.NoError(t, err)
		assert.Equal(t, result, resultExpected)
	})

	t.Run("repo failed#1", func(t *testing.T) {
		repo.EXPECT().getAllLessons(ctx).Return(nil, repoErr)
		result, err := srv.GetAllLessonsWithExercises(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("repo failed#2", func(t *testing.T) {
		repo.EXPECT().getAllLessons(ctx).Return(lessons, nil)
		exercisesService.EXPECT().GetExercisesByCodes(ctx, lessons[0].Exercises).Return(nil, repoErr)
		result, err := srv.GetAllLessonsWithExercises(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
