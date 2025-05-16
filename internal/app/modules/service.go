package modules

import (
	"context"
	"time"
	"uiren/internal/app/lessons"
	"uiren/pkg/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source service.go -destination service_mock.go -package modules

type repository interface {
	createModule(ctx context.Context, dto CreateModuleDTO) (primitive.ObjectID, error)
	updateModule(ctx context.Context, code string, dto UpdateModuleDTO) error
	deleteModule(ctx context.Context, code string) error
	getModule(ctx context.Context, code string) (Module, error)
	addLessonToList(ctx context.Context, code, lessonCode string) error
	deleteLessonFromList(ctx context.Context, code, lessonCode string) error

	getAllModules(ctx context.Context) ([]Module, error)
}

type lessonsService interface {
	GetLessonsByCodes(ctx context.Context, codes []string) ([]lessons.LessonDTO, error)
	LessonExists(ctx context.Context, code string) (bool, error)
}

type ModulesService struct {
	repo           repository
	lessonsService lessonsService
}

func NewModulesService(repo repository, lessonsService lessonsService) *ModulesService {
	return &ModulesService{
		repo:           repo,
		lessonsService: lessonsService,
	}
}

func (s ModulesService) GetModule(ctx context.Context, code string) (ModuleWithLessons, error) {
	logger.Info("ModulesService.GetModule new request")

	module, err := s.repo.getModule(ctx, code)
	if err != nil {
		logger.Error("ModulesService.GetModule s.modulesRepository.getModule: ", err)
		return ModuleWithLessons{}, err
	}

	lessonsList, err := s.lessonsService.GetLessonsByCodes(ctx, module.Lessons)
	if err != nil {
		logger.Error("ModulesService.GetModule s.lessonsRepository.GetLessons: ", err)
		return ModuleWithLessons{}, err
	}

	return module.toDTO(lessonsList), nil
}

func (s ModulesService) CreateModule(ctx context.Context, dto CreateModuleDTO) (primitive.ObjectID, error) {
	logger.Info("ModulesService.CreateModule new request")
	dto.CreatedAt = time.Now()
	dto.Lessons = make([]string, 0)

	id, err := s.repo.createModule(ctx, dto)
	if err != nil {
		logger.Error("ModulesService.CreateModule repo.createModule: ", err)
		return primitive.NilObjectID, err
	}

	return id, nil
}

func (s ModulesService) DeleteModule(ctx context.Context, code string) error {
	logger.Info("ModulesService.DeleteModule new request")

	if err := s.repo.deleteModule(ctx, code); err != nil {
		logger.Error("ModulesService.DeleteModule repo.deleteModule: ", err)
		return err
	}

	return nil
}

func (s ModulesService) UpdateModule(ctx context.Context, code string, dto UpdateModuleDTO) error {
	logger.Info("ModulesService.UpdateModule new request")

	if err := s.repo.updateModule(ctx, code, dto); err != nil {
		logger.Error("ModulesService.UpdateModule repo.updateModule: ", err)
		return err
	}

	return nil
}

func (s ModulesService) AddLessonToList(ctx context.Context, code, lessonCode string) error {
	logger.Info("ModulesService.AddLessonToList new request")

	exists, err := s.lessonsService.LessonExists(ctx, lessonCode)
	if err != nil {
		logger.Error("ModulesService.AddLessonToList lessonsService.LessonExists: ", err)
		return err
	}

	if !exists {
		return lessons.ErrNotFound
	}

	if err := s.repo.addLessonToList(ctx, code, lessonCode); err != nil {
		logger.Error("ModulesService.AddLessonToList repo.addLesson: ", err)
		return err
	}

	return nil
}

func (s ModulesService) DeleteLessonFromList(ctx context.Context, code, lessonCode string) error {
	logger.Info("ModulesService.DeleteLessonFromList new request")

	if err := s.repo.deleteLessonFromList(ctx, code, lessonCode); err != nil {
		logger.Error("ModulesService.DeleteLessonFromList repo.deleteLessonFromList: ", err)
		return err
	}

	return nil
}

//todo: write tests

func (s ModulesService) GetModulesList(ctx context.Context) ([]Module, error) {
	logger.Info("ModulesService.GetModules new request")
	// todo: change to func with pagination
	return s.repo.getAllModules(ctx)
}

// todo: write tests
func (s ModulesService) GetAllModulesWithLessons(ctx context.Context) ([]ModuleWithLessons, error) {
	logger.Info("ModulesService.GetModulesWithLessons new request")

	modules, err := s.repo.getAllModules(ctx)
	if err != nil {
		logger.Error("ModulesService.GetModulesWithLessons repo.getAllModules: ", err)
		return nil, err
	}

	var result []ModuleWithLessons
	for _, module := range modules {
		lessons, err := s.lessonsService.GetLessonsByCodes(ctx, module.Lessons)
		if err != nil {
			logger.Error("ModulesService.GetModulesWithLessons lessonsService.GetLessonsByCodes: ", err)
			return nil, err
		}
		result = append(result, module.toDTO(lessons))
	}

	return result, nil
}
