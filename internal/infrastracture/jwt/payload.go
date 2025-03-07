package jwt_maker

import "time"

type PayloadDTO struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	IsAdmin   bool   `json:"isAdmin"`
	Duration  time.Duration
}
