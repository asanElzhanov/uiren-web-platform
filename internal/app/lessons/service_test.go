package lessons

import (
	"context"
	"errors"
	"testing"
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

	repo.EXPECT().addExerciseToList(ctx, lessonCode, exerciseCode).Return(ErrExerciseAlreadyInSet)

	err := srv.AddExerciseToList(ctx, lessonCode, exerciseCode)
	assert.Error(t, err)
	assert.Equal(t, ErrExerciseAlreadyInSet, err)
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
