package dto

type ReqRegister struct {
	Name     string `json:"name" example:"Alma'ruf Hidayat"`
	Email    string `json:"email" example:"almarufhidayat99@gmail.com"`
	Password string `json:"password" example:"Admin@1234"`
	Role     string `json:"role"`
}

type ResRegister struct {
	ID   int32  `json:"user_id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type ReqLogin struct {
	Email    string `json:"email" binding:"required,email" example:"almarufhidayat99@gmail.com"`
	Password string `json:"password" binding:"min=8" example:"Admin@12345"`
	Token    string `json:"-"`
	Role     string `json:"-"`
	Device   string `json:"-"`
}

type ResLogin struct {
	Token string `json:"token"`
}

type ReqChangePassword struct {
	ID               int32
	Email            string `json:"email" example:"almarufhidayat99@gmail.com"`
	New_password     string `json:"new_password" example:"Admin21234"`
	Confirm_password string `json:"confirm_password" example:"Admin21234"`
}
