package repo

import (
	msgerr "backend/EventHub/internal/message"
	"backend/EventHub/internal/model"
	"backend/EventHub/pkg"
	"context"
	"errors"

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

// func (er *EventsRepo) GetAllEventsRepo(ctx context.Context, search string, categories []string, page int) ([]model.DetailEvenstMDL, error) {

// 	limit := 6

// 	if page < 1 {
// 		page = 1
// 	}

// 	offset := (page - 1) * limit

// 	query := `
// 		SELECT
// 			e.id,
// 			c.title AS community,
// 			e.organizer_id,

// 			(
// 				SELECT COUNT(*)
// 				FROM joined_events_users je
// 				WHERE je.event_id = e.id
// 			) AS attendee,

// 			e.title,
// 			e.location,
// 			e.description,
// 			e.image,
// 			e.capacity,
// 			e.start_time,
// 			e.end_time,
// 			e.created_at,
// 			e.updated_at

// 		FROM events e
// 		JOIN communities c
// 			ON c.id = e.community_id
// 	`

// 	var args []interface{}
// 	var conditions []string
// 	argPosition := 1

// 	if search != "" {
// 		conditions = append(conditions, fmt.Sprintf(`
// 			(
// 				POSITION(LOWER($%d) IN LOWER(e.title)) > 0
// 				OR
// 				POSITION(LOWER($%d) IN LOWER(c.title)) > 0
// 				OR
// 				POSITION(LOWER($%d) IN LOWER(e.location)) > 0
// 			)
// 		`, argPosition, argPosition, argPosition))

// 		args = append(args, search)
// 		argPosition++
// 	}

// 	if len(categories) > 0 {
// 		conditions = append(conditions, fmt.Sprintf(`
// 			(
// 				SELECT COUNT(DISTINCT ca.id)
// 				FROM events_categories ec
// 				JOIN categories ca
// 					ON ca.id = ec.category_id
// 				WHERE ec.event_id = e.id
// 				AND EXISTS (
// 					SELECT 1
// 					FROM unnest($%d::text[]) AS search_category
// 					WHERE POSITION(
// 						LOWER(search_category)
// 						IN LOWER(ca.name)
// 					) > 0
// 				)
// 			) = cardinality($%d::text[])
// 		`, argPosition, argPosition))

// 		args = append(args, categories)
// 		argPosition++
// 	}

// 	if len(conditions) > 0 {
// 		query += " WHERE " + strings.Join(conditions, " AND ")
// 	}

// 	query += `
// 		ORDER BY attendee DESC
// 	`

// 	query += fmt.Sprintf(`
// 		LIMIT %d
// 		OFFSET $%d
// 	`, limit, argPosition)

// 	args = append(args, offset)

// 	rows, err := er.db.Query(ctx, query, args...)
// 	if err != nil {
// 		return nil, pkg.ParseError(err)
// 	}

// 	var events []model.DetailEvenstMDL
// 	var eventIDs []int32

// 	for rows.Next() {
// 		var event model.DetailEvenstMDL

// 		err := rows.Scan(
// 			&event.ID,
// 			&event.Community,
// 			// &event.OrganizerID,
// 			&event.Attendee,
// 			&event.Title,
// 			&event.Location,
// 			&event.Description,
// 			&event.Image,
// 			&event.Capacity,
// 			&event.StartTime,
// 			&event.EndTime,
// 			&event.CreatedAt,
// 			&event.UpdatedAt,
// 		)

// 		if err != nil {
// 			rows.Close()
// 			return nil, pkg.ParseError(err)
// 		}

// 		event.Speakers = []SpeakerMDL{}
// 		event.Categories = []string{}

// 		events = append(events, event)
// 		eventIDs = append(eventIDs, event.ID)
// 	}

// 	if err := rows.Err(); err != nil {
// 		rows.Close()
// 		return nil, pkg.ParseError(err)
// 	}

// 	rows.Close()

// 	if len(eventIDs) == 0 {
// 		return []EventsMDL{}, nil
// 	}

// 	// Get speakers
// 	querySpeakers := `
// 		SELECT
// 			se.event_id,
// 			s.name,
// 			s.job,
// 			s.office
// 		FROM speakers_events se
// 		JOIN speakers s
// 			ON s.id = se.speaker_id
// 		WHERE se.event_id = ANY($1)
// 		ORDER BY se.event_id, s.name
// 	`

// 	speakerRows, err := er.db.Query(ctx, querySpeakers, eventIDs)
// 	if err != nil {
// 		return nil, pkg.ParseError(err)
// 	}

// 	speakersMap := make(map[int32][]model.SpeakerMDL)

// 	for speakerRows.Next() {
// 		var (
// 			eventID int32
// 			speaker model.SpeakerMDL
// 		)

// 		err := speakerRows.Scan(
// 			&eventID,
// 			&speaker.Name,
// 			&speaker.Job,
// 			&speaker.Office,
// 		)

// 		if err != nil {
// 			speakerRows.Close()
// 			return nil, pkg.ParseError(err)
// 		}

// 		speakersMap[eventID] = append(
// 			speakersMap[eventID],
// 			speaker,
// 		)
// 	}

// 	if err := speakerRows.Err(); err != nil {
// 		speakerRows.Close()
// 		return nil, pkg.ParseError(err)
// 	}

// 	speakerRows.Close()

// 	// Get categories
// 	queryCategories := `
// 		SELECT
// 			ec.event_id,
// 			c.name
// 		FROM events_categories ec
// 		JOIN categories c
// 			ON c.id = ec.category_id
// 		WHERE ec.event_id = ANY($1)
// 		ORDER BY ec.event_id, c.name
// 	`

// 	categoryRows, err := er.db.Query(ctx, queryCategories, eventIDs)
// 	if err != nil {
// 		return nil, pkg.ParseError(err)
// 	}

// 	categoriesMap := make(map[int32][]string)

// 	for categoryRows.Next() {
// 		var (
// 			eventID  int32
// 			category string
// 		)

// 		err := categoryRows.Scan(
// 			&eventID,
// 			&category,
// 		)

// 		if err != nil {
// 			categoryRows.Close()
// 			return nil, pkg.ParseError(err)
// 		}

// 		categoriesMap[eventID] = append(
// 			categoriesMap[eventID],
// 			category,
// 		)
// 	}

// 	if err := categoryRows.Err(); err != nil {
// 		categoryRows.Close()
// 		return nil, pkg.ParseError(err)
// 	}

// 	categoryRows.Close()

// 	for i := range events {
// 		if speakers, ok := speakersMap[events[i].ID]; ok {
// 			events[i].Speakers = speakers
// 		}

// 		if categories, ok := categoriesMap[events[i].ID]; ok {
// 			events[i].Categories = categories
// 		}
// 	}

// 	return events, nil
// }
