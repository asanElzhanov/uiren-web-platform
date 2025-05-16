package exercises

import "errors"

var (
	ErrNotFound              = errors.New("exercise not found")
	ErrCodeAlreadyExists     = errors.New("code already exists")
	ErrIncorrectType         = errors.New("incorrect exercise type")
	ErrOptionsRequired       = errors.New("correct options required")
	ErrCorrectAnswerRequired = errors.New("correct answer required")
	ErrPairsRequired         = errors.New("correct pairs required")
	ErrCorrectOrderRequired  = errors.New("correct order required")
	ErrNoFieldsToUpdate      = errors.New("no fields to update")
)
