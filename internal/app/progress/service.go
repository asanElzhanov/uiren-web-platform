package progress

import (
	"context"
	"uiren/internal/app/achievements"
	"uiren/pkg/logger"
)

//go:generate mockgen -source service.go -destination service_mock.go -package progress

type progressReceiverRepo interface {
	getAllBadges(ctx context.Context) ([]Badge, error)

	getUserBadges(ctx context.Context, id string) ([]string, error)
	getXP(ctx context.Context, id string) (int, error)
	getAchievementsProgress(ctx context.Context, id string) ([]UserAchievement, error)
	getAchievementProgress(ctx context.Context, userID string, achID int) (UserAchievement, error)
	getXPLeaderboard(ctx context.Context, limit int) (XPLeaderboard, error)
}

type progressUpdaterRepo interface {
	beginTransaction(ctx context.Context) (transaction, error)

	insertBadge(ctx context.Context, req Badge) error

	addBadges(ctx context.Context, tx transaction, req AddBadgesRequest) error
	addXP(ctx context.Context, tx transaction, req AddXPRequest) error
	updateAchievementProgress(ctx context.Context, tx transaction, req UpdateAchievementProgressRequest) error
}

type achievementService interface {
	GetAchievement(ctx context.Context, id int) (achievements.AchievementDTO, error)
}

type ProgressService struct {
	receiverRepo progressReceiverRepo
	updaterRepo  progressUpdaterRepo
	achService   achievementService
}

func NewProgressService(receiverRepo progressReceiverRepo, updaterRepo progressUpdaterRepo, achService achievementService) *ProgressService {
	return &ProgressService{
		receiverRepo: receiverRepo,
		updaterRepo:  updaterRepo,
		achService:   achService,
	}
}

func (s *ProgressService) GetBadges(ctx context.Context, user_id string) ([]string, error) {
	logger.Info("ProgressService.GetBadges new request")
	badges, err := s.receiverRepo.getUserBadges(ctx, user_id)
	if err != nil {
		logger.Error("ProgressService.GetBadges userProgressRepo.getBadges: ", err)
		return nil, err
	}
	return badges, nil
}

func (s *ProgressService) GetXP(ctx context.Context, user_id string) (int, error) {
	logger.Info("ProgressService.GetXP new request")
	xp, err := s.receiverRepo.getXP(ctx, user_id)
	if err != nil {
		logger.Error("ProgressService.GetXP userProgressRepo.getXP: ", err)
		return 0, err
	}
	return xp, nil
}

func (s *ProgressService) GetAchievements(ctx context.Context, user_id string) ([]UserAchievement, error) {
	logger.Info("ProgressService.GetAchievements new request")
	achievements, err := s.receiverRepo.getAchievementsProgress(ctx, user_id)
	if err != nil {
		logger.Error("ProgressService.GetAchievements userProgressRepo.getAchievements: ", err)
		return nil, err
	}
	return achievements, nil
}

func (s *ProgressService) UpdateUserProgress(ctx context.Context, req UpdateUserProgressRequest) error {
	logger.Info("ProgressService.UpdateUserProgress new request")

	tx, err := s.updaterRepo.beginTransaction(ctx)
	if err != nil {
		logger.Error("ProgressService.UpdateUserProgress beginTransaction: ", err)
		return err
	}

	commited := false
	defer func() {
		if !commited {
			_ = tx.Rollback(ctx)
		}
	}()

	if req.NewBadges == nil || len(req.NewBadges) == 0 {
		return ErrBadgeNotProvided
	}
	if err := s.updaterRepo.addBadges(ctx, tx, AddBadgesRequest{
		UserID: req.UserID,
		Badges: req.NewBadges,
	}); err != nil {
		logger.Error("ProgressService.UpdateUserProgress insertBadges: ", err)
		return err
	}

	if err := s.updaterRepo.addXP(ctx, tx, AddXPRequest{
		UserID: req.UserID,
		XP:     req.XP,
	}); err != nil {
		logger.Error("ProgressService.UpdateUserProgress addXP: ", err)
		return err
	}

	if req.AchievementsProgress != nil {
		for _, achievement := range req.AchievementsProgress {
			if err := s.updateAchievementProgress(ctx, tx, UpdateAchievementProgressRequest{
				UserID:   req.UserID,
				Progress: achievement,
			}); err != nil {
				logger.Error("ProgressService.UpdateUserProgress updateUserAchievementProgress: ", err)
				return err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error("ProgressService.UpdateUserProgress commit: ", err)
		return err
	}
	commited = true

	return nil
}

func (s *ProgressService) updateAchievementProgress(ctx context.Context, tx transaction, req UpdateAchievementProgressRequest) error {
	req.Progress.NewLevel = 0
	currentLevel, err := s.receiverRepo.getAchievementProgress(ctx, req.UserID, req.Progress.AchievementID)
	if err == ErrAchievementProgressNotFound {
		currentLevel = UserAchievement{
			Progress: 0,
		}
	} else if err != nil {
		logger.Error("ProgressService.updateAchievementProgress getAchievementProgress: ", err)
		return err
	}

	ach, err := s.achService.GetAchievement(ctx, req.Progress.AchievementID)
	if err != nil {
		logger.Error("ProgressService.updateAchievementProgress GetAchievement: ", err)
		return err
	}

	newLevel := 0
	newProgress := currentLevel.Progress + req.Progress.EarnedProgress
	for _, level := range ach.Levels {
		if level.Threshold <= newProgress {
			continue
		}
		newLevel = level.Level
		break
	}
	if newLevel == 0 && len(ach.Levels) > 0 {
		newLevel = ach.Levels[len(ach.Levels)-1].Level
	}
	req.Progress.NewLevel = newLevel

	if err := s.updaterRepo.updateAchievementProgress(ctx, tx, req); err != nil {
		logger.Error("ProgressService.updateAchievementProgress updateUserAchievementProgress: ", err)
		return err
	}

	return nil
}

func (s *ProgressService) RegisterNewBadge(ctx context.Context, req Badge) error {
	logger.Info("ProgressService.RegisterNewBadge new request")
	if err := s.updaterRepo.insertBadge(ctx, req); err != nil {
		logger.Error("ProgressService.RegisterNewBadge insertBadge: ", err)
		return err
	}
	return nil
}

func (s *ProgressService) GetAllBadges(ctx context.Context) ([]Badge, error) {
	logger.Info("ProgressService.GetAllBadges new request")

	return s.receiverRepo.getAllBadges(ctx)
}

func (s *ProgressService) GetXPLeaderboard(ctx context.Context, limit int) (XPLeaderboard, error) {
	logger.Info("ProgressService.GetXPLeaderboard new request")

	leaderboard, err := s.receiverRepo.getXPLeaderboard(ctx, limit)
	if err != nil {
		logger.Error("ProgressService.GetXPLeaderboard getXPLeaderboard")
		return XPLeaderboard{}, err
	}

	return leaderboard, nil
}
