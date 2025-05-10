package admin

// users
type (
	CreateUserReq struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	UpdateUserReq struct {
		Firstname   string `json:"first_name"`
		Lastname    string `json:"last_name"`
		Phone       string `json:"phone"`
		PhoneRegion string `json:"phone_region"`
	}
)

// auth
type (
	SignInParams struct {
		Identificator string `json:"identificator"`
		Password      string `json:"password"`
	}
)

// achievements
type (
	CreateAchievementRequest struct {
		Name string `json:"name"`
	}

	UpdateAchievementRequest struct {
		ID      int    `json:"id"`
		NewName string `json:"name"`
	}
	AddAchievementLevelRequest struct {
		AchID       int    `json:"achievement_id"`
		Description string `json:"description"`
		Threshold   int    `json:"threshold"`
	}
	DeleteAchievementLevelRequest struct {
		AchID int `json:"achievement_id"`
		Level int `json:"level"`
	}
)
