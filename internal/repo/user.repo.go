package repo

import (
	"backend/EventHub/internal/dto"
	"backend/EventHub/internal/model"
	"backend/EventHub/pkg"
	"context"
	"errors"

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

func (ur *UserRepo) SetProfiles(ctx context.Context, ID int32, body dto.DtoSetProfile) (*model.ProfileMDL, error) {
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
            description
    `
	args := []any{
		body.Name,
		body.Address,
		body.Job,
		body.Office,
		body.Image,
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
	)

	if err != nil {
		parsedErr := pkg.ParseError(err)
		if errors.Is(parsedErr, pkg.ErrNotFound) {
			return nil, errors.New("profile not found")
		}
		return nil, parsedErr
	}

	return &profile, nil
}
