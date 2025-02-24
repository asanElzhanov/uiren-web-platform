package achievements

import "errors"

var (
	ErrAchievementNotFound      = errors.New("achievement not found")
	ErrAchievementLevelNotFound = errors.New("achievement level not found")
)
