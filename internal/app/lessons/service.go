package lessons

import (
	"context"
	"time"
	"uiren/internal/app/exercises"
	"uiren/pkg/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source service.go -destination service_mock.go -package lessons

type repository interface {
	createLesson(ctx context.Context, dto CreateLessonDTO) (primitive.ObjectID, error)
	updateLesson(ctx context.Context, code string, dto UpdateLessonDTO) error
	deleteLesson(ctx context.Context, code string) error
	getLesson(ctx context.Context, code string) (lesson, error)
	getLessonsByCodes(ctx context.Context, codes []string) ([]lesson, error)
	addExerciseToList(ctx context.Context, code, exerciseCode string) error
	deleteExerciseFromList(ctx context.Context, code, exerciseCode string) error

	getAllLessons(ctx context.Context) ([]lesson, error)
}

type exerciseService interface {
	GetExercisesByCodes(ctx context.Context, codes []string) ([]exercises.Exercise, error)
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
			logger.Error("LessonsService.GetLessonsByCodes exerciseService.GetExercisesByCodes: ", err)
			return nil, err
		}
		response = append(response, lesson.toDTO(exerciseList))
	}
	return response, nil
}

func (s LessonsService) GetLesson(ctx context.Context, code string) (LessonDTO, error) {
	logger.Info("LessonsService.GetLesson new request")

	lesson, err := s.repo.getLesson(ctx, code)
	if err != nil {
		logger.Error("LessonsService.GetLesson repo.getLesson: ", err)
		return LessonDTO{}, err
	}

	exerciseList, err := s.exerciseService.GetExercisesByCodes(ctx, lesson.Exercises)
	if err != nil {
		logger.Error("LessonsService.GetLessonsByCodes exerciseService.GetExercisesByCodes: ", err)
		return LessonDTO{}, err
	}

	return lesson.toDTO(exerciseList), nil
}

func (s LessonsService) CreateLesson(ctx context.Context, dto CreateLessonDTO) (primitive.ObjectID, error) {
	logger.Info("LessonsService.CreateLesson new request")
	dto.CreatedAt = time.Now()
	dto.Exercises = make([]string, 0)

	oid, err := s.repo.createLesson(ctx, dto)
	if err != nil {
		logger.Error("LessonsService.CreateLesson repo.createLesson: ", err)
		return primitive.NilObjectID, err
	}

	return oid, nil
}

func (s LessonsService) UpdateLesson(ctx context.Context, code string, dto UpdateLessonDTO) error {
	logger.Info("LessonsService.UpdateLesson new request")

	if err := s.repo.updateLesson(ctx, code, dto); err != nil {
		logger.Error("LessonsService.UpdateLesson repo.updateLesson: ", err)
		return err
	}

	return nil
}

func (s LessonsService) DeleteLesson(ctx context.Context, code string) error {
	logger.Info("LessonsService.DeleteLesson new request")

	if err := s.repo.deleteLesson(ctx, code); err != nil {
		logger.Error("LessonsService.DeleteService repo.deleteLesson: ", err)
		return err
	}

	return nil
}

func (s LessonsService) AddExerciseToList(ctx context.Context, code, exerciseCode string) error {
	logger.Info("LessonsService.AddExerciseToList new request")

	//todo check if exists
	if err := s.repo.addExerciseToList(ctx, code, exerciseCode); err != nil {
		logger.Error("LessonsService.AddExerciseToList repo.addExerciseToList: ", err)
		return err
	}

	return nil
}

func (s LessonsService) DeleteExerciseFromList(ctx context.Context, code, exerciseCode string) error {
	logger.Info("LessonsService.DeleteExerciseFromList new request")

	if err := s.repo.deleteExerciseFromList(ctx, code, exerciseCode); err != nil {
		logger.Error("LessonsService.DeleteExerciseFromList repo.addExerciseToList: ", err)
		return err
	}

	return nil
}

// todo: write tests
func (s LessonsService) GetAllLessonsWithExercises(ctx context.Context) ([]LessonDTO, error) {
	logger.Info("LessonsService.GetAllLessonsWithExercises new request")

	lessons, err := s.repo.getAllLessons(ctx)
	if err != nil {
		logger.Error("LessonsService.GetAllLessonsWithExercises repo.getAllLessons: ", err)
		return nil, err
	}

	var result []LessonDTO
	for _, lesson := range lessons {
		exercises, err := s.exerciseService.GetExercisesByCodes(ctx, lesson.Exercises)
		if err != nil {
			logger.Error("LessonsService.GetAllLessonsWithExercises exerciseService.GetExercisesByCodes: ", err)
			return nil, err
		}
		result = append(result, lesson.toDTO(exercises))
	}
	return result, nil
}
