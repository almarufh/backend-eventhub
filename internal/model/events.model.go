package model

import "time"

type OrganizerMDL struct {
	ID     int32  `db:"organizer_id"`
	Name   string `db:"name"`
	Job    string `db:"job"`
	Office string `db:"office"`
	Image  string `db:"image"`
}

type SpeakerMDL struct {
	ID     int32  `db:"organizer_id"`
	Name   string `db:"name"`
	Job    string `db:"job"`
	Office string `db:"office"`
}

type SpeakersMDL struct {
	ID   int32 `db:"d"`
	data []SpeakerMDL
}

type DetailEvenstMDL struct {
	ID          int32        `db:"id"`
	Community   string       `db:"community"`
	Organizer   OrganizerMDL `db:"-"`
	Speakers    SpeakersMDL  `db:"-"`
	Title       string       `db:"title"`
	Location    string       `db:"location"`
	Description string       `db:"description"`
	Image       string       `db:"image"`
	Capacity    int32        `db:"capacity"`
	StartTime   time.Time    `db:"start_time"`
	Attendee    int32        `db:"attendee"`
	EndTime     time.Time    `db:"end_time"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
}
