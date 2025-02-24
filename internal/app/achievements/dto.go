package achievements

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type (
	achievement struct {
		id        int
		name      string
		levels    []achievementLevel
		createdAt time.Time
		updatedAt time.Time
		deletedAt sql.NullTime
	}

	AchievementDTO struct {
		ID        int
		Name      string
		Levels    []AchievementLevelDTO
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt time.Time
	}
)

func (a achievement) toDTO() AchievementDTO {
	var deletedAt time.Time
	if a.deletedAt.Valid {
		deletedAt = a.deletedAt.Time
	}

	levelsList := make([]AchievementLevelDTO, 0, len(a.levels))

	for _, level := range a.levels {
		levelsList = append(levelsList, level.toDTO())
	}

	return AchievementDTO{
		ID:        a.id,
		Name:      a.name,
		Levels:    levelsList,
		CreatedAt: a.createdAt,
		UpdatedAt: a.updatedAt,
		DeletedAt: deletedAt,
	}
}

type (
	achievementLevel struct {
		id          int
		level       int
		description string
		toComplete  int
		createdAt   time.Time
		updatedAt   time.Time
		deletedAt   sql.NullTime
	}

	AchievementLevelDTO struct {
		ID          int
		Level       int
		Description string
		ToComplete  int
		CreatedAt   time.Time
		UpdatedAt   time.Time
		DeletedAt   time.Time
	}
)

func (l achievementLevel) toDTO() AchievementLevelDTO {
	var deletedAt time.Time
	if l.deletedAt.Valid {
		deletedAt = l.deletedAt.Time
	}

	return AchievementLevelDTO{
		ID:          l.id,
		Level:       l.level,
		Description: l.description,
		ToComplete:  l.toComplete,
		CreatedAt:   l.createdAt,
		UpdatedAt:   l.updatedAt,
		DeletedAt:   deletedAt,
	}
}

type (
	userAchievement struct {
		userID          uuid.UUID
		achievementName string
		level           int
		progress        float64
		createdAt       time.Time
		updatedAt       time.Time
		deletedAt       sql.NullTime
	}

	UserAchievementDTO struct {
		UserID          string
		AchievementName string
		Level           int
		Progress        float64
		CreatedAt       time.Time
		UpdatedAt       time.Time
		DeletedAt       time.Time
	}
)

func (ua userAchievement) toDTO() UserAchievementDTO {
	var deletedAt time.Time
	if ua.deletedAt.Valid {
		deletedAt = ua.deletedAt.Time
	}

	return UserAchievementDTO{
		UserID:          ua.userID.String(),
		AchievementName: ua.achievementName,
		Level:           ua.level,
		Progress:        ua.progress,
		CreatedAt:       ua.createdAt,
		UpdatedAt:       ua.updatedAt,
		DeletedAt:       deletedAt,
	}
}

type UpdateAchievementDTO struct {
	ID      int
	NewName string
}

type CreateAchievementLevelDTO struct {
	AchievementID int
	Level         int
	Description   string
	ToComplete    int
}

type UpdateAchievementLevelDTO struct {
	Description string
	ToComplete  int
}

type UpdateUserAchievementDTO struct {
	UserID, AchievementID, Progress int
}
