package exercises

import (
	"fmt"
	"slices"
)

type createExerciseDTOMatcher struct {
	dto CreateExerciseDTO
}

func (m createExerciseDTOMatcher) Matches(x interface{}) bool {
	dto, ok := x.(CreateExerciseDTO)
	if !ok {
		return false
	}

	if dto.Code != m.dto.Code {
		return false
	}
	if dto.CorrectAnswer != nil && m.dto.CorrectAnswer != nil && *dto.CorrectAnswer != *m.dto.CorrectAnswer {
		return false
	}
	if dto.CorrectOrder != nil && m.dto.CorrectOrder != nil && !slices.Equal(dto.CorrectOrder, m.dto.CorrectOrder) {
		return false
	}
	if dto.ExerciseType != m.dto.ExerciseType {
		return false
	}
	if dto.Explanation != m.dto.Explanation {
		return false
	}
	if !slices.Equal(dto.Hints, m.dto.Hints) {
		return false
	}
	if dto.Options != nil && m.dto.Options != nil && !slices.Equal(dto.Options, m.dto.Options) {
		return false
	}
	if dto.Pairs != nil && m.dto.Pairs != nil && !slices.Equal(dto.Pairs, m.dto.Pairs) {
		return false
	}
	if dto.Question != dto.Question {
		return false
	}

	return true
}

func (m createExerciseDTOMatcher) String() string {
	return fmt.Sprintf("matches structs")
}
