package exercises

import (
	"context"
	"time"
	"uiren/pkg/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source service.go -destination service_mock.go -package exercises

type repository interface {
	getExercisesByCodes(ctx context.Context, codes []string) ([]Exercise, error)
	getExercise(ctx context.Context, code string) (Exercise, error)
	createExercise(ctx context.Context, dto CreateExerciseDTO) (primitive.ObjectID, error)
	updateExercise(ctx context.Context, code string, dto UpdateExerciseDTO) error
	getExerciseType(ctx context.Context, code string) (string, error)
	deleteExercise(ctx context.Context, code string) error

	getAllExercises(ctx context.Context) ([]Exercise, error)
}

type ExerciseService struct {
	repo repository
}

func NewExerciseService(repo repository) *ExerciseService {
	return &ExerciseService{
		repo: repo,
	}
}

func (s ExerciseService) GetExercisesByCodes(ctx context.Context, codes []string) ([]Exercise, error) {
	logger.Info("ExerciseService.GetExercisesbyCodes new request")

	exercises, err := s.repo.getExercisesByCodes(ctx, codes)
	if err != nil {
		logger.Error("ExerciseService.GetExercisesbyCodes repo.getExercisesByCodes: ", err)
		return nil, err
	}

	return exercises, nil
}

func (s ExerciseService) GetExercise(ctx context.Context, code string) (Exercise, error) {
	logger.Info("ExerciseService.GetExercise new request")

	exercise, err := s.repo.getExercise(ctx, code)
	if err != nil {
		logger.Error("ExerciseService.GetExercise repo.getExercise: ", err)
		return Exercise{}, err
	}

	return exercise, nil
}

func (s ExerciseService) DeleteExercise(ctx context.Context, code string) error {
	logger.Info("ExerciseService.DeleteExercise new request")

	if err := s.repo.deleteExercise(ctx, code); err != nil {
		logger.Error("ExerciseService.DeleteExercise repo.deleteExercise: ", err)
		return err
	}

	return nil
}

func (s ExerciseService) CreateExercise(ctx context.Context, dto CreateExerciseDTO) (primitive.ObjectID, error) {
	logger.Info("ExerciseService.CreateExercise new request")

	var newDTO CreateExerciseDTO
	newDTO.Code = dto.Code
	newDTO.ExerciseType = dto.ExerciseType
	newDTO.Question = dto.Question
	newDTO.Hints = dto.Hints
	newDTO.Explanation = dto.Explanation
	newDTO.CreatedAt = time.Now()
	newDTO.DeletedAt = nil

	var err error
	switch dto.ExerciseType {
	case multipleChoiceType:
		err = normalizeMultipleChoiceExerciseDTO(&dto, &newDTO)
	case manualTypingType:
		err = normalizeManualTypingExerciseDTO(&dto, &newDTO)
	case matchPairsType:
		err = normalizeMatchPairsExerciseDTO(&dto, &newDTO)
	case orderWordsType:
		err = normalizeOrderWordsExerciseDTO(&dto, &newDTO)
	default:
		err = ErrIncorrectType
	}
	if err != nil {
		logger.Error("ExerciseService.CreateExercise normalizeExerciseDTO: ", err)
		return primitive.NilObjectID, err
	}

	oid, err := s.repo.createExercise(ctx, newDTO)
	if err != nil {
		logger.Error("ExerciseService.CreateExercise repo.createExercise: ", err)
		return primitive.NilObjectID, err
	}

	return oid, nil
}

func (s ExerciseService) UpdateExercise(ctx context.Context, code string, dto UpdateExerciseDTO) error {
	logger.Info("ExerciseService.UpdateExercise new request")

	var newDTO UpdateExerciseDTO
	newDTO.Question = dto.Question
	newDTO.Hints = dto.Hints
	newDTO.Explanation = dto.Explanation
	exerciseType, err := s.repo.getExerciseType(ctx, code)
	if err != nil {
		logger.Error("ExerciseService.UpdateExercise repo.getExerciseType: ", err)
		return err
	}

	switch exerciseType {
	case multipleChoiceType:
		err = normalizeMultipleChoiceExerciseDTO(&dto, &newDTO)
	case manualTypingType:
		err = normalizeManualTypingExerciseDTO(&dto, &newDTO)
	case matchPairsType:
		err = normalizeMatchPairsExerciseDTO(&dto, &newDTO)
	case orderWordsType:
		err = normalizeOrderWordsExerciseDTO(&dto, &newDTO)
	default:
		err = ErrIncorrectType
	}
	if err != nil {
		logger.Error("ExerciseService.UpdateExercise normalizeExerciseDTO: ", err)
		return err
	}

	if err := s.repo.updateExercise(ctx, code, newDTO); err != nil {
		logger.Error("ExerciseService.UpdateExercise updateExercise: ", err)
		return err
	}

	return nil
}

// todo write tests
func (s ExerciseService) GetAllExercises(ctx context.Context) ([]Exercise, error) {
	logger.Info("ExerciseService.GetAllExercises new requests")

	exercises, err := s.repo.getAllExercises(ctx)
	if err != nil {
		logger.Error("ExerciseService.GetAllExercises repo.getAllExercises: ", err)
		return nil, err
	}

	return exercises, nil
}
