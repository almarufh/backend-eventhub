package repo

import (
	msgerr "backend/EventHub/internal/message"
	"backend/EventHub/internal/model"
	"backend/EventHub/pkg"
	"context"
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventsRepo struct {
	db *pgxpool.Pool
}

func NewEventsRepo(db *pgxpool.Pool) *EventsRepo {
	return &EventsRepo{
		db: db,
	}
}

func (er *EventsRepo) GetDetailEvents(ctx context.Context, id int32) (model.DetailEvenstMDL, error) {
	query := `
        SELECT 
            e.id,
            c.title AS community,
            e.organizer_id,
            COUNT(je.user_id) AS attendee,
            e.title,
            e.location,
            e.description,
            e.image,
            e.capacity,
            e.start_time,
            e.end_time,
            e.created_at,
            e.updated_at
        FROM events e
        JOIN communities c ON c.id = e.community_id
        LEFT JOIN joined_events_users je ON je.event_id = e.id
        WHERE e.id = $1
        GROUP BY e.id, c.title
    `
	var event model.DetailEvenstMDL
	err := er.db.QueryRow(ctx, query, id).Scan(
		&event.ID,
		&event.Community,
		&event.Organizer.ID,
		&event.Attendee,
		&event.Title,
		&event.Location,
		&event.Description,
		&event.Image,
		&event.Capacity,
		&event.StartTime,
		&event.EndTime,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		parsedErr := pkg.ParseError(err)
		if errors.Is(parsedErr, pkg.ErrNotFound) {
			return model.DetailEvenstMDL{}, errors.New("event not found")
		}
		return model.DetailEvenstMDL{}, parsedErr
	}
	return event, nil
}

func (er *EventsRepo) GetSpeakersByEventID(ctx context.Context, eventID int32) ([]model.SpeakerMDL, error) {
	query := `
		SELECT
			s.name,
			s.job,
			s.office
		FROM speakers_events se
		JOIN speakers s ON s.id = se.speaker_id
		WHERE se.event_id = $1
	`

	rows, err := er.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, pkg.ParseError(err)
	}
	defer rows.Close()

	var speakers []model.SpeakerMDL

	for rows.Next() {
		var speaker model.SpeakerMDL
		if err := rows.Scan(
			&speaker.Name,
			&speaker.Job,
			&speaker.Office,
		); err != nil {
			return nil, err
		}
		speakers = append(speakers, speaker)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return speakers, nil
}

func (er *EventsRepo) GetCategoryEvents(ctx context.Context, eventID int32) ([]string, error) {
	query := `
		SELECT
			c.name
		FROM categories c
		JOIN events_categories ec ON c.id = ec.category_id
		WHERE ec.event_id = $1
	`

	rows, err := er.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, pkg.ParseError(err)
	}
	defer rows.Close()

	var categories []string

	for rows.Next() {
		var category string
		if err := rows.Scan(
			&category,
		); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (er *EventsRepo) JoinEventRepo(ctx context.Context, userID int32, eventID int32) error {
	query := `
        INSERT INTO joined_events_users (user_id, event_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, event_id) DO NOTHING;
    `
	rows, err := er.db.Exec(ctx, query, userID, eventID)
	if err != nil {
		return err
	}

	if rows.RowsAffected() == 0 {
		return msgerr.AlredyJoinedEvents
	}

	return nil
}

func (er *EventsRepo) SaveEventRepo(ctx context.Context, userID int32, eventID int32) error {
	query := `
        INSERT INTO saved_events_users (user_id, event_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, event_id) DO NOTHING;
    `
	rows, err := er.db.Exec(ctx, query, userID, eventID)
	if err != nil {
		return err
	}

	if rows.RowsAffected() == 0 {
		return msgerr.AlredySavedEvent
	}

	return nil
}

func (er *EventsRepo) LeaveEventRepo(ctx context.Context, userID int32, eventID int32) error {
	query := `
		DELETE FROM joined_events_users 
		WHERE user_id = $1 AND event_id = $2;
	`
	rows, err := er.db.Exec(ctx, query, userID, eventID)
	if err != nil {
		return err
	}

	if rows.RowsAffected() == 0 {
		return msgerr.NotJoinedEvents
	}

	return nil
}

func (er *EventsRepo) UnsaveEventRepo(ctx context.Context, userID int32, eventID int32) error {
	query := `
		DELETE FROM saved_events_users 
		WHERE user_id = $1 AND event_id = $2;
	`
	rows, err := er.db.Exec(ctx, query, userID, eventID)
	if err != nil {
		return err
	}

	if rows.RowsAffected() == 0 {
		return msgerr.NotSavedEvent
	}

	return nil
}

func (er *EventsRepo) GetAllEvents(ctx context.Context, filter EventFilter) ([]model.DetailEvenstMDL, PaginationMeta, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	baseQuery := `
        SELECT 
            e.id,
            c.title AS community,
            e.organizer_id,
            COUNT(DISTINCT je_all.user_id) AS attendee,
            e.title,
            e.location,
            e.description,
            e.image,
            e.capacity,
            e.start_time,
            e.end_time,
            e.created_at,
            e.updated_at,
            COUNT(*) OVER() AS total_count
        FROM events e
        JOIN communities c ON c.id = e.community_id
        LEFT JOIN categories cat ON cat.id = e.category_id
        LEFT JOIN joined_events_users je_all ON je_all.event_id = e.id
        WHERE 1=1
    `
	var args []any
	argIdx := 1

	if filter.Category != "" {
		baseQuery += fmt.Sprintf(" AND cat.name ILIKE $%d", argIdx)
		args = append(args, filter.Category)
		argIdx++
	}

	if filter.Location != "" {
		baseQuery += fmt.Sprintf(" AND e.location ILIKE $%d", argIdx)
		args = append(args, "%"+filter.Location+"%")
		argIdx++
	}

	if filter.Search != "" {
		baseQuery += fmt.Sprintf(" AND e.title ILIKE $%d", argIdx)
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}

	baseQuery += ` GROUP BY e.id, c.title`

	switch filter.SortBy {
	case "most_popular":
		baseQuery += ` ORDER BY attendee DESC, e.start_time ASC`

	case "almost_full":
		baseQuery += ` ORDER BY (COUNT(DISTINCT je_all.user_id)::float / NULLIF(e.capacity, 0)) DESC`

	case "recently_added":
		baseQuery += ` ORDER BY e.created_at DESC`

	case "upcoming":
		fallthrough
	default:
		baseQuery += ` AND e.start_time >= NOW() ORDER BY e.start_time ASC`
	}

	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, filter.Limit, offset)

	rows, err := er.db.Query(ctx, baseQuery, args...)
	if err != nil {
		log.Printf("[UserRepo.GetAllEvents] Query Error: %v", err)
		return nil, PaginationMeta{}, fmt.Errorf("gagal mengambil data events: %w", err)
	}
	defer rows.Close()

	events := make([]model.DetailEvenstMDL, 0)
	var totalItems int64 = 0

	for rows.Next() {
		var event model.DetailEvenstMDL
		err := rows.Scan(
			&event.ID,
			&event.Community,
			&event.Organizer.ID,
			&event.Attendee,
			&event.Title,
			&event.Location,
			&event.Description,
			&event.Image,
			&event.Capacity,
			&event.StartTime,
			&event.EndTime,
			&event.CreatedAt,
			&event.UpdatedAt,
			&totalItems,
		)
		if err != nil {
			return nil, PaginationMeta{}, fmt.Errorf("gagal memindai data event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, PaginationMeta{}, fmt.Errorf("kesalahan membaca baris event: %w", err)
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(filter.Limit)))
	}

	meta := PaginationMeta{
		CurrentPage: filter.Page,
		TotalPages:  totalPages,
		Limit:       filter.Limit,
		TotalItems:  totalItems,
	}

	return events, meta, nil
}

type EventFilter struct {
	Category string
	Location string
	Search   string
	SortBy   string
	Page     int
	Limit    int
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	Limit       int   `json:"limit"`
	TotalItems  int64 `json:"total_items"`
}

type PaginatedEventsResponse struct {
	Events []model.DetailEvenstMDL `json:"events"`
	Meta   PaginationMeta          `json:"meta"`
}
