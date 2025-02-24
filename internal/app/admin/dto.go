package admin

// users
type (
	User struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

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
