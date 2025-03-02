package lessons

import (
	"time"
	"uiren/internal/app/exercises"
)

type lesson struct {
	Code        string     `bson:"code" json:"code"`
	Title       string     `bson:"title" json:"title"`
	Description string     `bson:"description" json:"description"`
	Exercises   []string   `bson:"exercises" json:"exercises"`
	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`
	DeletedAt   *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type LessonDTO struct {
	Code        string
	Title       string
	Description string
	Exercises   []exercises.ExerciseDTO
	CreatedAt   time.Time
	DeletedAt   time.Time
}

func (lesson lesson) toDTO(exercises []exercises.ExerciseDTO) LessonDTO {
	var lessonDeletedAt time.Time
	if lesson.DeletedAt != nil {
		lessonDeletedAt = *lesson.DeletedAt
	}
	return LessonDTO{
		Code:        lesson.Code,
		Title:       lesson.Title,
		Description: lesson.Description,
		Exercises:   exercises,
		CreatedAt:   lesson.CreatedAt,
		DeletedAt:   lessonDeletedAt,
	}
}

// for repo

type CreateLessonDTO struct {
	Code        string     `bson:"code" json:"code"`
	Title       string     `bson:"title" json:"title"`
	Description string     `bson:"description" json:"description"`
	Exercises   []string   `bson:"exercises"`
	CreatedAt   time.Time  `bson:"created_at"`
	DeletedAt   *time.Time `bson:"deleted_at"`
}

type UpdateLessonDTO struct {
	Title       *string `bson:"title" json:"title"`
	Description *string `bson:"description" json:"description"`
}
