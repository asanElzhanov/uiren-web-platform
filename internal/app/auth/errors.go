package auth

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrVerificationExpired  = errors.New("verification code expired")
	ErrVerificationInvalid  = errors.New("verification code invalid")
	ErrInvalidToken         = errors.New("invalid token")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)
