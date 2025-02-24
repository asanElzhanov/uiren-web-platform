package achievements

import (
	"context"
	"errors"
	"uiren/pkg/logger"
)

type achievementRepo interface {
	createAchievement(ctx context.Context, name string) (AchievementDTO, error)                             //todo
	updateAchievement(ctx context.Context, dto UpdateAchievementDTO) (AchievementDTO, error)                //todo
	getAchievement(ctx context.Context, id int) (AchievementDTO, error)                                     //todo
	deleteAchievement(ctx context.Context, id int) error                                                    //todo
	createAchievementLevel(ctx context.Context, dto CreateAchievementLevelDTO) (AchievementLevelDTO, error) //todo
	updateAchievementLevel(ctx context.Context, dto UpdateAchievementLevelDTO) (AchievementLevelDTO, error) //todo
	getAchievementLevels(ctx context.Context, achID int) ([]AchievementLevelDTO, error)                     //todo
	deleteAchievementLevel(ctx context.Context, achID int, level int) error                                 //todo
	deleteLevelsByAchievementID(ctx context.Context, achID int) error                                       //todo
}

type userAchievementRepo interface {
	increaseUserProgress(ctx context.Context, dto UpdateUserAchievementDTO) error //todo
	decreaseUserProgress(ctx context.Context, dto UpdateUserAchievementDTO) error //todo
	tryIncreaseLevel(ctx context.Context, dto UpdateUserAchievementDTO) error     //todo
}

type AchievementService struct {
	achievementRepo     achievementRepo
	userAchievementRepo userAchievementRepo
}

func NewAchievementService(achievementRepo achievementRepo, userAchievementRepo userAchievementRepo) *AchievementService {
	return &AchievementService{
		achievementRepo:     achievementRepo,
		userAchievementRepo: userAchievementRepo,
	}
}

func (s *AchievementService) CreateAchievement(ctx context.Context, name string) (AchievementDTO, error) {
	logger.Info("AchievementService.CreateAchievement new request")

	achievement, err := s.achievementRepo.createAchievement(ctx, name)
	if err != nil {
		logger.Error("AchievementService.achievementRepo.createAchievement: ", err)
		return AchievementDTO{}, err
	}

	return achievement, nil
}

func (s *AchievementService) UpdateAchievement(ctx context.Context, dto UpdateAchievementDTO) (AchievementDTO, error) {
	logger.Info("AchievementService.UpdateAchievement new request")

	updatedAchievement, err := s.achievementRepo.updateAchievement(ctx, dto)
	if err != nil {
		logger.Error("AchievementService.achievementRepo.updateAchievement: ", err)
		return AchievementDTO{}, nil
	}

	return updatedAchievement, nil
}

func (s *AchievementService) GetAchievementInfo(ctx context.Context, id int) (AchievementDTO, error) {
	logger.Info("AchievementService.GetAchievementInfo new request")

	achievementInfo, err := s.achievementRepo.getAchievement(ctx, id)
	if err != nil {
		logger.Error("AchievementService.achievementRepo.getAchievement: ", err)
		return AchievementDTO{}, nil
	}

	achievementInfo.Levels, err = s.achievementRepo.getAchievementLevels(ctx, id)
	if err != nil {
		logger.Error("AchievementService.achievementRepo.getAchievementLevels: ", err)
		return AchievementDTO{}, nil
	}

	return achievementInfo, nil
}

func (s *AchievementService) DeleteAchievement(ctx context.Context, id int) error {
	logger.Info("AchievementService.DeleteAchievement new request")

	if err := s.achievementRepo.deleteAchievement(ctx, id); err != nil {
		logger.Error("AchievementService.achievementRepo.deleteAchievement: ", err)
		return err
	}

	if err := s.achievementRepo.deleteLevelsByAchievementID(ctx, id); err != nil {
		logger.Error("AchievementService.achievementRepo.deleteLevelsByAchievementID: ", err)
	}

	return nil
}

func (s *AchievementService) AddAchievementLevel(ctx context.Context, dto CreateAchievementLevelDTO) (AchievementLevelDTO, error) {
	logger.Info("AchievementService.AddAchievementLevel new request")

	if err := s.achievementRepo.deleteAchievementLevel(ctx, dto.AchievementID, dto.Level); err != nil && !errors.Is(err, ErrAchievementLevelNotFound) {
		logger.Error("AchievementService.AddAchievementLevel deleteAchievementLevel: ", err)
		return AchievementLevelDTO{}, err
	}

	newLevel, err := s.achievementRepo.createAchievementLevel(ctx, dto)
	if err != nil {
		logger.Error("AchievementService.AddAchievementLevel createAchievementLevel: ", err)
		return AchievementLevelDTO{}, err
	}

	return newLevel, nil
}
