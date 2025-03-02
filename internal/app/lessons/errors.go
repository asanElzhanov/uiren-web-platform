package lessons

import "errors"

var (
	ErrCodeAlreadyExists    = errors.New("code already exists")
	ErrNoFieldsToUpdate     = errors.New("no fields to update")
	ErrNotFound             = errors.New("lesson not found")
	ErrExerciseAlreadyInSet = errors.New("exercise already in list")
	ErrExerciseNotInList    = errors.New("exercise is not in list")
)
