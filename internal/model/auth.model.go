package model

import "time"

type RegisterMDL struct {
	ID   int32  `db:"user_id"`
	Name string `db:"name"`
	Role string `db:"role"`
}

type LoginMDL struct {
	ID      int32  `db:"id"`
	User_id int32  `db:"user_id"`
	Token   string `db:"email"`
}

type UserMDL struct {
	ID             int32     `db:"user_id"`
	Email          string    `db:"email"`
	Password       string    `db:"password"`
	Name           string    `db:"name"`
	Role           string    `db:"role"`
	DarkPreference bool      `db:"dark_preference"`
	Address        *string   `db:"address"`
	Job            *string   `db:"job"`
	Office         *string   `db:"office"`
	Image          *string   `db:"image"`
	Description    *string   `db:"description"`
	CreatedAt      time.Time `db:"created_at"`
}

type AuthUser struct {
	// ID       int32   `db:"id"`
	User_id  int32  `db:"user_id"`
	IsActive bool   `db:"is_active"`
	Token    string `db:"token"`
	// Device   *string `db:"device"`
}

type AuthChangePassword struct {
	Email    string `db:"email"`
	Password string `db:"password"`
}
