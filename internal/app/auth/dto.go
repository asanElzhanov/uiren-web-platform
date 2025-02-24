package auth

import "uiren/internal/app/users"

type LoginParams struct {
	Identificator string `json:"ident"`
	Password      string `json:"pass"`
}

type RegisterParams struct {
	DTO users.CreateUserDTO
}
