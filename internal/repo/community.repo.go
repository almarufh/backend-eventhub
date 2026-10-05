package repo

import (
	"backend/EventHub/internal/model"
	"backend/EventHub/pkg"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommuntyRepo struct {
	db *pgxpool.Pool
}

func NewCommunityRepo(db *pgxpool.Pool) *CommuntyRepo {
	return &CommuntyRepo{
		db: db,
	}
}

func (ur *CommuntyRepo) ToggleJoinCommunityRepo(ctx context.Context, ID int32, param int32) (model.JoinCommunityMDL, error) {
	query := `
		WITH deleted AS (
		    DELETE FROM communities_users 
		    WHERE user_id = $1 AND community_id = $2
		    RETURNING community_id
		),
		inserted AS (
		    INSERT INTO communities_users (user_id, community_id)
		    SELECT $1, $2
		    WHERE NOT EXISTS (SELECT 1 FROM deleted)
		    RETURNING community_id
		)
		SELECT 
		    c.id, 
		    c.title,
		    CASE 
		        WHEN EXISTS (SELECT 1 FROM deleted) 
				THEN 'Unjoined'
		        ELSE 'Joined'
		    END AS status
		FROM communities c
		WHERE c.id = $2
    `
	args := []any{ID, param}
	var communities model.JoinCommunityMDL
	err := ur.db.QueryRow(
		ctx,
		query,
		args...,
	).Scan(
		&communities.ID,
		&communities.Title,
		&communities.Status,
	)

	if err != nil {
		parsedErr := pkg.ParseError(err)
		if errors.Is(parsedErr, pkg.ErrNotFound) {
			return model.JoinCommunityMDL{}, errors.New("communities not found")
		}
		return model.JoinCommunityMDL{}, parsedErr
	}

	return communities, nil
}

func (ur *CommuntyRepo) GetDetailCommunity(ctx context.Context, ID int32) (*model.DetailCommunitieMDL, error) {
	query := `
		SELECT
		    c.id,
		    c.title,
		    c.description,
		    c.image,
		    c.status,
		    COALESCE(ARRAY_AGG(DISTINCT cg.name)) AS "categories",
		    (
		        SELECT COUNT(user_id) 
		        FROM communities_users 
		        WHERE community_id = c.id
		    ) AS "members"
		FROM communities c
		LEFT JOIN communities_categories cc ON c.id = cc.community_id
		LEFT JOIN categories cg ON cg.id = cc.category_id
		WHERE c.id = $1
		GROUP BY c.id;
    `
	var communities model.DetailCommunitieMDL
	args := []any{ID}
	if err := ur.db.QueryRow(ctx, query, args...).Scan(
		&communities.ID,
		&communities.Title,
		&communities.Description,
		&communities.Image,
		&communities.Status,
		&communities.Categories,
		&communities.Members,
	); err != nil {
		parsedErr := pkg.ParseError(err)
		if errors.Is(parsedErr, pkg.ErrNotFound) {
			return nil, errors.New("communities not found")
		}
		return nil, parsedErr
	}
	return &communities, nil
}

func (ur *CommuntyRepo) GetMembersCommunity(ctx context.Context, ID int32) (*[]model.MembersCommunityMDL, error) {
	query := `
		SELECT
			p.name,
			p.job,
			p.office,
			p.image
		FROM communities_users cu
		JOIN profiles p ON p.user_id = cu.user_id
		WHERE cu.community_id = $1
    `
	var membersCommunities []model.MembersCommunityMDL
	args := []any{ID}
	rows, err := ur.db.Query(ctx, query, args...)
	if err != nil {
		parsedErr := pkg.ParseError(err)
		if errors.Is(parsedErr, pkg.ErrNotFound) {
			return nil, errors.New("communities not found")
		}
		return nil, parsedErr
	}

	for rows.Next() {
		var users model.MembersCommunityMDL
		if err := rows.Scan(
			&users.Name,
			&users.Job,
			&users.Office,
			&users.Image,
		); err != nil {
			return nil, err
		}
		membersCommunities = append(membersCommunities, users)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &membersCommunities, nil
}

func (er *EventsRepo) GetCategoryCommunity(ctx context.Context, eventID int32) ([]string, error) {
	query := `
        SELECT
			c.name
		FROM categories c
		JOIN communities_categories cc ON c.id = cc.category_id
		WHERE cc.community_id = $1
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
