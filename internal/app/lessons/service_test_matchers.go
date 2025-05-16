package lessons

import (
	"slices"

	"github.com/golang/mock/gomock"
)

type CreateLessonDTOMatcher struct {
	expected CreateLessonDTO
}

func (m CreateLessonDTOMatcher) Matches(x interface{}) bool {
	actual, ok := x.(CreateLessonDTO)
	if !ok {
		return false
	}

	return actual.Code == m.expected.Code &&
		actual.Title == m.expected.Title &&
		actual.Description == m.expected.Description &&
		slices.Equal(m.expected.Exercises, actual.Exercises)
}

func (m CreateLessonDTOMatcher) String() string {
	return "matches CreateLessonDTO (ignoring CreatedAt and DeletedAt)"
}

func MatchCreateLessonDTO(expected CreateLessonDTO) gomock.Matcher {
	return CreateLessonDTOMatcher{expected: expected}
}

type LessonDTOMatcher struct {
	expected LessonDTO
}

func (m LessonDTOMatcher) Matches(x interface{}) bool {
	actual, ok := x.(LessonDTO)
	if !ok {
		return false
	}

	return actual.Code == m.expected.Code &&
		actual.Title == m.expected.Title &&
		actual.Description == m.expected.Description &&
		len(m.expected.Exercises) == len(actual.Exercises)
}

func (m LessonDTOMatcher) String() string {
	return "matches LessonDTO (ignoring CreatedAt and DeletedAt)"
}

func MatchLessonDTO(expected LessonDTO) gomock.Matcher {
	return LessonDTOMatcher{expected: expected}
}
