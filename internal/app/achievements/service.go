package achievements

import (
	"context"
	"uiren/pkg/logger"

	"github.com/jackc/pgx/v5/pgconn"
)

//go:generate mockgen -source service.go -destination service_mock.go -package achievements

type achievementRepo interface {
	getAllAchievements(ctx context.Context) ([]achievement, error)
	getAchievement(ctx context.Context, id int) (achievement, error)
	createAchievement(ctx context.Context, name string) (achievement, error)
	updateAchievement(ctx context.Context, dto UpdateAchievementDTO) (string, error)
	deleteAchievement(ctx context.Context, id int) error

	beginTransaction(ctx context.Context) (transaction, error)

	getLevelsByAchievementID(ctx context.Context, achID int) ([]AchievementLevel, error)
	getLevel(ctx context.Context, achID, level int) (AchievementLevel, error)
	getLastLevelAndTreshold(ctx context.Context, achID int) (LevelData, error)
	addLevel(ctx context.Context, dto AddAchievementLevelDTO) error
	deleteLevel(ctx context.Context, tx transaction, dto DeleteAchievementLevelDTO) error
	decrementUpperLevels(ctx context.Context, tx transaction, dto DeleteAchievementLevelDTO) error
}

type transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
}

type AchievementService struct {
	achievementRepo achievementRepo
}

func NewAchievementService(achievementRepo achievementRepo) *AchievementService {
	return &AchievementService{
		achievementRepo: achievementRepo,
	}
}

func (s AchievementService) CreateAchievement(ctx context.Context, name string) (AchievementDTO, error) {
	logger.Info("AchievementService.CreateAchievement new request")

	achievement, err := s.achievementRepo.createAchievement(ctx, name)
	if err != nil {
		logger.Error("AchievementService.CreateAchievement achievementRepo.createAchievement: ", err)
		return AchievementDTO{}, err
	}

	emptyLevelsList := []AchievementLevel{}
	return achievement.toDTO(emptyLevelsList), nil
}

func (s AchievementService) UpdateAchievement(ctx context.Context, dto UpdateAchievementDTO) (string, error) {
	logger.Info("AchievementService.UpdateAchievement new request")

	newName, err := s.achievementRepo.updateAchievement(ctx, dto)
	if err != nil {
		logger.Error("AchievementService.UpdateAchievement achievementRepo.updateAchievement: ", err)
		return "", err
	}

	return newName, nil
}

func (s AchievementService) DeleteAchievement(ctx context.Context, id int) error {
	logger.Info("AchievementService.DeleteAchievement new request")

	if err := s.achievementRepo.deleteAchievement(ctx, id); err != nil {
		logger.Error("AchievementService.DeleteAchievement achievementRepo.deleteAchievement: ", err)
		return err
	}
	//todo: delete achievement levels
	return nil
}

func (s AchievementService) GetAchievement(ctx context.Context, id int) (AchievementDTO, error) {
	logger.Info("AchievementService.GetAchievement new request")

	achievement, err := s.achievementRepo.getAchievement(ctx, id)
	if err != nil {
		logger.Error("AchievementService.GetAchievement achievementRepo.getAchievement: ", err)
		return AchievementDTO{}, err
	}

	levels, err := s.achievementRepo.getLevelsByAchievementID(ctx, id)
	if err != nil {
		logger.Error("AchievementService.GetAchievement achievementRepo.getLevelsByAchievementID: ", err)
		return AchievementDTO{}, err
	}

	return achievement.toDTO(levels), nil
}

func (s AchievementService) AddAchievementLevel(ctx context.Context, dto AddAchievementLevelDTO) error {
	logger.Info("AchievementService.AddAchievementLevel new request")
	if dto.Threshold <= 0 {
		return ErrInvalidThreshold
	}

	levelData, err := s.achievementRepo.getLastLevelAndTreshold(ctx, dto.AchID)
	if err != nil {
		logger.Error("AchievementService.AddAchievementLevel achievementRepo.getLastLevel: ", err)
		return err
	}

	if levelData.Threshold >= dto.Threshold {
		return ErrLowThreshold
	}
	dto.Level = levelData.Level + 1

	if err := s.achievementRepo.addLevel(ctx, dto); err != nil {
		logger.Error("AchievementService.AddAchievementLevel achievementRepo.addLevel: ", err)
		return err
	}
	return nil
}

func (s AchievementService) DeleteAchievementLevel(ctx context.Context, dto DeleteAchievementLevelDTO) error {
	logger.Info("AchievementService.DeleteAchievementLevel new request")
	tx, err := s.achievementRepo.beginTransaction(ctx)
	if err != nil {
		logger.Error("AchievementService.DeleteAchievementLevel achievementRepo.beginTransaction: ", err)
		return err
	}

	commited := false
	defer func() {
		if !commited {
			_ = tx.Rollback(ctx)
		}
	}()

	if err := s.achievementRepo.deleteLevel(ctx, tx, dto); err != nil {
		logger.Error("AchievementService.DeleteAchievementLevel achievementRepo.deleteLevel: ", err)
		return err
	}

	if err := s.achievementRepo.decrementUpperLevels(ctx, tx, dto); err != nil {
		logger.Error("AchievementService.DeleteAchievementLevel achievementRepo.decrementUpperLevels: ", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error("AchievementService.DeleteAchievementLevel tx.Commit: ", err)
		return err
	}

	commited = true
	return nil
}

func (s AchievementService) GetLevel(ctx context.Context, achID, level int) (AchievementLevel, error) {
	logger.Info("AchievementService.GetLevel new request")

	achievementLevel, err := s.achievementRepo.getLevel(ctx, achID, level)
	if err != nil {
		logger.Error("AchievementService.GetLevel achievementRepo.getLevel: ", err)
		return AchievementLevel{}, err
	}

	return achievementLevel, nil
}

// todo write tests
func (s AchievementService) GetAllAchievements(ctx context.Context) ([]AchievementDTO, error) {
	logger.Info("AchievementService.GetAllAchievements new request")

	achievements, err := s.achievementRepo.getAllAchievements(ctx)
	if err != nil {
		logger.Error("AchievementService.GetAllAchievements achievementRepo.getAllAchievements: ", err)
		return nil, err
	}

	var result []AchievementDTO
	for _, achievement := range achievements {
		levels, err := s.achievementRepo.getLevelsByAchievementID(ctx, achievement.id)
		if err != nil {
			logger.Error("AchievementService.GetAllAchievements achievementRepo.getLevelsByAchievementID: ", err)
			return nil, err
		}

		result = append(result, achievement.toDTO(levels))
	}

	return result, nil
}
