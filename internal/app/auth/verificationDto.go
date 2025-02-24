package auth

import "time"

type CreateVerificationCodeRequest struct {
	Username string
	Email    string
	Code     string
	Duration time.Duration
}

type Verification struct {
	Username  string
	Email     string
	Code      string
	ExpiresAt time.Time
}
