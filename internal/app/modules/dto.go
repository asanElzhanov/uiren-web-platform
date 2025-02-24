package modules

import (
	"time"
	"uiren/internal/app/lessons"
)

type Difficulty string

var BeginnerDifficulty Difficulty = "beginner"
var MediumDifficulty Difficulty = "medium"
var HighDifficulty Difficulty = "high"

type module struct {
	Code        string             `bson:"code" json:"code"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Goal        string             `bson:"goal" json:"goal"`
	Difficulty  string             `bson:"difficulty" json:"difficulty"`
	UnlockReq   UnlockRequirements `bson:"unlock_requirements" json:"unlock_requirements"`
	Reward      Reward             `bson:"reward" json:"reward"`
	Lessons     []string           `bson:"lessons" json:"lessons"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	DeletedAt   *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type ModuleDTO struct {
	Code        string
	Title       string
	Description string
	Goal        string
	Difficulty  string
	UnlockReq   UnlockRequirements
	Reward      Reward
	Lessons     []lessons.LessonDTO
	CreatedAt   time.Time
	DeletedAt   time.Time
}

func (module module) toDTO(lessons []lessons.LessonDTO) ModuleDTO {
	var moduleDeletedAt time.Time
	if module.DeletedAt != nil {
		moduleDeletedAt = *module.DeletedAt
	}
	return ModuleDTO{
		Code:        module.Code,
		Title:       module.Title,
		Description: module.Description,
		Goal:        module.Goal,
		Difficulty:  module.Difficulty,
		UnlockReq:   module.UnlockReq,
		Reward:      module.Reward,
		Lessons:     lessons,
		CreatedAt:   module.CreatedAt,
		DeletedAt:   moduleDeletedAt,
	}
}

type UnlockRequirements struct {
	PrevModuleCode string
	MinimumXP      float64
}

type Reward struct {
	XP    float64
	Badge string
}

// for repo

type CreateModuleDTO struct {
	Code        string
	Title       string
	Description string
	Goal        string
	Difficulty  Difficulty
	UnlockReq   UnlockRequirements
	Reward      Reward
}

type UpdateModuleDTO struct {
	Title       string
	Description string
	Goal        string
	Difficulty  string
	UnlockReq   UnlockRequirements
	Reward      Reward
}
