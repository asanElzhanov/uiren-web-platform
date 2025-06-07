package modules

import (
	"time"
	"uiren/internal/app/lessons"
)

type Module struct {
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

type ModuleWithLessons struct {
	Code        string              `json:"code"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Goal        string              `json:"goal"`
	Difficulty  string              `json:"difficulty"`
	UnlockReq   UnlockRequirements  `json:"unlock_requirements"`
	Reward      Reward              `json:"reward"`
	Lessons     []lessons.LessonDTO `json:"lessons"`
	CreatedAt   time.Time           `json:"created_at"`
	DeletedAt   time.Time           `json:"deleted_at"`
}

func (module Module) toDTO(lessons []lessons.LessonDTO) ModuleWithLessons {
	var moduleDeletedAt time.Time
	if module.DeletedAt != nil {
		moduleDeletedAt = *module.DeletedAt
	}
	return ModuleWithLessons{
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
	PrevModuleCode string  `bson:"previous_module" json:"previous_module"`
	MinimumXP      float64 `bson:"min_xp" json:"min_xp"`
}

type Reward struct {
	XP    float64 `bson:"xp" json:"xp"`
	Badge string  `bson:"badge" json:"badge"`
}

// for repo

type CreateModuleDTO struct {
	Code        string             `bson:"code" json:"code"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Goal        string             `bson:"goal" json:"goal"`
	Difficulty  string             `bson:"difficulty" json:"difficulty"`
	UnlockReq   UnlockRequirements `bson:"unlock_requirements" json:"unlock_requirements"`
	Reward      Reward             `bson:"reward" json:"reward"`
	Lessons     []string           `bson:"lessons"`
	CreatedAt   time.Time          `bson:"created_at"`
	DeletedAt   *time.Time         `bson:"deleted_at"`
}

type UpdateModuleDTO struct {
	Title       *string             `bson:"title,omitempty" json:"title,omitempty"`
	Description *string             `bson:"description,omitempty" json:"description,omitempty"`
	Goal        *string             `bson:"goal,omitempty" json:"goal,omitempty"`
	Difficulty  *string             `bson:"difficulty,omitempty" json:"difficulty,omitempty"`
	UnlockReq   *UnlockRequirements `bson:"unlock_requirements,omitempty" json:"unlock_requirements,omitempty"`
	Reward      *Reward             `bson:"reward,omitempty" json:"reward,omitempty"`
}
