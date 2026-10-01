package repo

import (
	"context"
	"errors"

	"backend/EventHub/internal/dto"
	"backend/EventHub/internal/model"
	"backend/EventHub/pkg"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrRegisterFailed     = errors.New("failed to register user")
)

type AuthRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (u *AuthRepo) RegisterAuthRepo(ctx context.Context, body dto.ReqRegister) (*model.RegisterMDL, error) {
	query := `
        WITH new_user AS (
            INSERT INTO users (
                email, 
                password
            )
            VALUES (
                $1, 
                $2
            )
            RETURNING id
        )
        INSERT INTO profiles (user_id, name, role)
        SELECT 
            id, 
            $3, 
            $4
        FROM new_user
        RETURNING 
            user_id, 
            name,
            role
    `
	role := body.Role
	if role == "" {
		role = "attendee"
	}

	args := []any{body.Email, body.Password, body.Name, role}

	var profile model.RegisterMDL
	err := u.db.QueryRow(ctx, query, args...).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Role,
	)

	if err != nil {
		return nil, pkg.ParseError(err)
	}

	return &profile, nil
}

func (u *AuthRepo) Login(ctx context.Context, body dto.ReqLogin) (*model.LoginMDL, error) {
	query := `
		WITH new_user_login AS (
		    SELECT u.id
		    FROM users u
		    WHERE u.email = $1 AND u.password = $2
		)
		INSERT INTO auth_users (
		    user_id,
		    token,
		    device
		)
		SELECT
		    id,
		    $3,
		    $4
		FROM new_user_login
		RETURNING id, user_id, token
    `
	args := []any{body.Email, body.Password, body.Token, body.Device}

	var authUser model.LoginMDL
	err := u.db.QueryRow(ctx, query, args...).Scan(
		&authUser.ID,
		&authUser.User_id,
		&authUser.Token,
	)

	if err != nil {
		parsedErr := pkg.ParseError(err)
		if errors.Is(parsedErr, pkg.ErrNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, parsedErr
	}

	return &authUser, nil
}

func (u *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*model.UserMDL, error) {
	query := `
		SELECT 
		    u.id, 
		    u.email, 
		    u.password,
		    p.name,
		    p.role,
		    p.dark_preference,
		    p.job,
		    p.office,
		    p.address,
		    p.image,
		    p.description
		FROM users u
		JOIN profiles p ON u.id = p.user_id
		WHERE u.email = $1;
	`
	var user model.UserMDL
	err := u.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.DarkPreference,
		&user.Address,
		&user.Job,
		&user.Office,
		&user.Image,
		&user.Description,
	)

	if err != nil {
		parsedErr := pkg.ParseError(err)
		if errors.Is(parsedErr, pkg.ErrNotFound) {
			return nil, errors.New("email dosn't exist")
		}
		return nil, parsedErr
	}
	return &user, nil
}

func (u *AuthRepo) GetUserById(ctx context.Context, id int32) (*model.UserMDL, error) {
	query := `
		SELECT 
		    u.id, 
		    u.email, 
		    u.password,
		    p.name,
		    p.role,
		    p.dark_preference,
		    p.job,
		    p.office,
		    p.address,
		    p.image,
		    p.description
		FROM users u
		JOIN profiles p ON u.id = p.user_id
		WHERE u.id = $1;
	`
	var user model.UserMDL
	err := u.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.DarkPreference,
		&user.Address,
		&user.Job,
		&user.Office,
		&user.Image,
		&user.Description,
	)

	if err != nil {
		parsedErr := pkg.ParseError(err)
		if errors.Is(parsedErr, pkg.ErrNotFound) {
			return nil, errors.New("email dosn't exist")
		}
		return nil, parsedErr
	}
	return &user, nil
}

func (u *AuthRepo) GetAuthUser(ctx context.Context, id int32) (*model.AuthUser, error) {
	query := `
        SELECT user_id, is_active, token
        FROM auth_users
        WHERE id = $1
    `

	var authUser model.AuthUser
	err := u.db.QueryRow(ctx, query, id).Scan(
		&authUser.User_id,
		&authUser.IsActive,
		&authUser.Token,
	)

	if err != nil {
		parsedErr := pkg.ParseError(err)
		if errors.Is(parsedErr, pkg.ErrNotFound) {
			return nil, errors.New("auth session not found")
		}
		return nil, parsedErr
	}

	return &authUser, nil
}

func (u *AuthRepo) Logout(ctx context.Context, id int32) error {
	query := `
		UPDATE auth_users
		SET is_active = FALSE,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND is_active = TRUE
	`

	tag, err := u.db.Exec(ctx, query, id)
	if err != nil {
		return pkg.ParseError(err)
	}

	if tag.RowsAffected() == 0 {
		return errors.New("auth session not found or already logged out")
	}

	return nil
}

func (u *AuthRepo) ChangePassword(ctx context.Context, password string, ID int32) error {
	query := `
        UPDATE users
        SET password = $1,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $2
    `
	args := []any{password, ID}
	tag, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return pkg.ParseError(err)
	}

	if tag.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}
