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
	ID        int
	Name      string
	Levels    []AchievementLevel
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
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
	achID       int
	achName     string
	level       int
	description string
	threshold   int
	createdAt   time.Time
	updatedAt   time.Time
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
