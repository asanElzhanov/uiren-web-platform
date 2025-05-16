package auth

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrVerificationExpired  = errors.New("verification code expired")
	ErrVerificationInvalid  = errors.New("verification code invalid")
	ErrVerificationNotFound = errors.New("verification not found")
	ErrInvalidToken         = errors.New("invalid token")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)
