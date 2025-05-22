package data

import (
	"time"
	"uiren/internal/app/modules"
	"uiren/internal/app/progress"
	"uiren/internal/app/users"
)

// for service
type UserInfo struct {
	ID        string              `json:"id"`
	Username  string              `json:"username"`
	Firstname string              `json:"first_name"`
	Lastname  string              `json:"last_name"`
	Email     string              `json:"email"`
	Phone     string              `json:"phone"`
	Progress  *users.UserProgress `json:"progress"`
	CreatedAt time.Time           `json:"created_at"`
}

type ModulesList struct {
	Modules []modules.Module `json:"modules"`
	Total   int              `json:"total"`
}

type XPLeaderboard struct {
	Board progress.XPLeaderboard
}
