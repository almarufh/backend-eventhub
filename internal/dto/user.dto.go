package dto

import (
	"mime/multipart"
	"time"
)

type ReqNewPassword struct {
	ID           int32  `json:"-"`
	Old_password string `json:"old_password" binding:"min=8"`
	New_password string `json:"new_password" binding:"min=8"`
}

type ResMyProfle struct {
	ID             int32     `json:"-"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	Role           string    `json:"role"`
	DarkPreference bool      `json:"dark_preference"`
	Address        *string   `json:"address"`
	Job            *string   `json:"job"`
	Office         *string   `json:"office"`
	Image          *string   `json:"image"`
	Description    *string   `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
}

type ReqSetProfile struct {
	Name        *string               `form:"name" example:"Alma'ruf Hidayat" binding:"omitempty,min=3,max=100"`
	Address     *string               `form:"address" example:"Masamba, Luwu Utara, Sulawesi Selatan" binding:"omitempty,max=255"`
	Job         *string               `form:"job" example:"Fullstack Developer" binding:"omitempty,max=100"`
	Office      *string               `form:"office" example:"PT. Best Life Ummah" binding:"omitempty,max=100"`
	Image       *multipart.FileHeader `form:"image" swaggerignore:"true" binding:"omitempty"`
	Description *string               `form:"description" example:"Senior Fullstack Developer" binding:"omitempty,max=500"`
}

type ResSetProfile struct {
	Name        *string   `json:"name"`
	Address     *string   `json:"address"`
	Job         *string   `json:"job"`
	Office      *string   `json:"office"`
	Image       *string   `json:"image"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type ResJoinedEvents struct {
	Total uint32           `json:"total"`
	Data  []ResDetailEvent `json:"data"`
}

type ResSavedEvents struct {
	Total uint32           `json:"total"`
	Data  []ResDetailEvent `json:"data"`
}

type ResJoinedCommunities struct {
	Total uint32                `json:"total"`
	Data  []ResDetailCommunitiy `json:"data"`
}
