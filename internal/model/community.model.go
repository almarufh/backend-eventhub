package model

type JoinCommunityMDL struct {
	ID     int32  `db:"user_id"`
	Title  string `db:"title"`
	Status string `db:"status"`
}

type CommunitieMDL struct {
	ID         int32    `db:"user_id"`
	Categories []string `db:"title"`
	Status     string   `db:"status"`
}

type CategorieMDL struct {
	ID         int32    `db:"id"`
	Categories []string `db:"name"`
}

type DetailCommunitieMDL struct {
	ID          int32  `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Image       string `db:"image"`
	Status      string `db:"status"`
	Categories  string `db:"categories"`
	Members     int32  `db:"members"`
}

type MembersCommunityMDL struct {
	Name   string  `db:"name"`
	Job    *string `db:"job"`
	Office *string `db:"office"`
	Image  *string `db:"image"`
}
