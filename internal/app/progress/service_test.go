package progress

import (
	"context"
	"errors"
	"testing"
	"uiren/internal/app/achievements"
	"uiren/pkg/logger"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func init() {
	logger.InitLogger("info")
}

func Test_ProgressService_GetBadges_success(t *testing.T) {
	t.Parallel()
	var (
		ctx     = context.TODO()
		ctrl    = gomock.NewController(t)
		repo    = NewMockprogressReceiverRepo(ctrl)
		service = &ProgressService{receiverRepo: repo}
		userID  = "user123"
		resp    = []string{"badge1", "badge2"}
	)

	repo.EXPECT().getBadges(ctx, userID).Return(resp, nil)

	badges, err := service.GetBadges(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, resp, badges)
}

func Test_ProgressService_GetBadges_repoFailed(t *testing.T) {
	t.Parallel()
	var (
		ctx     = context.TODO()
		ctrl    = gomock.NewController(t)
		repo    = NewMockprogressReceiverRepo(ctrl)
		service = &ProgressService{receiverRepo: repo}
		userID  = "user123"
		errRepo = errors.New("db error")
	)

	repo.EXPECT().getBadges(ctx, userID).Return(nil, errRepo)

	_, err := service.GetBadges(ctx, userID)
	assert.Equal(t, errRepo, err)
}

func Test_ProgressService_GetXP_success(t *testing.T) {
	t.Parallel()
	var (
		ctx     = context.TODO()
		ctrl    = gomock.NewController(t)
		repo    = NewMockprogressReceiverRepo(ctrl)
		service = &ProgressService{receiverRepo: repo}
		userID  = "user123"
		xp      = 150
	)

	repo.EXPECT().getXP(ctx, userID).Return(xp, nil)

	result, err := service.GetXP(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, xp, result)
}

func Test_ProgressService_GetXP_repoFailed(t *testing.T) {
	t.Parallel()
	var (
		ctx     = context.TODO()
		ctrl    = gomock.NewController(t)
		repo    = NewMockprogressReceiverRepo(ctrl)
		service = &ProgressService{receiverRepo: repo}
		userID  = "user123"
		errRepo = errors.New("db error")
	)

	repo.EXPECT().getXP(ctx, userID).Return(0, errRepo)

	_, err := service.GetXP(ctx, userID)
	assert.Equal(t, errRepo, err)
}

func Test_ProgressService_GetAchievements_success(t *testing.T) {
	t.Parallel()
	var (
		ctx     = context.TODO()
		ctrl    = gomock.NewController(t)
		repo    = NewMockprogressReceiverRepo(ctrl)
		service = &ProgressService{receiverRepo: repo}
		userID  = "user123"
		resp    = []UserAchievement{
			{Progress: 100, AchievementName: "ach1"},
		}
	)

	repo.EXPECT().getAchievementsProgress(ctx, userID).Return(resp, nil)

	achievements, err := service.GetAchievements(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, resp, achievements)
}

func Test_ProgressService_GetAchievements_repoFailed(t *testing.T) {
	t.Parallel()
	var (
		ctx     = context.TODO()
		ctrl    = gomock.NewController(t)
		repo    = NewMockprogressReceiverRepo(ctrl)
		service = &ProgressService{receiverRepo: repo}
		userID  = "user123"
		errRepo = errors.New("db error")
	)

	repo.EXPECT().getAchievementsProgress(ctx, userID).Return(nil, errRepo)

	_, err := service.GetAchievements(ctx, userID)
	assert.Equal(t, errRepo, err)
}

func Test_ProgressService_RegisterNewBadge_success(t *testing.T) {
	t.Parallel()
	var (
		ctx     = context.TODO()
		ctrl    = gomock.NewController(t)
		repo    = NewMockprogressUpdaterRepo(ctrl)
		service = &ProgressService{updaterRepo: repo}
		req     = InsertBadgeRequest{
			Badge:       "badge1",
			Description: "description1",
		}
	)

	repo.EXPECT().insertBadge(ctx, req).Return(nil)

	err := service.RegisterNewBadge(ctx, req)
	assert.NoError(t, err)
}

func Test_ProgressService_RegisterNewBadge_repoFailed(t *testing.T) {
	t.Parallel()
	var (
		ctx     = context.TODO()
		ctrl    = gomock.NewController(t)
		repo    = NewMockprogressUpdaterRepo(ctrl)
		service = &ProgressService{updaterRepo: repo}
		req     = InsertBadgeRequest{
			Badge:       "badge1",
			Description: "description1",
		}
		errRepo = errors.New("db error")
	)

	repo.EXPECT().insertBadge(ctx, req).Return(errRepo)

	err := service.RegisterNewBadge(ctx, req)
	assert.Equal(t, errRepo, err)
}

func Test_ProgressService_UpdateProgress_success(t *testing.T) {
	t.Parallel()
	var (
		ctx        = context.TODO()
		ctrl       = gomock.NewController(t)
		updateRepo = NewMockprogressUpdaterRepo(ctrl)
		selectRepo = NewMockprogressReceiverRepo(ctrl)
		achSrv     = NewMockachievementService(ctrl)
		service    = &ProgressService{updaterRepo: updateRepo, receiverRepo: selectRepo, achService: achSrv}
		tx         = NewMocktransaction(ctrl)
		req        = UpdateUserProgressRequest{
			UserID:    "user123",
			NewBadges: []string{"badge1", "badge2"},
			XP:        100,
			AchievementsProgress: []AchievementProgress{
				{AchievementID: 1, EarnedProgress: 50},
				{AchievementID: 2, EarnedProgress: 100},
			},
		}

		achievementsList = map[int]achievements.AchievementDTO{
			1: {
				ID:   1,
				Name: "ach1",
				Levels: []achievements.AchievementLevel{
					{Level: 1, Threshold: 51},
					{Level: 2, Threshold: 100},
				},
			},
			2: {
				ID:   2,
				Name: "ach2",
				Levels: []achievements.AchievementLevel{
					{Level: 1, Threshold: 51},
					{Level: 2, Threshold: 100},
				},
			},
		}

		progresses = map[int]UserAchievement{
			1: {
				AchievementName:  achievementsList[1].Name,
				Level:            1,
				Progress:         0,
				Threshold:        51,
				LevelDescription: "fdsafads",
			},
			2: {
				AchievementName:  achievementsList[2].Name,
				Level:            1,
				Progress:         0,
				Threshold:        51,
				LevelDescription: "fdsaads",
			},
		}
	)

	updateRepo.EXPECT().beginTransaction(ctx).Return(tx, nil)
	updateRepo.EXPECT().addBadges(ctx, tx, AddBadgesRequest{
		UserID: req.UserID,
		Badges: req.NewBadges,
	}).Return(nil)
	updateRepo.EXPECT().addXP(ctx, tx, AddXPRequest{
		UserID: req.UserID,
		XP:     req.XP,
	}).Return(nil)

	selectRepo.EXPECT().getAchievementProgress(ctx, req.UserID, req.AchievementsProgress[0].AchievementID).Return(progresses[1], nil)
	achSrv.EXPECT().GetAchievement(ctx, 1).Return(achievementsList[1], nil)
	updateRepo.EXPECT().updateAchievementProgress(ctx, tx, UpdateAchievementProgressRequest{
		UserID: req.UserID,
		Progress: AchievementProgress{
			AchievementID:  req.AchievementsProgress[0].AchievementID,
			EarnedProgress: req.AchievementsProgress[0].EarnedProgress,
			NewLevel:       1,
		},
	})

	selectRepo.EXPECT().getAchievementProgress(ctx, req.UserID, req.AchievementsProgress[1].AchievementID).Return(UserAchievement{}, ErrAchievementProgressNotFound)
	achSrv.EXPECT().GetAchievement(ctx, 2).Return(achievementsList[2], nil)
	updateRepo.EXPECT().updateAchievementProgress(ctx, tx, UpdateAchievementProgressRequest{
		UserID: req.UserID,
		Progress: AchievementProgress{
			AchievementID:  req.AchievementsProgress[1].AchievementID,
			EarnedProgress: req.AchievementsProgress[1].EarnedProgress,
			NewLevel:       2,
		},
	})

	tx.EXPECT().Commit(ctx).Return(nil)

	err := service.UpdateUserProgress(ctx, req)
	assert.NoError(t, err)
}

func Test_ProgressService_UpdateProgress_success_no_achievements(t *testing.T) {
	t.Parallel()
	var (
		ctx        = context.TODO()
		ctrl       = gomock.NewController(t)
		updateRepo = NewMockprogressUpdaterRepo(ctrl)
		selectRepo = NewMockprogressReceiverRepo(ctrl)
		achSrv     = NewMockachievementService(ctrl)
		service    = &ProgressService{updaterRepo: updateRepo, receiverRepo: selectRepo, achService: achSrv}
		tx         = NewMocktransaction(ctrl)
		req        = UpdateUserProgressRequest{
			UserID:               "user123",
			NewBadges:            []string{"badge1", "badge2"},
			XP:                   100,
			AchievementsProgress: nil,
		}
	)

	updateRepo.EXPECT().beginTransaction(ctx).Return(tx, nil)
	updateRepo.EXPECT().addBadges(ctx, tx, AddBadgesRequest{
		UserID: req.UserID,
		Badges: req.NewBadges,
	}).Return(nil)
	updateRepo.EXPECT().addXP(ctx, tx, AddXPRequest{
		UserID: req.UserID,
		XP:     req.XP,
	}).Return(nil)

	tx.EXPECT().Commit(ctx).Return(nil)

	err := service.UpdateUserProgress(ctx, req)
	assert.NoError(t, err)
}

func Test_ProgressService_UpdateProgress_beginTransaction_beginTransaction_failed(t *testing.T) {
	t.Parallel()
	var (
		ctx        = context.TODO()
		ctrl       = gomock.NewController(t)
		updateRepo = NewMockprogressUpdaterRepo(ctrl)
		selectRepo = NewMockprogressReceiverRepo(ctrl)
		achSrv     = NewMockachievementService(ctrl)
		service    = &ProgressService{updaterRepo: updateRepo, receiverRepo: selectRepo, achService: achSrv}
		req        = UpdateUserProgressRequest{
			UserID:               "user123",
			NewBadges:            []string{"badge1", "badge2"},
			XP:                   100,
			AchievementsProgress: nil,
		}
		errToReturn = errors.New("transaction error")
	)

	updateRepo.EXPECT().beginTransaction(ctx).Return(nil, errToReturn)
	err := service.UpdateUserProgress(ctx, req)
	assert.Equal(t, err, errToReturn)
}

func Test_ProgressService_UpdateProgress_no_badges(t *testing.T) {
	t.Parallel()
	var (
		ctx        = context.TODO()
		ctrl       = gomock.NewController(t)
		updateRepo = NewMockprogressUpdaterRepo(ctrl)
		selectRepo = NewMockprogressReceiverRepo(ctrl)
		achSrv     = NewMockachievementService(ctrl)
		tx         = NewMocktransaction(ctrl)
		service    = &ProgressService{updaterRepo: updateRepo, receiverRepo: selectRepo, achService: achSrv}
		req        = UpdateUserProgressRequest{
			UserID:               "user123",
			NewBadges:            nil,
			XP:                   100,
			AchievementsProgress: nil,
		}
	)

	updateRepo.EXPECT().beginTransaction(ctx).Return(tx, nil)
	tx.EXPECT().Rollback(ctx)
	err := service.UpdateUserProgress(ctx, req)
	assert.Equal(t, err, ErrBadgeNotProvided)
}

func Test_ProgressService_UpdateProgress_addBadges_failed(t *testing.T) {
	t.Parallel()
	var (
		ctx        = context.TODO()
		ctrl       = gomock.NewController(t)
		updateRepo = NewMockprogressUpdaterRepo(ctrl)
		selectRepo = NewMockprogressReceiverRepo(ctrl)
		achSrv     = NewMockachievementService(ctrl)
		tx         = NewMocktransaction(ctrl)
		service    = &ProgressService{updaterRepo: updateRepo, receiverRepo: selectRepo, achService: achSrv}
		req        = UpdateUserProgressRequest{
			UserID:               "user123",
			NewBadges:            []string{"badge1", "badge2"},
			XP:                   100,
			AchievementsProgress: nil,
		}
		errToReturn = errors.New("db query failed")
	)

	updateRepo.EXPECT().beginTransaction(ctx).Return(tx, nil)
	updateRepo.EXPECT().addBadges(ctx, tx, AddBadgesRequest{
		UserID: req.UserID,
		Badges: req.NewBadges,
	}).Return(errToReturn)
	tx.EXPECT().Rollback(ctx)
	err := service.UpdateUserProgress(ctx, req)
	assert.Equal(t, err, errToReturn)
}

func Test_ProgressService_UpdateProgress_addXP_failed(t *testing.T) {
	t.Parallel()
	var (
		ctx        = context.TODO()
		ctrl       = gomock.NewController(t)
		updateRepo = NewMockprogressUpdaterRepo(ctrl)
		selectRepo = NewMockprogressReceiverRepo(ctrl)
		achSrv     = NewMockachievementService(ctrl)
		tx         = NewMocktransaction(ctrl)
		service    = &ProgressService{updaterRepo: updateRepo, receiverRepo: selectRepo, achService: achSrv}
		req        = UpdateUserProgressRequest{
			UserID:               "user123",
			NewBadges:            []string{"badge1", "badge2"},
			XP:                   100,
			AchievementsProgress: nil,
		}
		errToReturn = errors.New("db query failed")
	)

	updateRepo.EXPECT().beginTransaction(ctx).Return(tx, nil)
	updateRepo.EXPECT().addBadges(ctx, tx, AddBadgesRequest{
		UserID: req.UserID,
		Badges: req.NewBadges,
	}).Return(nil)
	updateRepo.EXPECT().addXP(ctx, tx, AddXPRequest{
		UserID: req.UserID,
		XP:     req.XP,
	}).Return(errToReturn)
	tx.EXPECT().Rollback(ctx)
	err := service.UpdateUserProgress(ctx, req)
	assert.Equal(t, err, errToReturn)
}

func Test_ProgressService_UpdateProgress_getAchievementProgress_failed(t *testing.T) {
	t.Parallel()
	var (
		ctx        = context.TODO()
		ctrl       = gomock.NewController(t)
		updateRepo = NewMockprogressUpdaterRepo(ctrl)
		selectRepo = NewMockprogressReceiverRepo(ctrl)
		achSrv     = NewMockachievementService(ctrl)
		tx         = NewMocktransaction(ctrl)
		service    = &ProgressService{updaterRepo: updateRepo, receiverRepo: selectRepo, achService: achSrv}
		req        = UpdateUserProgressRequest{
			UserID:    "user123",
			NewBadges: []string{"badge1", "badge2"},
			XP:        100,
			AchievementsProgress: []AchievementProgress{
				{
					AchievementID:  1,
					EarnedProgress: 50,
				},
			},
		}
		errToReturn = errors.New("db query failed")
	)

	updateRepo.EXPECT().beginTransaction(ctx).Return(tx, nil)
	updateRepo.EXPECT().addBadges(ctx, tx, AddBadgesRequest{
		UserID: req.UserID,
		Badges: req.NewBadges,
	}).Return(nil)
	updateRepo.EXPECT().addXP(ctx, tx, AddXPRequest{
		UserID: req.UserID,
		XP:     req.XP,
	}).Return(nil)
	selectRepo.EXPECT().getAchievementProgress(ctx, req.UserID, req.AchievementsProgress[0].AchievementID).Return(UserAchievement{}, errToReturn)
	tx.EXPECT().Rollback(ctx)
	err := service.UpdateUserProgress(ctx, req)
	assert.Equal(t, err, errToReturn)
}

func Test_ProgressService_UpdateProgress_GetAchievement_failed(t *testing.T) {
	t.Parallel()
	var (
		ctx        = context.TODO()
		ctrl       = gomock.NewController(t)
		updateRepo = NewMockprogressUpdaterRepo(ctrl)
		selectRepo = NewMockprogressReceiverRepo(ctrl)
		achSrv     = NewMockachievementService(ctrl)
		tx         = NewMocktransaction(ctrl)
		service    = &ProgressService{updaterRepo: updateRepo, receiverRepo: selectRepo, achService: achSrv}
		req        = UpdateUserProgressRequest{
			UserID:    "user123",
			NewBadges: []string{"badge1", "badge2"},
			XP:        100,
			AchievementsProgress: []AchievementProgress{
				{
					AchievementID:  1,
					EarnedProgress: 50,
				},
			},
		}
		errToReturn = errors.New("db query failed")
	)

	updateRepo.EXPECT().beginTransaction(ctx).Return(tx, nil)
	updateRepo.EXPECT().addBadges(ctx, tx, AddBadgesRequest{
		UserID: req.UserID,
		Badges: req.NewBadges,
	}).Return(nil)
	updateRepo.EXPECT().addXP(ctx, tx, AddXPRequest{
		UserID: req.UserID,
		XP:     req.XP,
	}).Return(nil)
	selectRepo.EXPECT().getAchievementProgress(ctx, req.UserID, req.AchievementsProgress[0].AchievementID).Return(UserAchievement{
		AchievementName: "name",
		Progress:        123}, nil)
	achSrv.EXPECT().GetAchievement(ctx, req.AchievementsProgress[0].AchievementID).Return(achievements.AchievementDTO{}, errToReturn)
	tx.EXPECT().Rollback(ctx)
	err := service.UpdateUserProgress(ctx, req)
	assert.Equal(t, err, errToReturn)
}

func Test_ProgressService_UpdateProgress_updateAchievementProgress_failed(t *testing.T) {
	t.Parallel()
	var (
		ctx        = context.TODO()
		ctrl       = gomock.NewController(t)
		updateRepo = NewMockprogressUpdaterRepo(ctrl)
		selectRepo = NewMockprogressReceiverRepo(ctrl)
		achSrv     = NewMockachievementService(ctrl)
		tx         = NewMocktransaction(ctrl)
		service    = &ProgressService{updaterRepo: updateRepo, receiverRepo: selectRepo, achService: achSrv}
		req        = UpdateUserProgressRequest{
			UserID:    "user123",
			NewBadges: []string{"badge1", "badge2"},
			XP:        100,
			AchievementsProgress: []AchievementProgress{
				{
					AchievementID:  1,
					EarnedProgress: 50,
				},
			},
		}
		errToReturn = errors.New("db query failed")
	)

	updateRepo.EXPECT().beginTransaction(ctx).Return(tx, nil)
	updateRepo.EXPECT().addBadges(ctx, tx, AddBadgesRequest{
		UserID: req.UserID,
		Badges: req.NewBadges,
	}).Return(nil)
	updateRepo.EXPECT().addXP(ctx, tx, AddXPRequest{
		UserID: req.UserID,
		XP:     req.XP,
	}).Return(nil)
	selectRepo.EXPECT().getAchievementProgress(ctx, req.UserID, req.AchievementsProgress[0].AchievementID).Return(UserAchievement{
		AchievementName: "name",
		Progress:        123}, nil)
	achSrv.EXPECT().GetAchievement(ctx, req.AchievementsProgress[0].AchievementID).Return(achievements.AchievementDTO{
		ID:     req.AchievementsProgress[0].AchievementID,
		Levels: []achievements.AchievementLevel{},
	}, nil)
	updateRepo.EXPECT().updateAchievementProgress(ctx, tx, UpdateAchievementProgressRequest{
		UserID:   req.UserID,
		Progress: req.AchievementsProgress[0],
	}).Return(errToReturn)
	tx.EXPECT().Rollback(ctx)
	err := service.UpdateUserProgress(ctx, req)
	assert.Equal(t, err, errToReturn)
}

func Test_ProgressService_updateAchievementProgress_level_calculation(t *testing.T) {
	t.Parallel()
	var (
		ctx        = context.TODO()
		ctrl       = gomock.NewController(t)
		updateRepo = NewMockprogressUpdaterRepo(ctrl)
		selectRepo = NewMockprogressReceiverRepo(ctrl)
		achSrv     = NewMockachievementService(ctrl)
		tx         = NewMocktransaction(ctrl)
		service    = &ProgressService{updaterRepo: updateRepo, receiverRepo: selectRepo, achService: achSrv}
	)

	t.Run("#1 test", func(t *testing.T) {
		req := UpdateAchievementProgressRequest{
			UserID: "uid",
			Progress: AchievementProgress{
				AchievementID:  1,
				EarnedProgress: 20,
			},
		}
		achievement := achievements.AchievementDTO{
			ID:   1,
			Name: "some)name",
			Levels: []achievements.AchievementLevel{
				{
					Level:     1,
					Threshold: 20,
				},
				{
					Level:     2,
					Threshold: 40,
				},
				{
					Level:     3,
					Threshold: 60,
				},
				{
					Level:     4,
					Threshold: 80,
				},
			},
		}
		selectRepo.EXPECT().getAchievementProgress(ctx, req.UserID, req.Progress.AchievementID).Return(UserAchievement{Level: 2, Threshold: 60, Progress: 50}, nil)
		achSrv.EXPECT().GetAchievement(ctx, 1).Return(achievement, nil)
		req.Progress.NewLevel = 4
		updateRepo.EXPECT().updateAchievementProgress(ctx, tx, req)
		err := service.updateAchievementProgress(ctx, tx, req)
		assert.NoError(t, err)
	})

	t.Run("#2 test", func(t *testing.T) {
		req := UpdateAchievementProgressRequest{
			UserID: "uid",
			Progress: AchievementProgress{
				AchievementID:  1,
				EarnedProgress: 10,
			},
		}
		achievement := achievements.AchievementDTO{
			ID:   1,
			Name: "some)name",
			Levels: []achievements.AchievementLevel{
				{
					Level:     1,
					Threshold: 20,
				},
				{
					Level:     2,
					Threshold: 40,
				},
				{
					Level:     3,
					Threshold: 60,
				},
				{
					Level:     4,
					Threshold: 80,
				},
			},
		}
		selectRepo.EXPECT().getAchievementProgress(ctx, req.UserID, req.Progress.AchievementID).Return(UserAchievement{}, ErrAchievementProgressNotFound)
		achSrv.EXPECT().GetAchievement(ctx, 1).Return(achievement, nil)
		req.Progress.NewLevel = 1
		updateRepo.EXPECT().updateAchievementProgress(ctx, tx, req)
		err := service.updateAchievementProgress(ctx, tx, req)
		assert.NoError(t, err)
	})

	t.Run("#3 test", func(t *testing.T) {
		req := UpdateAchievementProgressRequest{
			UserID: "uid",
			Progress: AchievementProgress{
				AchievementID:  1,
				EarnedProgress: 80,
			},
		}
		achievement := achievements.AchievementDTO{
			ID:   1,
			Name: "some)name",
			Levels: []achievements.AchievementLevel{
				{
					Level:     1,
					Threshold: 20,
				},
				{
					Level:     2,
					Threshold: 40,
				},
				{
					Level:     3,
					Threshold: 60,
				},
				{
					Level:     4,
					Threshold: 80,
				},
			},
		}
		selectRepo.EXPECT().getAchievementProgress(ctx, req.UserID, req.Progress.AchievementID).Return(UserAchievement{}, ErrAchievementProgressNotFound)
		achSrv.EXPECT().GetAchievement(ctx, 1).Return(achievement, nil)
		req.Progress.NewLevel = 4
		updateRepo.EXPECT().updateAchievementProgress(ctx, tx, req)
		err := service.updateAchievementProgress(ctx, tx, req)
		assert.NoError(t, err)
	})
}
