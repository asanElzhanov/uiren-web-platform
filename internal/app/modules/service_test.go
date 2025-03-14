package modules

import (
	"context"
	"errors"
	"testing"
	"uiren/internal/app/lessons"
	"uiren/pkg/logger"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	logger.InitLogger("info")
}

func Test_ModulesService_CreateModule_success(t *testing.T) {
	t.Parallel()
	var (
		ctx            = context.TODO()
		ctrl           = gomock.NewController(t)
		repo           = NewMockrepository(ctrl)
		lessonsService = NewMocklessonsService(ctrl)
		srv            = NewModulesService(repo, lessonsService)
		dto            = CreateModuleDTO{
			Code:        "module1",
			Title:       "Module 1",
			Description: "Desc",
		}
		expectedID = primitive.NewObjectID()
	)

	repo.EXPECT().createModule(ctx, CreateModuleWithLessonsMatcher{expected: dto}).Return(expectedID, nil)

	id, err := srv.CreateModule(ctx, dto)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
}

func Test_ModulesService_CreateModule_fail(t *testing.T) {
	t.Parallel()
	var (
		ctx            = context.TODO()
		ctrl           = gomock.NewController(t)
		repo           = NewMockrepository(ctrl)
		lessonsService = NewMocklessonsService(ctrl)
		srv            = NewModulesService(repo, lessonsService)
		dto            = CreateModuleDTO{Code: "module1"}
	)

	repo.EXPECT().createModule(ctx, CreateModuleWithLessonsMatcher{expected: dto}).Return(primitive.NilObjectID, ErrCodeAlreadyExists)

	id, err := srv.CreateModule(ctx, dto)
	assert.Error(t, err)
	assert.Equal(t, primitive.NilObjectID, id)
	assert.Equal(t, err, ErrCodeAlreadyExists)
}

func Test_ModulesService_DeleteModule_success(t *testing.T) {
	t.Parallel()
	var (
		ctx            = context.TODO()
		ctrl           = gomock.NewController(t)
		repo           = NewMockrepository(ctrl)
		lessonsService = NewMocklessonsService(ctrl)

		srv  = NewModulesService(repo, lessonsService)
		code = "module1"
	)

	repo.EXPECT().deleteModule(ctx, code).Return(nil)

	err := srv.DeleteModule(ctx, code)
	assert.NoError(t, err)
}

func Test_ModulesService_DeleteModule_fail(t *testing.T) {
	t.Parallel()
	var (
		ctx            = context.TODO()
		ctrl           = gomock.NewController(t)
		repo           = NewMockrepository(ctrl)
		lessonsService = NewMocklessonsService(ctrl)

		srv  = NewModulesService(repo, lessonsService)
		code = "module1"
	)

	repo.EXPECT().deleteModule(ctx, code).Return(ErrNotFound)

	err := srv.DeleteModule(ctx, code)
	assert.Error(t, err)
	assert.Equal(t, err, ErrNotFound)
}

func Test_ModulesService_UpdateModule_success(t *testing.T) {
	t.Parallel()
	var (
		ctx            = context.TODO()
		ctrl           = gomock.NewController(t)
		repo           = NewMockrepository(ctrl)
		lessonsService = NewMocklessonsService(ctrl)

		srv      = NewModulesService(repo, lessonsService)
		code     = "module1"
		newTitle = "new title"
		dto      = UpdateModuleDTO{Title: &newTitle}
	)

	repo.EXPECT().updateModule(ctx, code, dto).Return(nil)

	err := srv.UpdateModule(ctx, code, dto)
	assert.NoError(t, err)
}

func Test_ModulesService_UpdateModule_fail(t *testing.T) {
	t.Parallel()
	var (
		ctx            = context.TODO()
		ctrl           = gomock.NewController(t)
		repo           = NewMockrepository(ctrl)
		lessonsService = NewMocklessonsService(ctrl)

		srv   = NewModulesService(repo, lessonsService)
		code  = "module1"
		title = "new"
		dto   = UpdateModuleDTO{Title: &title}
	)

	repo.EXPECT().updateModule(ctx, code, dto).Return(ErrNoFieldsToUpdate)

	err := srv.UpdateModule(ctx, code, dto)
	assert.Error(t, err)
	assert.Equal(t, err, ErrNoFieldsToUpdate)
}

func Test_ModulesService_AddLessonToList_success(t *testing.T) {
	t.Parallel()
	var (
		ctx            = context.TODO()
		ctrl           = gomock.NewController(t)
		repo           = NewMockrepository(ctrl)
		lessonsService = NewMocklessonsService(ctrl)

		srv        = NewModulesService(repo, lessonsService)
		moduleCode = "module1"
		lessonCode = "lesson1"
	)

	repo.EXPECT().addLessonToList(ctx, moduleCode, lessonCode).Return(nil)

	err := srv.AddLessonToList(ctx, moduleCode, lessonCode)
	assert.NoError(t, err)
}

func Test_ModulesService_AddLessonToList_fail(t *testing.T) {
	t.Parallel()
	var (
		ctx            = context.TODO()
		ctrl           = gomock.NewController(t)
		repo           = NewMockrepository(ctrl)
		lessonsService = NewMocklessonsService(ctrl)

		srv        = NewModulesService(repo, lessonsService)
		moduleCode = "module1"
		lessonCode = "lesson1"
	)

	repo.EXPECT().addLessonToList(ctx, moduleCode, lessonCode).Return(ErrLessonAlreadyInSet)

	err := srv.AddLessonToList(ctx, moduleCode, lessonCode)
	assert.Error(t, err)
	assert.Equal(t, err, ErrLessonAlreadyInSet)
}

func Test_ModulesService_DeleteLessonFromList_success(t *testing.T) {
	t.Parallel()
	var (
		ctx            = context.TODO()
		ctrl           = gomock.NewController(t)
		repo           = NewMockrepository(ctrl)
		lessonsService = NewMocklessonsService(ctrl)

		srv        = NewModulesService(repo, lessonsService)
		moduleCode = "module1"
		lessonCode = "lesson1"
	)

	repo.EXPECT().deleteLessonFromList(ctx, moduleCode, lessonCode).Return(nil)

	err := srv.DeleteLessonFromList(ctx, moduleCode, lessonCode)
	assert.NoError(t, err)
}

func Test_ModulesService_DeleteLessonFromList_fail(t *testing.T) {
	t.Parallel()
	var (
		ctx            = context.TODO()
		ctrl           = gomock.NewController(t)
		repo           = NewMockrepository(ctrl)
		lessonsService = NewMocklessonsService(ctrl)

		srv        = NewModulesService(repo, lessonsService)
		moduleCode = "module1"
		lessonCode = "lesson1"
	)

	repo.EXPECT().deleteLessonFromList(ctx, moduleCode, lessonCode).Return(ErrLessonNotInList)

	err := srv.DeleteLessonFromList(ctx, moduleCode, lessonCode)
	assert.Error(t, err)
	assert.Equal(t, err, ErrLessonNotInList)
}

func Test_ModulesService_GetModule_success(t *testing.T) {
	t.Parallel()
	var (
		ctx       = context.TODO()
		ctrl      = gomock.NewController(t)
		repo      = NewMockrepository(ctrl)
		lessonSrv = NewMocklessonsService(ctrl)
		srv       = NewModulesService(repo, lessonSrv)
		code      = "module1"
		module    = Module{
			Code:    code,
			Title:   "Module 1",
			Lessons: []string{"lesson1", "lesson2"},
		}
		lessonsList = []lessons.LessonDTO{
			{Code: "lesson1", Title: "Lesson 1"},
			{Code: "lesson2", Title: "Lesson 2"},
		}
		expectedDTO = module.toDTO(lessonsList)
	)

	repo.EXPECT().getModule(ctx, code).Return(module, nil)
	lessonSrv.EXPECT().GetLessonsByCodes(ctx, module.Lessons).Return(lessonsList, nil)

	dto, err := srv.GetModule(ctx, code)
	assert.NoError(t, err)
	assert.Equal(t, expectedDTO, dto)
}

func Test_ModulesService_GetModule_fail_repo(t *testing.T) {
	t.Parallel()
	var (
		ctx       = context.TODO()
		ctrl      = gomock.NewController(t)
		repo      = NewMockrepository(ctrl)
		lessonSrv = NewMocklessonsService(ctrl)
		srv       = NewModulesService(repo, lessonSrv)
		code      = "module1"
	)

	repo.EXPECT().getModule(ctx, code).Return(Module{}, ErrNotFound)

	dto, err := srv.GetModule(ctx, code)
	assert.Error(t, err)
	assert.Equal(t, ModuleWithLessons{}, dto)
	assert.Equal(t, err, ErrNotFound)
}

func Test_ModulesService_GetModule_fail_lessons(t *testing.T) {
	t.Parallel()
	var (
		ctx       = context.TODO()
		ctrl      = gomock.NewController(t)
		repo      = NewMockrepository(ctrl)
		lessonSrv = NewMocklessonsService(ctrl)
		srv       = NewModulesService(repo, lessonSrv)
		code      = "module1"
		module    = Module{
			Code:    code,
			Title:   "Module 1",
			Lessons: []string{"lesson1", "lesson2"},
		}
	)

	repo.EXPECT().getModule(ctx, code).Return(module, nil)
	lessonSrv.EXPECT().GetLessonsByCodes(ctx, module.Lessons).Return(nil, errors.New("lessons error"))

	dto, err := srv.GetModule(ctx, code)
	assert.Error(t, err)
	assert.Equal(t, ModuleWithLessons{}, dto)
}
