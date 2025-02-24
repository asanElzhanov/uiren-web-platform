package exercises

import (
	"errors"
	"time"
)

type ExerciseType string

const (
	multipleChoiceType ExerciseType = "multiple_choice"
	fillInTheBlankType ExerciseType = "fill_in_the_blank"
	translateType      ExerciseType = "translate"
	matchPairsType     ExerciseType = "match_pairs"
	orderWordsType     ExerciseType = "order_words"
)

func GetValidExerciseType(str string) (ExerciseType, error) {
	switch ExerciseType(str) {
	case multipleChoiceType, fillInTheBlankType, translateType, matchPairsType, orderWordsType:
		return ExerciseType(str), nil
	default:
		return "", errors.New("invalid exercise type")
	}
}

type Pair struct {
	Term  string `bson:"term" json:"term"`   // match_pairs
	Match string `bson:"match" json:"match"` // match_pairs
}

type exercise struct {
	Code          string     `bson:"code" json:"code"`
	ExerciseType  string     `bson:"type" json:"type"`
	Question      string     `bson:"question" json:"question"`
	Options       []string   `bson:"options,omitempty" json:"options,omitempty"`               // multiple_choice, fill_in_the_blank, order_words
	CorrectAnswer string     `bson:"correct_answer,omitempty" json:"correct_answer,omitempty"` // multiple_choice, fill_in_the_blank, translate
	CorrectOrder  []string   `bson:"correct_order,omitempty" json:"correct_order,omitempty"`   // order_words
	Pairs         []Pair     `bson:"pairs,omitempty" json:"pairs,omitempty"`                   // match_pairs
	Explanation   string     `bson:"explanation" json:"explanation"`
	Hints         []string   `bson:"hints" json:"hints"`
	CreatedAt     time.Time  `bson:"created_at" json:"created_at"`
	DeletedAt     *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type ExerciseDTO struct {
	Code          string
	ExerciseType  ExerciseType
	Question      string
	Options       []string
	CorrectAnswer string
	CorrectOrder  []string
	Pairs         []Pair
	Explanation   string
	Hints         []string
	CreatedAt     time.Time
	DeletedAt     time.Time
}

func (exercise exercise) toDTO(exerciseType ExerciseType) ExerciseDTO {
	var exerciseDeletedAt time.Time
	if exercise.DeletedAt != nil {
		exerciseDeletedAt = *exercise.DeletedAt
	}
	var exerciseOptions []string
	if exercise.Options != nil {
		exerciseOptions = exercise.Options
	}

	var exerciseCorrectAnswer string
	if exercise.CorrectAnswer != "" {
		exerciseCorrectAnswer = exercise.CorrectAnswer
	}

	var exerciseCorrectOrder []string
	if exercise.CorrectOrder != nil {
		exerciseCorrectOrder = exercise.CorrectOrder
	}

	var exercisePairs []Pair
	if exercise.Pairs != nil {
		exercisePairs = exercise.Pairs
	}

	return ExerciseDTO{
		Code:          exercise.Code,
		ExerciseType:  exerciseType,
		Question:      exercise.Question,
		Options:       exerciseOptions,
		CorrectAnswer: exerciseCorrectAnswer,
		CorrectOrder:  exerciseCorrectOrder,
		Pairs:         exercisePairs,
		Explanation:   exercise.Explanation,
		Hints:         exercise.Hints,
		CreatedAt:     exercise.CreatedAt,
		DeletedAt:     exerciseDeletedAt,
	}
}
