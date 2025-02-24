package jwt_maker

import "time"

type PayloadDTO struct {
	Username  string
	Firstname string
	Lastname  string
	IsAdmin   bool
	Duration  time.Duration
}
