package lessons

import (
	"context"
	"uiren/internal/app/exercises"
	"uiren/pkg/logger"
)

type repository interface {
	//createLesson(ctx context.Context, dto CreateLessonDTO) (lesson, error)
	//updateLesson(ctx context.Context, dto UpdateLessonDTO) (lesson, error)
	//deleteLesson(ctx context.Context, code string) error
	//getLesson(ctx context.Context, code string) (lesson, error)
	getLessonsByCodes(ctx context.Context, codes []string) ([]lesson, error)
	//addExercise(ctx context.Context, exerciseCode string) error
	//deleteExercise(ctx context.Context, exerciseCode string) error
}

type exerciseService interface {
	GetExercisesByCodes(ctx context.Context, codes []string) ([]exercises.ExerciseDTO, error)
}

type LessonsService struct {
	repo            repository
	exerciseService exerciseService
}

func NewLessonsService(repo repository, exerciseService exerciseService) *LessonsService {
	return &LessonsService{
		repo:            repo,
		exerciseService: exerciseService,
	}
}

func (s LessonsService) GetLessonsByCodes(ctx context.Context, codes []string) ([]LessonDTO, error) {
	logger.Info("LessonsService.GetLessonsByCodes new request")

	lessons, err := s.repo.getLessonsByCodes(ctx, codes)
	if err != nil {
		logger.Error("LessonsService.GetLessonsByCodes repo.getLessonsByCodes: ", err)
		return nil, err
	}

	var response []LessonDTO
	for _, lesson := range lessons {
		exerciseList, err := s.exerciseService.GetExercisesByCodes(ctx, lesson.Exercises)
		if err != nil {
			logger.Error("LessonsService.GetLessonsByCodes exerciseService: ", err)
			return nil, err
		}
		response = append(response, lesson.toDTO(exerciseList))
	}
	return response, nil
}
