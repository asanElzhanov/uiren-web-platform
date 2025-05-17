package exercises

import (
	"context"
	"errors"
	"testing"
	"time"
	"uiren/pkg/logger"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	logger.InitLogger("info")
}

func Test_exerciseService_GetExercisesByCodes_success(t *testing.T) {
	t.Parallel()
	var (
		ctx  = context.TODO()
		ctrl = gomock.NewController(t)
		repo = NewMockrepository(ctrl)
		srv  = NewExerciseService(repo)
	)

	t.Run("get filled list", func(t *testing.T) {
		req := []string{"code1", "code2", "code3"}
		exercises := []Exercise{
			{Code: req[0]},
			{Code: req[1]},
			{Code: req[2]},
		}
		repo.EXPECT().getExercisesByCodes(ctx, req).Return(exercises, nil)
		list, err := srv.GetExercisesByCodes(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, list, exercises)
	})

	t.Run("get empty list", func(t *testing.T) {
		req := []string{}
		exercises := []Exercise{}
		repo.EXPECT().getExercisesByCodes(ctx, req).Return([]Exercise{}, nil)
		list, err := srv.GetExercisesByCodes(ctx, []string{})
		assert.NoError(t, err)
		assert.Equal(t, list, exercises)
	})
}

func Test_exerciseService_GetExercisesByCodes_repoError(t *testing.T) {
	t.Parallel()
	var (
		ctx       = context.TODO()
		ctrl      = gomock.NewController(t)
		repo      = NewMockrepository(ctrl)
		srv       = NewExerciseService(repo)
		repoError = errors.New("database exploded while query")
	)

	repo.EXPECT().getExercisesByCodes(ctx, []string{}).Return([]Exercise{}, repoError)

	_, err := srv.GetExercisesByCodes(ctx, []string{})
	assert.Error(t, err)
	assert.Equal(t, err, repoError)
}

func Test_exerciseService_GetExercise_success(t *testing.T) {
	t.Parallel()
	var (
		ctx      = context.TODO()
		ctrl     = gomock.NewController(t)
		repo     = NewMockrepository(ctrl)
		srv      = NewExerciseService(repo)
		exercise = Exercise{Code: "test_code"}
	)

	repo.EXPECT().getExercise(ctx, "test_code").Return(exercise, nil)

	result, err := srv.GetExercise(ctx, "test_code")
	assert.NoError(t, err)
	assert.Equal(t, exercise, result)
}

func Test_exerciseService_GetExercise_repoError(t *testing.T) {
	t.Parallel()
	var (
		ctx       = context.TODO()
		ctrl      = gomock.NewController(t)
		repo      = NewMockrepository(ctrl)
		srv       = NewExerciseService(repo)
		repoError = ErrNotFound
	)

	repo.EXPECT().getExercise(ctx, "test_code").Return(Exercise{}, repoError)

	result, err := srv.GetExercise(ctx, "test_code")
	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	assert.Equal(t, Exercise{}, result)
}

func Test_exerciseService_DeleteExercise_success(t *testing.T) {
	t.Parallel()
	var (
		ctx  = context.TODO()
		ctrl = gomock.NewController(t)
		repo = NewMockrepository(ctrl)
		srv  = NewExerciseService(repo)
	)

	repo.EXPECT().deleteExercise(ctx, "test_code").Return(nil)

	err := srv.DeleteExercise(ctx, "test_code")
	assert.NoError(t, err)
}

func Test_exerciseService_DeleteExercise_repoError(t *testing.T) {
	t.Parallel()
	var (
		ctx       = context.TODO()
		ctrl      = gomock.NewController(t)
		repo      = NewMockrepository(ctrl)
		srv       = NewExerciseService(repo)
		repoError = ErrNotFound
	)

	repo.EXPECT().deleteExercise(ctx, "test_code").Return(repoError)

	err := srv.DeleteExercise(ctx, "test_code")
	assert.Error(t, err)
	assert.Equal(t, repoError, err)
}

func Test_exerciseService_CreateExercise_success(t *testing.T) {
	t.Parallel()
	var (
		ctx  = context.TODO()
		ctrl = gomock.NewController(t)
		repo = NewMockrepository(ctrl)
		srv  = NewExerciseService(repo)

		correctAnswer = "1"
		correctOrder  = []string{"1", "2", "3"}
		pairs         = []Pair{
			{Term: "1",
				Match: "0",
			},
		}
	)
	/*
		dto := CreateExerciseDTO{
			Code:          "mc",
			ExerciseType:  multipleChoiceType,
			Question:      "some questions?",
			Hints:         []string{"some hint"},
			Explanation:   "explanation",
			Options:       []string{"1", "2", "3"},
			CorrectAnswer: &correctAnswer,
			CorrectOrder:  correctOrder,
			Pairs:         pairs,
		}
	*/

	t.Run("valid multiple choice exercise", func(t *testing.T) {
		dto := CreateExerciseDTO{
			Code:          "code",
			ExerciseType:  multipleChoiceType,
			Question:      "some questions?",
			Hints:         []string{"some hint"},
			Explanation:   "explanation",
			Options:       []string{"1", "2", "3"},
			CorrectAnswer: &correctAnswer,
		}
		repo.EXPECT().createExercise(ctx, createExerciseDTOMatcher{dto: dto}).Return(primitive.ObjectID{}, nil)

		id, err := srv.CreateExercise(ctx, dto)

		assert.NoError(t, err)
		assert.Equal(t, id, primitive.ObjectID{})
	})

	t.Run("valid manual typing exercise", func(t *testing.T) {
		dto := CreateExerciseDTO{
			Code:          "code",
			ExerciseType:  manualTypingType,
			Question:      "some questions?",
			Hints:         []string{"some hint"},
			Explanation:   "explanation",
			CorrectAnswer: &correctAnswer,
		}
		repo.EXPECT().createExercise(ctx, createExerciseDTOMatcher{dto: dto}).Return(primitive.ObjectID{}, nil)

		id, err := srv.CreateExercise(ctx, dto)

		assert.NoError(t, err)
		assert.Equal(t, id, primitive.ObjectID{})
	})

	t.Run("valid match pairs exercise", func(t *testing.T) {
		dto := CreateExerciseDTO{
			Code:         "code",
			ExerciseType: matchPairsType,
			Question:     "some questions?",
			Hints:        []string{"some hint"},
			Explanation:  "explanation",
			Pairs:        pairs,
		}
		repo.EXPECT().createExercise(ctx, createExerciseDTOMatcher{dto: dto}).Return(primitive.ObjectID{}, nil)

		id, err := srv.CreateExercise(ctx, dto)

		assert.NoError(t, err)
		assert.Equal(t, id, primitive.ObjectID{})
	})

	t.Run("valid order words exercise", func(t *testing.T) {
		dto := CreateExerciseDTO{
			Code:         "code",
			ExerciseType: orderWordsType,
			Question:     "some questions?",
			Hints:        []string{"some hint"},
			Explanation:  "explanation",
			Options:      []string{"1", "2", "3"},
			CorrectOrder: correctOrder,
		}
		repo.EXPECT().createExercise(ctx, createExerciseDTOMatcher{dto: dto}).Return(primitive.ObjectID{}, nil)

		id, err := srv.CreateExercise(ctx, dto)

		assert.NoError(t, err)
		assert.Equal(t, id, primitive.ObjectID{})
	})
}

func Test_exerciseService_CreateExercise_fail(t *testing.T) {
	t.Parallel()
	var (
		ctx  = context.TODO()
		ctrl = gomock.NewController(t)
		repo = NewMockrepository(ctrl)
		srv  = NewExerciseService(repo)

		correctAnswer = "1"
		correctOrder  = []string{"1", "2", "3"}
	)

	t.Run("invalid multiple choice exercise#1", func(t *testing.T) {
		dto := CreateExerciseDTO{
			Code:          "code",
			ExerciseType:  multipleChoiceType,
			Question:      "some questions?",
			Hints:         []string{"some hint"},
			Explanation:   "explanation",
			CorrectAnswer: &correctAnswer,
		}

		_, err := srv.CreateExercise(ctx, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrOptionsRequired)
	})

	t.Run("invalid multiple choice exercise#2", func(t *testing.T) {
		dto := CreateExerciseDTO{
			Code:         "code",
			ExerciseType: multipleChoiceType,
			Question:     "some questions?",
			Hints:        []string{"some hint"},
			Explanation:  "explanation",
			Options:      []string{"1", "2", "3"},
		}

		_, err := srv.CreateExercise(ctx, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrCorrectAnswerRequired)
	})

	t.Run("vinalid manual typing exercise#1", func(t *testing.T) {
		dto := CreateExerciseDTO{
			Code:         "code",
			ExerciseType: manualTypingType,
			Question:     "some questions?",
			Hints:        []string{"some hint"},
			Explanation:  "explanation",
		}

		_, err := srv.CreateExercise(ctx, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrCorrectAnswerRequired)
	})

	t.Run("invalid match pairs exercise#1", func(t *testing.T) {
		dto := CreateExerciseDTO{
			Code:         "code",
			ExerciseType: matchPairsType,
			Question:     "some questions?",
			Hints:        []string{"some hint"},
			Explanation:  "explanation",
		}

		_, err := srv.CreateExercise(ctx, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrPairsRequired)
	})

	t.Run("valid order words exercise#1", func(t *testing.T) {
		dto := CreateExerciseDTO{
			Code:         "code",
			ExerciseType: orderWordsType,
			Question:     "some questions?",
			Hints:        []string{"some hint"},
			Explanation:  "explanation",
			CorrectOrder: correctOrder,
		}

		_, err := srv.CreateExercise(ctx, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrOptionsRequired)
	})

	t.Run("valid order words exercise#2", func(t *testing.T) {
		dto := CreateExerciseDTO{
			Code:         "code",
			ExerciseType: orderWordsType,
			Question:     "some questions?",
			Hints:        []string{"some hint"},
			Explanation:  "explanation",
			Options:      []string{"1", "2", "3"},
		}

		_, err := srv.CreateExercise(ctx, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrCorrectOrderRequired)
	})

	t.Run("repo code already exists", func(t *testing.T) {
		dto := CreateExerciseDTO{
			Code:          "code",
			ExerciseType:  manualTypingType,
			Question:      "some questions?",
			Hints:         []string{"some hint"},
			Explanation:   "explanation",
			CorrectAnswer: &correctAnswer,
		}
		repo.EXPECT().createExercise(ctx, createExerciseDTOMatcher{dto: dto}).Return(primitive.NilObjectID, ErrCodeAlreadyExists)

		_, err := srv.CreateExercise(ctx, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrCodeAlreadyExists)
	})

	t.Run("repo some error", func(t *testing.T) {
		someError := errors.New("some")

		dto := CreateExerciseDTO{
			Code:          "code",
			ExerciseType:  manualTypingType,
			Question:      "some questions?",
			Hints:         []string{"some hint"},
			Explanation:   "explanation",
			CorrectAnswer: &correctAnswer,
		}
		repo.EXPECT().createExercise(ctx, createExerciseDTOMatcher{dto: dto}).Return(primitive.NilObjectID, someError)

		_, err := srv.CreateExercise(ctx, dto)

		assert.Error(t, err)
		assert.Equal(t, err, someError)
	})

	t.Run("error incorrect type", func(t *testing.T) {
		dto := CreateExerciseDTO{
			Code:          "code",
			ExerciseType:  "incorrectType",
			Question:      "some questions?",
			Hints:         []string{"some hint"},
			Explanation:   "explanation",
			CorrectAnswer: &correctAnswer,
		}

		_, err := srv.CreateExercise(ctx, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrIncorrectType)
	})
}

func Test_exerciseService_UpdateExercise_success(t *testing.T) {
	t.Parallel()
	var (
		ctx  = context.TODO()
		ctrl = gomock.NewController(t)
		repo = NewMockrepository(ctrl)
		srv  = NewExerciseService(repo)

		correctAnswer = "1"
		correctOrder  = []string{"1", "2", "3"}
		pairs         = []Pair{
			{Term: "1",
				Match: "0",
			},
		}
		question    = "quest"
		explanation = "exp"
		code        = "code"
	)

	t.Run("valid multiple choice exercise", func(t *testing.T) {
		dto := UpdateExerciseDTO{
			Question:      &question,
			Hints:         []string{"some hint"},
			Explanation:   &explanation,
			Options:       []string{"1", "2", "3"},
			CorrectAnswer: &correctAnswer,
		}
		repo.EXPECT().getExerciseType(ctx, code).Return(multipleChoiceType, nil)
		repo.EXPECT().updateExercise(ctx, code, dto).Return(nil)

		err := srv.UpdateExercise(ctx, code, dto)
		assert.NoError(t, err)
	})

	t.Run("valid manual typing exercise", func(t *testing.T) {
		dto := UpdateExerciseDTO{
			Question:      &question,
			Hints:         []string{"some hint"},
			Explanation:   &explanation,
			CorrectAnswer: &correctAnswer,
		}
		repo.EXPECT().getExerciseType(ctx, code).Return(manualTypingType, nil)
		repo.EXPECT().updateExercise(ctx, code, dto).Return(nil)

		err := srv.UpdateExercise(ctx, code, dto)
		assert.NoError(t, err)
	})

	t.Run("valid match pairs exercise", func(t *testing.T) {
		dto := UpdateExerciseDTO{
			Question:    &question,
			Hints:       []string{"some hint"},
			Explanation: &explanation,
			Pairs:       pairs,
		}
		repo.EXPECT().getExerciseType(ctx, code).Return(matchPairsType, nil)
		repo.EXPECT().updateExercise(ctx, code, dto).Return(nil)

		err := srv.UpdateExercise(ctx, code, dto)
		assert.NoError(t, err)
	})

	t.Run("valid order words exercise", func(t *testing.T) {
		dto := UpdateExerciseDTO{
			Question:     &question,
			Hints:        []string{"some hint"},
			Explanation:  &explanation,
			Options:      []string{"1", "2"},
			CorrectOrder: correctOrder,
		}
		repo.EXPECT().getExerciseType(ctx, code).Return(orderWordsType, nil)
		repo.EXPECT().updateExercise(ctx, code, dto).Return(nil)

		err := srv.UpdateExercise(ctx, code, dto)
		assert.NoError(t, err)
	})
}

func Test_exerciseService_UpdateExercise_fail(t *testing.T) {
	t.Parallel()
	var (
		ctx  = context.TODO()
		ctrl = gomock.NewController(t)
		repo = NewMockrepository(ctrl)
		srv  = NewExerciseService(repo)

		correctAnswer = "1"
		correctOrder  = []string{"1", "2", "3"}
		question      = "quest"
		explanation   = "exp"
		code          = "code"
	)

	t.Run("invalid multiple choice exercise#1", func(t *testing.T) {
		dto := UpdateExerciseDTO{
			Question:      &question,
			Hints:         []string{"some hint"},
			Explanation:   &explanation,
			CorrectAnswer: &correctAnswer,
		}

		repo.EXPECT().getExerciseType(ctx, code).Return(multipleChoiceType, nil)

		err := srv.UpdateExercise(ctx, code, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrOptionsRequired)
	})

	t.Run("invalid multiple choice exercise#2", func(t *testing.T) {
		dto := UpdateExerciseDTO{
			Question:    &question,
			Hints:       []string{"some hint"},
			Explanation: &explanation,
			Options:     []string{"1", "2", "3"},
		}
		repo.EXPECT().getExerciseType(ctx, code).Return(multipleChoiceType, nil)

		err := srv.UpdateExercise(ctx, code, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrCorrectAnswerRequired)
	})

	t.Run("vinalid manual typing exercise#1", func(t *testing.T) {
		dto := UpdateExerciseDTO{
			Question:    &question,
			Hints:       []string{"some hint"},
			Explanation: &explanation,
		}
		repo.EXPECT().getExerciseType(ctx, code).Return(manualTypingType, nil)
		err := srv.UpdateExercise(ctx, code, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrCorrectAnswerRequired)
	})

	t.Run("invalid match pairs exercise#1", func(t *testing.T) {
		dto := UpdateExerciseDTO{
			Question:    &question,
			Hints:       []string{"some hint"},
			Explanation: &explanation,
		}
		repo.EXPECT().getExerciseType(ctx, code).Return(matchPairsType, nil)

		err := srv.UpdateExercise(ctx, code, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrPairsRequired)
	})

	t.Run("valid order words exercise#1", func(t *testing.T) {
		dto := UpdateExerciseDTO{
			Question:     &question,
			Hints:        []string{"some hint"},
			Explanation:  &explanation,
			CorrectOrder: correctOrder,
		}
		repo.EXPECT().getExerciseType(ctx, code).Return(orderWordsType, nil)

		err := srv.UpdateExercise(ctx, code, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrOptionsRequired)
	})

	t.Run("valid order words exercise#2", func(t *testing.T) {
		dto := UpdateExerciseDTO{
			Question:    &question,
			Hints:       []string{"some hint"},
			Explanation: &explanation,
			Options:     []string{"1", "2", "3"},
		}
		repo.EXPECT().getExerciseType(ctx, code).Return(orderWordsType, nil)

		err := srv.UpdateExercise(ctx, code, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrCorrectOrderRequired)
	})

	t.Run("getExerciseType repo some error", func(t *testing.T) {

		dto := UpdateExerciseDTO{
			Question:    &question,
			Hints:       []string{"some hint"},
			Explanation: &explanation,
			Options:     []string{"1", "2", "3"},
		}
		repo.EXPECT().getExerciseType(ctx, code).Return("", ErrNotFound)

		err := srv.UpdateExercise(ctx, code, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNotFound)
	})
	t.Run("updateExercise repo some error", func(t *testing.T) {

		dto := UpdateExerciseDTO{
			Question:      &question,
			Hints:         []string{"some hint"},
			Explanation:   &explanation,
			Options:       []string{"1", "2", "3"},
			CorrectAnswer: &correctAnswer,
		}
		repo.EXPECT().getExerciseType(ctx, code).Return(multipleChoiceType, nil)
		repo.EXPECT().updateExercise(ctx, code, dto).Return(ErrNoFieldsToUpdate)

		err := srv.UpdateExercise(ctx, code, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrNoFieldsToUpdate)
	})

	t.Run("error incorrect type", func(t *testing.T) {
		dto := UpdateExerciseDTO{
			Question:      &question,
			Hints:         []string{"some hint"},
			Explanation:   &explanation,
			CorrectAnswer: &correctAnswer,
		}
		repo.EXPECT().getExerciseType(ctx, code).Return("incorrectType", nil)

		err := srv.UpdateExercise(ctx, code, dto)

		assert.Error(t, err)
		assert.Equal(t, err, ErrIncorrectType)
	})

}

func Test_exerciseService_GetAllExercises(t *testing.T) {
	var (
		ctx           = context.TODO()
		ctrl          = gomock.NewController(t)
		repo          = NewMockrepository(ctrl)
		srv           = NewExerciseService(repo)
		MockExercises = []Exercise{
			{
				Code:          "ex1",
				ExerciseType:  "multiple_choice",
				Question:      "What is the capital of Kazakhstan?",
				Hints:         []string{"It's a planned city", "Renamed in 2019"},
				Explanation:   "The capital is Astana.",
				Options:       []string{"Astana", "Almaty", "Shymkent"},
				CorrectAnswer: "Astana",
				CreatedAt:     time.Now(),
			},
			{
				Code:          "ex2",
				ExerciseType:  "manual_typing",
				Question:      "Type the Kazakh word for 'hello'.",
				Hints:         []string{"Starts with 'С'"},
				Explanation:   "'Сәлем' is the Kazakh word for hello.",
				CorrectAnswer: "Сәлем",
				CreatedAt:     time.Now(),
			},
			{
				Code:         "ex3",
				ExerciseType: "order_words",
				Question:     "Arrange the words: 'I love Kazakhstan'.",
				Hints:        []string{"Start with 'Мен'"},
				Explanation:  "The correct order is: Мен Қазақстанды жақсы көремін.",
				Options:      []string{"Мен", "жақсы", "Қазақстанды", "көремін"},
				CorrectOrder: []string{"Мен", "Қазақстанды", "жақсы", "көремін"},
				CreatedAt:    time.Now(),
			},
			{
				Code:         "ex4",
				ExerciseType: "match_pairs",
				Question:     "Match Kazakh food words with their English translations.",
				Hints:        []string{"All are food-related"},
				Explanation:  "Basic vocabulary for common food items.",
				Pairs: []Pair{
					{Term: "Сүт", Match: "Milk"},
					{Term: "Нан", Match: "Bread"},
					{Term: "Су", Match: "Water"},
				},
				CreatedAt: time.Now(),
			},
		}
		errRepo = errors.New("repo er")
	)

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().getAllExercises(ctx).Return(MockExercises, nil)
		result, err := srv.GetAllExercises(ctx)
		assert.NoError(t, err)
		assert.Equal(t, result, MockExercises)
	})
	t.Run("repo failed", func(t *testing.T) {
		repo.EXPECT().getAllExercises(ctx).Return(nil, errRepo)
		result, err := srv.GetAllExercises(ctx)
		assert.Error(t, err)
		assert.Equal(t, err, errRepo)
		assert.Nil(t, result)
	})
}

func Test_exerciseService_ExerciseExists(t *testing.T) {
	var (
		ctx     = context.TODO()
		ctrl    = gomock.NewController(t)
		repo    = NewMockrepository(ctrl)
		srv     = NewExerciseService(repo)
		code    = "code"
		errRepo = errors.New("err repo")
	)

	t.Run("success-true", func(t *testing.T) {
		repo.EXPECT().exerciseExists(ctx, code).Return(true, nil)
		exists, err := srv.ExerciseExists(ctx, code)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("success-false", func(t *testing.T) {
		repo.EXPECT().exerciseExists(ctx, code).Return(false, nil)
		exists, err := srv.ExerciseExists(ctx, code)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("repo failed", func(t *testing.T) {
		repo.EXPECT().exerciseExists(ctx, code).Return(true, errRepo)
		exists, err := srv.ExerciseExists(ctx, code)
		assert.Error(t, err)
		assert.Equal(t, err, errRepo)
		assert.False(t, exists)
	})
}
