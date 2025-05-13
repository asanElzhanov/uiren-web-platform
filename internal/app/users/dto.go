package users

import (
	"database/sql"
	"time"
	"uiren/internal/app/progress"

	"github.com/google/uuid"
)

type UserProgress struct {
	Badges       []string                   `json:"badges"`
	XP           int32                      `json:"xp"`
	Achievements []progress.UserAchievement `json:"achievements"`
}

type user struct {
	id        uuid.UUID
	username  string
	email     string
	password  string
	firstname sql.NullString
	lastname  sql.NullString
	phone     sql.NullString
	isActive  bool
	isAdmin   bool
	createdAt time.Time
	updatedAt time.Time
	deletedAt sql.NullTime
}

type UserDTO struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Firstname string    `json:"first_name"`
	Lastname  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	IsActive  bool      `json:"is_active"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (user user) ToDTO() UserDTO {
	var deletedAt time.Time
	if user.deletedAt.Valid {
		deletedAt = user.deletedAt.Time
	}
	return UserDTO{
		ID:        user.id.String(),
		Username:  user.username,
		Email:     user.email,
		Password:  user.password,
		Firstname: user.firstname.String,
		Lastname:  user.lastname.String,
		Phone:     user.phone.String,
		IsActive:  user.isActive,
		IsAdmin:   user.isAdmin,
		CreatedAt: user.createdAt,
		UpdatedAt: user.updatedAt,
		DeletedAt: deletedAt,
	}
}

type CreateUserDTO struct {
	Username, Email, Password string
}

type UpdateUserDTO struct {
	ID, Firstname, Lastname, Phone, PhoneRegion string
}
