package progress

import "errors"

var (
	ErrBadgeNotExists              = errors.New("badge does not exist")
	ErrBadgeAlreadyExists          = errors.New("badge already exists")
	ErrBadgeNotProvided            = errors.New("badge not provided")
	ErrUserHasBadge                = errors.New("user already has badge")
	ErrAchievementProgressNotFound = errors.New("achievement progress not found")
	ErrNegativeProgress            = errors.New("negative achievement progress not allowed")
)
