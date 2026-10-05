package service

import (
	"context"
	"fmt"

	"backend/EventHub/internal/dto"
	"backend/EventHub/internal/repo"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventsService struct {
	er *repo.EventsRepo
	ar *repo.AuthRepo
	db *pgxpool.Pool
}

func NewEventsService(er *repo.EventsRepo, ar *repo.AuthRepo, db *pgxpool.Pool) *EventsService {
	return &EventsService{
		er: er,
		ar: ar,
		db: db,
	}
}

func (es *EventsService) GetDetailEvents(ctx context.Context, param int32) (dto.ResDetailEvent, error) {
	event, err := es.er.GetDetailEvents(ctx, param)
	if err != nil {
		return dto.ResDetailEvent{}, err
	}

	organizer, err := es.ar.GetUserById(ctx, es.db, event.Organizer.ID)
	if err != nil {
		return dto.ResDetailEvent{}, err
	}

	speakers, err := es.er.GetSpeakersByEventID(ctx, param)
	if err != nil {
		return dto.ResDetailEvent{}, err
	}

	speakersDTO := make([]dto.Speaker, 0, len(speakers))
	for _, s := range speakers {
		speakersDTO = append(speakersDTO, dto.Speaker{
			Name:   s.Name,
			Job:    s.Job,
			Office: s.Office,
		})
	}

	categories, err := es.er.GetCategoryEvents(ctx, event.Organizer.ID)
	if err != nil {
		return dto.ResDetailEvent{}, err
	}

	res := dto.ResDetailEvent{
		ID:        event.ID,
		Community: event.Community,
		Organizer: dto.Organizer{
			Name:   organizer.Name,
			Job:    *organizer.Job,
			Office: *organizer.Office,
			Image:  *organizer.Image,
		},
		Speakers:    speakersDTO,
		Categorys:   categories,
		Attendee:    event.Attendee,
		Title:       event.Title,
		Location:    event.Location,
		Description: event.Description,
		Image:       event.Image,
		Capacity:    event.Capacity,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
	}

	return res, nil
}

func (es *EventsService) JoinEventService(ctx context.Context, ID int32, eventID int32) error {
	event, err := es.er.GetDetailEvents(ctx, eventID)
	if err != nil {
		return err
	}

	if event.Attendee == event.Capacity {
		return fmt.Errorf("Event full capacity %d/%d", event.Attendee, event.Capacity)
	}

	auth, err := es.ar.GetAuthUser(ctx, es.db, ID)
	if err != nil {
		return err
	}
	return es.er.JoinEventRepo(ctx, auth.User_id, eventID)
}

func (es *EventsService) SaveEventService(ctx context.Context, ID int32, eventID int32) error {
	auth, err := es.ar.GetAuthUser(ctx, es.db, ID)
	if err != nil {
		return err
	}
	return es.er.SaveEventRepo(ctx, auth.User_id, eventID)
}

func (es *EventsService) LeaveEventService(ctx context.Context, ID int32, eventID int32) error {
	auth, err := es.ar.GetAuthUser(ctx, es.db, ID)
	if err != nil {
		return err
	}
	return es.er.LeaveEventRepo(ctx, auth.User_id, eventID)
}

func (es *EventsService) UnsaveEventService(ctx context.Context, ID int32, eventID int32) error {
	auth, err := es.ar.GetAuthUser(ctx, es.db, ID)
	if err != nil {
		return err
	}
	return es.er.UnsaveEventRepo(ctx, auth.User_id, eventID)
}
