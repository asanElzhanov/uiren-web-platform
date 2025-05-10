package users

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserAchievement struct {
	AchievementName  string `json:"ach_name"`
	Level            int    `json:"lvl"`
	LevelDescription string `json:"description"`
	Progress         int    `json:"progress"`
	Threshold        int    `json:"lvl_threshold"`
}

type UserProgress struct {
	Badges       []string          `json:"badges"`
	XP           int32             `json:"xp"`
	Achievements []UserAchievement `json:"achievements"`
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
	ID        string
	Username  string
	Email     string
	Password  string
	Firstname string
	Lastname  string
	Phone     string
	IsActive  bool
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
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
