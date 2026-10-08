package model

import "time"

type ProfileMDL struct {
	Name        *string `db:"name"`
	Address     *string `db:"address"`
	Job         *string `db:"job"`
	Office      *string `db:"office"`
	Image       *string `db:"image"`
	Description *string `db:"description"`
}

type UsersDB struct {
	ID        int32     `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ProfileDB struct {
	UserID         int32     `db:"user_id"`
	Name           string    `db:"name"`
	Role           string    `db:"role"`
	DarkPreference bool      `db:"dark_preference"`
	Address        *string   `db:"address"`
	Job            *string   `db:"job"`
	Office         *string   `db:"office"`
	Image          *string   `db:"image"`
	Description    *string   `db:"description"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type AuthDB struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	IsActive  bool      `db:"is_active"`
	Token     string    `db:"token"`
	Device    *string   `db:"device"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
