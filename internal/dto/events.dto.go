package dto

import "time"

type Organizer struct {
	Name   string `json:"name"`
	Job    string `json:"job"`
	Office string `json:"office"`
	Image  string `json:"image"`
}

type Speaker struct {
	Name   string `json:"name"`
	Job    string `json:"job"`
	Office string `json:"office"`
}

type ResDetailEvent struct {
	ID          int32     `json:"id"`
	Community   string    `json:"community"`
	Organizer   Organizer `json:"organizer"`
	Speakers    []Speaker `json:"speakers"`
	Categorys   []string  `json:"categorys"`
	Attendee    int32     `json:"attendee"`
	Title       string    `json:"title"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Capacity    int32     `json:"capacity"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}
