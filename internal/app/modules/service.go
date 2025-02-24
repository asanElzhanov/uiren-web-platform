package modules

import (
	"context"
	"uiren/internal/app/lessons"
	"uiren/pkg/logger"
)

type repository interface {
	//createModule(ctx context.Context, dto CreateModuleDTO) (module, error)
	//updateModule(ctx context.Context, dto UpdateModuleDTO) (module, error)
	//deleteModule(ctx context.Context, code string) error
	getModule(ctx context.Context, code string) (module, error)
	//addLesson(ctx context.Context, lessonCode string) error
	//deleteLesson(ctx context.Context, lessonCode string) error
}

type lessonsService interface {
	GetLessonsByCodes(ctx context.Context, codes []string) ([]lessons.LessonDTO, error)
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

func (s ModulesService) GetModule(ctx context.Context, code string) (ModuleDTO, error) {
	logger.Info("ModulesService.GetModule new request")

	module, err := s.repo.getModule(ctx, code)
	if err != nil {
		logger.Error("ModulesService.GetModule s.modulesRepository.getModule: ", err)
		return ModuleDTO{}, err
	}

	lessonsList, err := s.lessonsService.GetLessonsByCodes(ctx, module.Lessons)
	if err != nil {
		logger.Error("ModulesService.GetModule s.lessonsRepository.GetLessons: ", err)
		return ModuleDTO{}, err
	}

	return module.toDTO(lessonsList), nil
}
