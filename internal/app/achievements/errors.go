package achievements

import "errors"

var (
	ErrAchievementNotFound      = errors.New("achievement not found")
	ErrAchievementLevelNotFound = errors.New("achievement level not found")
	ErrAchievementNameExists    = errors.New("achievement name already exists")
	ErrLowThreshold             = errors.New("achievement threshold is lower than maximum")
	ErrInvalidThreshold         = errors.New("invalid threshold")
	ErrLevelExists              = errors.New("achievement level already exists")
)
