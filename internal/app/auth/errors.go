package auth

import "errors"

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrVerificationExpired = errors.New("verification code expired")
	ErrVerificationInvalid = errors.New("verification code invalid")
)
