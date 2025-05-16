package achievements

import (
	"time"
)

type achievement struct {
	id        int
	name      string
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

type AchievementDTO struct {
	ID        int                `json:"id"`
	Name      string             `json:"name"`
	Levels    []AchievementLevel `json:"levels"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt *time.Time         `json:"deleted_at"`
}

func (ach achievement) toDTO(levels []AchievementLevel) AchievementDTO {
	return AchievementDTO{
		ID:        ach.id,
		Name:      ach.name,
		Levels:    levels,
		CreatedAt: ach.createdAt,
		UpdatedAt: ach.updatedAt,
		DeletedAt: ach.deletedAt,
	}
}

type AchievementLevel struct {
	AchID       int       `json:"achievement_id"`
	AchName     string    `json:"achievement_name"`
	Level       int       `json:"level"`
	Description string    `json:"description"`
	Threshold   int       `json:"threshold"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LevelData struct {
	Level     int
	Threshold int
}

//for repo

type UpdateAchievementDTO struct {
	ID      int
	NewName string
}

type AddAchievementLevelDTO struct {
	AchID       int
	Description string
	Threshold   int
	Level       int
}

type DeleteAchievementLevelDTO struct {
	AchID int
	Level int
}
