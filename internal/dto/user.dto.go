package dto

type ReqNewPassword struct {
	ID           int32
	Old_password string `json:"old_password" binding:"min=8"`
	New_password string `json:"new_password" binding:"min=8"`
}

type ResMyProfle struct {
	ID             int32   `json:"user_id"`
	Email          string  `json:"email"`
	Name           string  `json:"name"`
	Role           string  `json:"role"`
	DarkPreference bool    `json:"dark_preference"`
	Address        *string `json:"address"`
	Job            *string `json:"job"`
	Office         *string `json:"office"`
	Image          *string `json:"image"`
	Description    *string `json:"description"`
}

type DtoSetProfile struct {
	Name        string  `json:"name"`
	Address     *string `json:"address"`
	Job         *string `json:"job"`
	Office      *string `json:"office"`
	Image       *string `json:"image"`
	Description *string `json:"description"`
}
