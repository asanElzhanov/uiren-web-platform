package modules

import "errors"

var (
	ErrNotFound           = errors.New("module not found")
	ErrCodeAlreadyExists  = errors.New("code already exists")
	ErrNoFieldsToUpdate   = errors.New("no fields to update")
	ErrInvalidCode        = errors.New("invalid code")
	ErrLessonAlreadyInSet = errors.New("lesson already in list")
	ErrLessonNotInList    = errors.New("lesson is not in list")
)
