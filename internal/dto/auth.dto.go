package dto

type ReqRegister struct {
	Name     string `json:"name" example:"Alma'ruf Hidayat"`
	Email    string `json:"email" binding:"required,email" example:"almarufhidayat99@gmail.com"`
	Password string `json:"password" binding:"min=8" example:"Admin@1234"`
	Role     string `json:"-"`
}

type ResRegister struct {
	ID   int32  `json:"user_id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type ReqLogin struct {
	Email    string `json:"email" binding:"required,email" example:"almarufhidayat99@gmail.com"`
	Password string `json:"password" binding:"min=8" example:"Admin@1234"`
	Token    string `json:"-"`
	Role     string `json:"-"`
	Device   string `json:"-"`
}

type ResLogin struct {
	Token          string  `json:"token"`
	Name           string  `json:"name"`
	Role           string  `json:"role"`
	DarkPreference bool    `json:"dark_preference"`
	Image          *string `json:"image"`
}

type ReqChangePassword struct {
	ID       int32  `json:"-"`
	Email    string `json:"email" binding:"required,email" example:"almarufhidayat99@gmail.com"`
	Password string `json:"password" binding:"min=8"  example:"Admin@1234"`
}
