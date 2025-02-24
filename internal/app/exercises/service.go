package exercises

import (
	"context"
	"uiren/pkg/logger"
)

type repository interface {
	getExercisesByCodes(ctx context.Context, codes []string) ([]exercise, error)
}

type ExerciseService struct {
	repo repository
}

func NewExerciseService(repo repository) *ExerciseService {
	return &ExerciseService{
		repo: repo,
	}
}

func (s ExerciseService) GetExercisesByCodes(ctx context.Context, codes []string) ([]ExerciseDTO, error) {
	logger.Info("ExerciseService.GetExercisesbyCodes new request")

	exercises, err := s.repo.getExercisesByCodes(ctx, codes)
	if err != nil {
		logger.Error("ExerciseService.GetExercisesbyCodes repo.getExercisesByCodes: ", err)
		return nil, err
	}

	var response []ExerciseDTO
	for _, exercise := range exercises {
		exerciseType, err := GetValidExerciseType(exercise.ExerciseType)
		if err != nil {
			logger.Error("ExerciseService.GetExercisesbyCodes GetValidExerciseType: ", err)
			return nil, err
		}
		response = append(response, exercise.toDTO(exerciseType))
	}

	return response, nil
}
