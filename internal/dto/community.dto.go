package dto

type ResJoinCommunity struct {
	ID     int32  `json:"user_id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type ReqJoinCommunity struct {
	ID int32 `json:"id"`
}

type ResDetailCommunitiy struct {
	ID          int32  `json:"community_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Status      string `json:"status"`
	Categories  string `json:"categories"`
	Members     int32  `json:"members"`
}

type ResMembersCommunity struct {
	Name   string  `json:"name"`
	Job    *string `json:"job"`
	Office *string `json:"office"`
	Image  *string `json:"image"`
}
