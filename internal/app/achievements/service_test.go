package achievements

import (
	"context"
	"errors"
	"testing"
	"time"
	"uiren/pkg/logger"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func init() {
	logger.InitLogger("debug")
}

func Test_achievementService_CreateAchievement_success(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)
		req             = "test_achievement_name"
		response        = AchievementDTO{
			ID:     1,
			Name:   req,
			Levels: []AchievementLevelDTO{},
		}
	)
	defer ctrl.Finish()

	achievementRepo.EXPECT().createAchievement(ctx, req).Return(achievement{
		id:   1,
		name: req,
	}, nil)

	result, err := service.CreateAchievement(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, response, result)
}

func Test_achievementService_CreateAchievement_repoFailed(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)
		req             = "test_achievement_name"
		response        = AchievementDTO{}
		errRepo         = ErrAchievementNameExists
	)
	defer ctrl.Finish()

	achievementRepo.EXPECT().createAchievement(ctx, req).Return(achievement{}, errRepo)

	result, err := service.CreateAchievement(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, response, result)
}

func Test_achievementService_UpdateAchievement_success(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)

		newName = "test_achievement_name"
		req     = UpdateAchievementDTO{
			ID:      1,
			NewName: newName,
		}
		response = newName
	)
	defer ctrl.Finish()

	achievementRepo.EXPECT().updateAchievement(ctx, req).Return(newName, nil)

	result, err := service.UpdateAchievement(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, response, result)
}

func Test_achievementService_UpdateAchievement_repoFailed(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)

		newName = "test_achievement_name"
		req     = UpdateAchievementDTO{
			ID:      1,
			NewName: newName,
		}
		response = ""
		errRepo  = ErrAchievementNameExists
	)
	defer ctrl.Finish()

	achievementRepo.EXPECT().updateAchievement(ctx, req).Return("", errRepo)

	result, err := service.UpdateAchievement(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, response, result)
}

func Test_achievementService_DeleteAchievement_success(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)

		req = 1
	)
	defer ctrl.Finish()

	achievementRepo.EXPECT().deleteAchievement(ctx, req).Return(nil)

	err := service.DeleteAchievement(ctx, req)

	assert.NoError(t, err)
}

func Test_achievementService_DeleteAchievement_repoFailed(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)

		req     = 1
		errRepo = ErrAchievementNotFound
	)
	defer ctrl.Finish()

	achievementRepo.EXPECT().deleteAchievement(ctx, req).Return(errRepo)

	err := service.DeleteAchievement(ctx, req)

	assert.Error(t, err)
}

func Test_achievementService_GetAchievement_Success(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)

		req  = 1
		info = achievement{
			id:        1,
			name:      "some_name",
			createdAt: time.Now(),
		}

		levels = []achievementLevel{
			{
				achID:       1,
				achName:     "some_name",
				level:       1,
				description: "some_description",
				threshold:   1,
				createdAt:   time.Now(),
			},
			{
				achID:       1,
				achName:     "some_name",
				level:       2,
				description: "some_description",
				threshold:   2,
				createdAt:   time.Now(),
			},
			{
				achID:       1,
				achName:     "some_name",
				level:       3,
				description: "some_description",
				threshold:   3,
				createdAt:   time.Now(),
			},
		}
		levelsDTO = []AchievementLevelDTO{levels[0].toDTO(), levels[1].toDTO(), levels[2].toDTO()}
		response  = info.toDTO(levelsDTO)
	)
	defer ctrl.Finish()

	achievementRepo.EXPECT().getAchievement(ctx, req).Return(info, nil)
	achievementRepo.EXPECT().getLevelsByAchievementID(ctx, req).Return(levelsDTO, nil)

	result, err := service.GetAchievement(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, response, result)
}

func Test_achievementService_GetAchievement_RepoFailed_getAchievement(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)

		req     = 1
		errRepo = ErrAchievementNotFound
	)
	defer ctrl.Finish()

	achievementRepo.EXPECT().getAchievement(ctx, req).Return(achievement{}, errRepo)

	result, err := service.GetAchievement(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, AchievementDTO{}, result)
}

func Test_achievementService_GetAchievement_RepoFailed_getLevelsByAchievementID(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)

		req  = 1
		info = achievement{
			id:        1,
			name:      "some_name",
			createdAt: time.Now(),
		}
		errRepo = errors.New("some_error")
	)
	defer ctrl.Finish()

	achievementRepo.EXPECT().getAchievement(ctx, req).Return(info, nil)
	achievementRepo.EXPECT().getLevelsByAchievementID(ctx, req).Return(nil, errRepo)

	result, err := service.GetAchievement(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, AchievementDTO{}, result)
}

func Test_achievementService_GetAchievement_getLevelsByAchievementID_emptyLevels(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)

		req  = 1
		info = achievement{
			id:        1,
			name:      "some_name",
			createdAt: time.Now(),
		}
		levelsDTO = []AchievementLevelDTO{}
		response  = info.toDTO(levelsDTO)
	)
	defer ctrl.Finish()

	achievementRepo.EXPECT().getAchievement(ctx, req).Return(info, nil)
	achievementRepo.EXPECT().getLevelsByAchievementID(ctx, req).Return(levelsDTO, nil)

	result, err := service.GetAchievement(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, response, result)
}

func Test_achievementService_AddAchievementLevel_Success(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)

		req = AddAchievementLevelDTO{
			AchID:       1,
			Description: "some_description",
			Threshold:   2,
		}
	)
	defer ctrl.Finish()

	achievementRepo.EXPECT().getLastLevelAndTreshold(ctx, req.AchID).Return(LevelData{
		Level:     1,
		Threshold: 1,
	}, nil)
	req.Level = 2
	achievementRepo.EXPECT().addLevel(ctx, req).Return(nil)

	err := service.AddAchievementLevel(ctx, req)

	assert.NoError(t, err)
}

func Test_achievementService_AddAchievementLevel_RepoFailed_getLastLevelAndTreshold(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)

		req = AddAchievementLevelDTO{
			AchID:       1,
			Description: "some_description",
			Threshold:   2,
		}
		errRepo = ErrAchievementNotFound
	)
	defer ctrl.Finish()

	achievementRepo.EXPECT().getLastLevelAndTreshold(ctx, req.AchID).Return(LevelData{}, errRepo)

	err := service.AddAchievementLevel(ctx, req)

	assert.Error(t, err)
}

func Test_achievementService_AddAchievementLevel_RepoFailed_addLevel_lowThreshold(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)

		req = AddAchievementLevelDTO{
			AchID:       1,
			Description: "some_description",
			Threshold:   2,
		}
	)
	defer ctrl.Finish()

	achievementRepo.EXPECT().getLastLevelAndTreshold(ctx, req.AchID).Return(LevelData{
		Level:     1,
		Threshold: 3,
	}, nil)

	err := service.AddAchievementLevel(ctx, req)

	assert.Error(t, err)
}

func Test_achievementService_AddAchievementLevel_RepoFailed_addLevel_repoFailed(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)

		req = AddAchievementLevelDTO{
			AchID:       1,
			Description: "some_description",
			Threshold:   2,
		}
		errRepo = ErrLevelExists
	)
	defer ctrl.Finish()

	achievementRepo.EXPECT().getLastLevelAndTreshold(ctx, req.AchID).Return(LevelData{
		Level:     1,
		Threshold: 1,
	}, nil)
	req.Level = 2
	achievementRepo.EXPECT().addLevel(ctx, req).Return(errRepo)

	err := service.AddAchievementLevel(ctx, req)

	assert.Error(t, err)
}

func Test_achievementService_DeleteAchievementLevel_Success(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)

		req = DeleteAchievementLevelDTO{
			AchID: 1,
			Level: 1,
		}
	)
	defer ctrl.Finish()

	tx := NewMocktransaction(ctrl)
	achievementRepo.EXPECT().beginTransaction(ctx).Return(tx, nil)
	achievementRepo.EXPECT().deleteLevel(ctx, tx, req).Return(nil)
	achievementRepo.EXPECT().decrementUpperLevels(ctx, tx, req).Return(nil)
	tx.EXPECT().Commit(ctx).Return(nil)
	err := service.DeleteAchievementLevel(ctx, req)

	assert.NoError(t, err)
}

func Test_achievementService_DeleteAchievementLevel_RepoFailed_beginTransaction(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)
		req             = DeleteAchievementLevelDTO{
			AchID: 1,
			Level: 1,
		}
		errRepo = errors.New("some_error")
	)
	defer ctrl.Finish()

	achievementRepo.EXPECT().beginTransaction(ctx).Return(nil, errRepo)

	err := service.DeleteAchievementLevel(ctx, req)

	assert.Error(t, err)
}

func Test_achievementService_DeleteAchievementLevel_RepoFailed_deleteLevel(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)
		req             = DeleteAchievementLevelDTO{
			AchID: 1,
			Level: 1,
		}
		errRepo = errors.New("some_error")
	)
	defer ctrl.Finish()

	tx := NewMocktransaction(ctrl)
	achievementRepo.EXPECT().beginTransaction(ctx).Return(tx, nil)
	achievementRepo.EXPECT().deleteLevel(ctx, tx, req).Return(errRepo)
	tx.EXPECT().Rollback(ctx)

	err := service.DeleteAchievementLevel(ctx, req)

	assert.Error(t, err)
}

func Test_achievementService_DeleteAchievementLevel_RepoFailed_decrementUpperLevels(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)
		req             = DeleteAchievementLevelDTO{
			AchID: 1,
			Level: 1,
		}
		errRepo = errors.New("some_error")
	)
	defer ctrl.Finish()

	tx := NewMocktransaction(ctrl)
	achievementRepo.EXPECT().beginTransaction(ctx).Return(tx, nil)
	achievementRepo.EXPECT().deleteLevel(ctx, tx, req).Return(nil)
	achievementRepo.EXPECT().decrementUpperLevels(ctx, tx, req).Return(errRepo)
	tx.EXPECT().Rollback(ctx)

	err := service.DeleteAchievementLevel(ctx, req)

	assert.Error(t, err)
}

func Test_achievementService_DeleteAchievementLevel_RepoFailed_commitTransaction(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)
		req             = DeleteAchievementLevelDTO{
			AchID: 1,
			Level: 1,
		}
		errRepo = errors.New("some_error")
	)
	defer ctrl.Finish()

	tx := NewMocktransaction(ctrl)
	achievementRepo.EXPECT().beginTransaction(ctx).Return(tx, nil)
	achievementRepo.EXPECT().deleteLevel(ctx, tx, req).Return(nil)
	achievementRepo.EXPECT().decrementUpperLevels(ctx, tx, req).Return(nil)
	tx.EXPECT().Commit(ctx).Return(errRepo)
	tx.EXPECT().Rollback(ctx)

	err := service.DeleteAchievementLevel(ctx, req)

	assert.Error(t, err)
}

func Test_achievementService_DeleteAchievementLevel_RepoFailed_commitTransaction_rollbackFailed(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)
		req             = DeleteAchievementLevelDTO{
			AchID: 1,
			Level: 1,
		}
		errCommit   = errors.New("commit")
		errRollback = errors.New("rollback")
	)
	defer ctrl.Finish()

	tx := NewMocktransaction(ctrl)
	achievementRepo.EXPECT().beginTransaction(ctx).Return(tx, nil)
	achievementRepo.EXPECT().deleteLevel(ctx, tx, req).Return(nil)
	achievementRepo.EXPECT().decrementUpperLevels(ctx, tx, req).Return(nil)
	tx.EXPECT().Commit(ctx).Return(errCommit)
	tx.EXPECT().Rollback(ctx).Return(errRollback)

	err := service.DeleteAchievementLevel(ctx, req)

	assert.ErrorIs(t, err, errCommit)
}

func Test_achievementService_DeleteAchievementLevel_RepoFailed_decrementUpperLevels_rollbackFailed(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)
		req             = DeleteAchievementLevelDTO{
			AchID: 1,
			Level: 1,
		}
		errDecrement = errors.New("decrement")
		errRollback  = errors.New("rollback")
	)
	defer ctrl.Finish()

	tx := NewMocktransaction(ctrl)
	achievementRepo.EXPECT().beginTransaction(ctx).Return(tx, nil)
	achievementRepo.EXPECT().deleteLevel(ctx, tx, req).Return(nil)
	achievementRepo.EXPECT().decrementUpperLevels(ctx, tx, req).Return(errDecrement)
	tx.EXPECT().Rollback(ctx).Return(errRollback)

	err := service.DeleteAchievementLevel(ctx, req)

	assert.ErrorIs(t, err, errDecrement)
}

func Test_achievementService_DeleteAchievementLevel_RepoFailed_deleteLevel_rollbackFailed(t *testing.T) {
	t.Parallel()
	var (
		ctx             = context.TODO()
		ctrl            = gomock.NewController(t)
		achievementRepo = NewMockachievementRepo(ctrl)
		service         = NewAchievementService(achievementRepo)
		req             = DeleteAchievementLevelDTO{
			AchID: 1,
			Level: 1,
		}
		errDelete   = errors.New("delete")
		errRollback = errors.New("rollback")
	)
	defer ctrl.Finish()

	tx := NewMocktransaction(ctrl)
	achievementRepo.EXPECT().beginTransaction(ctx).Return(tx, nil)
	achievementRepo.EXPECT().deleteLevel(ctx, tx, req).Return(errDelete)
	tx.EXPECT().Rollback(ctx).Return(errRollback)

	err := service.DeleteAchievementLevel(ctx, req)

	assert.ErrorIs(t, err, errDelete)
}
