package users

import "errors"

var (
	ErrIncorrectEmail    = errors.New("incorrect email format")
	ErrIncorrectPassword = errors.New("incorrect password format")
	ErrIncorrectPhone    = errors.New("incorrect phone format")
	ErrUserNotFound      = errors.New("user not found")
	ErrUsernameExists    = errors.New("username already exists")
	ErrEmailExists       = errors.New("email already exists")
)
