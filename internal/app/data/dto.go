package data

import (
	"uiren/internal/app/modules"
	"uiren/internal/app/users"
)

type UserInfo struct {
	ID        string              `json:"id"`
	Username  string              `json:"username"`
	Firstname string              `json:"firstname"`
	Lastname  string              `json:"lastname"`
	Email     string              `json:"email"`
	Phone     string              `json:"phone"`
	Progress  *users.UserProgress `json:"progress"`
}

type ModulesList struct {
	Modules []modules.Module `json:"modules"`
	Total   int              `json:"total"`
}
