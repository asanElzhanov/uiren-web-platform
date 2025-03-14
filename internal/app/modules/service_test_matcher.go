package modules

import (
	"fmt"
	"slices"
)

type CreateModuleWithLessonsMatcher struct {
	expected CreateModuleDTO
}

func (m CreateModuleWithLessonsMatcher) Matches(x interface{}) bool {
	dto, ok := x.(CreateModuleDTO)
	if !ok {
		return false
	}
	return dto.Code == m.expected.Code &&
		dto.Title == m.expected.Title &&
		dto.Description == m.expected.Description &&
		slices.Equal(dto.Lessons, m.expected.Lessons)
}

func (m CreateModuleWithLessonsMatcher) String() string {
	return fmt.Sprintf("matches CreateModuleWithLessons (ignoring CreatedAt)")
}
