package repo

import (
	"backend/EventHub/internal/dto"
	"backend/EventHub/internal/model"
	"backend/EventHub/pkg"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (ur *UserRepo) SetProfiles(ctx context.Context, ID int32, imageUrl string, body dto.ReqSetProfile) (model.ProfileMDL, error) {
	query := `
        UPDATE profiles
        SET 
            name = COALESCE(NULLIF($1, ''), name),
            address = COALESCE($2, address),
            job = COALESCE($3, job),
            office = COALESCE($4, office),
            image = COALESCE($5, image),
            description = COALESCE($6, description),
            updated_at = CURRENT_TIMESTAMP
        WHERE user_id = $7
        RETURNING 
            name, 
            address, 
            job, 
            office, 
            image, 
            description,
			created_at
    `
	args := []any{
		body.Name,
		body.Address,
		body.Job,
		body.Office,
		imageUrl,
		body.Description,
		ID,
	}
	var profile model.ProfileMDL
	err := ur.db.QueryRow(
		ctx,
		query,
		args...,
	).Scan(
		&profile.Name,
		&profile.Address,
		&profile.Job,
		&profile.Office,
		&profile.Image,
		&profile.Description,
		&profile.CreatedAt,
	)

	if err != nil {
		parsedErr := pkg.ParseError(err)
		if errors.Is(parsedErr, pkg.ErrNotFound) {
			return model.ProfileMDL{}, errors.New("profile not found")
		}
		return model.ProfileMDL{}, parsedErr
	}

	return profile, nil
}

func (ur *UserRepo) JoinedEvents(ctx context.Context, userID int32) ([]model.DetailEvenstMDL, error) {
	query := `
		SELECT 
			e.id,
			c.title AS community,
			e.organizer_id,
			COUNT(je_all.user_id) AS attendee,
			e.title,
			e.location,
			e.description,
			e.image,
			e.capacity,
			e.start_time,
			e.end_time,
			e.created_at,
			e.updated_at
		FROM joined_events_users je_user
		JOIN events e ON e.id = je_user.event_id
		JOIN communities c ON c.id = e.community_id
		LEFT JOIN joined_events_users je_all ON je_all.event_id = e.id
		WHERE je_user.user_id = $1
		GROUP BY e.id, c.title
		ORDER BY e.start_time ASC;
	`

	rows, err := ur.db.Query(ctx, query, userID)
	if err != nil {
		log.Printf("[UserRepo.SetProfiles] Query : %s", err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return []model.DetailEvenstMDL{}, nil
		}
		return nil, fmt.Errorf("gagal mengeksekusi query joined events: %w", err)
	}
	defer rows.Close()

	events := make([]model.DetailEvenstMDL, 0)

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
		)
		if err != nil {
			return nil, fmt.Errorf("gagal memindai data event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("terjadi kesalahan saat membaca baris event: %w", err)
	}

	return events, nil
}

func (ur *UserRepo) SavedEvents(ctx context.Context, userID int32) ([]model.DetailEvenstMDL, error) {
	query := `
		SELECT 
			e.id,
			c.title AS community,
			e.organizer_id,
			COUNT(se_all.user_id) AS attendee,
			e.title,
			e.location,
			e.description,
			e.image,
			e.capacity,
			e.start_time,
			e.end_time,
			e.created_at,
			e.updated_at
		FROM saved_events_users se_user
		JOIN events e ON e.id = se_user.event_id
		JOIN communities c ON c.id = e.community_id
		LEFT JOIN saved_events_users se_all ON se_all.event_id = e.id
		WHERE se_user.user_id = $1
		GROUP BY e.id, c.title
		ORDER BY e.start_time ASC
	`

	rows, err := ur.db.Query(ctx, query, userID)
	if err != nil {
		log.Printf("[UserRepo.SetProfiles] Query : %s", err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return []model.DetailEvenstMDL{}, nil
		}
		return nil, fmt.Errorf("gagal mengeksekusi query saved events: %w", err)
	}
	defer rows.Close()

	events := make([]model.DetailEvenstMDL, 0)

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
		)
		if err != nil {
			return nil, fmt.Errorf("gagal memindai data event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("terjadi kesalahan saat membaca baris event: %w", err)
	}

	return events, nil
}

func (ur *UserRepo) JoinedCommunities(ctx context.Context, userID int32) ([]model.DetailCommunitieMDL, error) {
	query := `
		SELECT 
		    c.id,
		    c.title,
		    c.description,
		    c.image,
		    c.status,
		    COALESCE(ARRAY_AGG(DISTINCT cg.name) FILTER (WHERE cg.name IS NOT NULL), '{}') AS "categories",
		    (
		        SELECT COUNT(*) 
		        FROM communities_users cu 
		        WHERE cu.community_id = c.id
		    ) AS members
		FROM communities_users cu_user
		JOIN communities c ON c.id = cu_user.community_id
		LEFT JOIN communities_categories cc ON c.id = cc.community_id
		LEFT JOIN categories cg ON cg.id = cc.category_id
		WHERE cu_user.user_id = $1
		GROUP BY c.id;
	`

	rows, err := ur.db.Query(ctx, query, userID)
	if err != nil {
		log.Printf("[UserRepo.SetProfiles] Query : %s", err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("gagal mengeksekusi query joined community: %w", err)
	}
	defer rows.Close()

	communities := make([]model.DetailCommunitieMDL, 0)

	for rows.Next() {
		var community model.DetailCommunitieMDL

		err := rows.Scan(
			&community.ID,
			&community.Title,
			&community.Description,
			&community.Image,
			&community.Status,
			&community.Categories,
			&community.Members,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal memindai data community: %w", err)
		}
		communities = append(communities, community)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("terjadi kesalahan saat membaca baris event: %w", err)
	}

	return communities, nil
}
