package service

import (
	"backend/EventHub/internal/dto"
	"backend/EventHub/internal/middleware"
	"backend/EventHub/internal/repo"
	"backend/EventHub/pkg"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var (
	ErrFieldEmpty  = errors.New("all fields are required")
	ErrFormatEmail = errors.New("wrong format email")
)

type AuthService struct {
	ar    *repo.AuthRepo
	conf  *pkg.ConfigHash
	redis *redis.Client
	db    *pgxpool.Pool
}

func NewAuthService(AuthRepo *repo.AuthRepo, redis *redis.Client, db *pgxpool.Pool) *AuthService {
	return &AuthService{
		ar:    AuthRepo,
		conf:  pkg.NewConfigHash(),
		redis: redis,
		db:    db,
	}
}

func (as *AuthService) RegisterAuthService(ctx context.Context, body dto.ReqRegister) (*dto.ResRegister, error) {
	if len(body.Email) == 0 || len(body.Name) == 0 || len(body.Password) <= 8 {
		return nil, ErrFieldEmpty
	}

	hashPassword, err := as.conf.HashPassword(body.Password)
	if err != nil {
		return nil, err
	}

	body.Password = hashPassword

	user, err := as.ar.RegisterAuthRepo(ctx, as.db, body)

	if err != nil {
		return nil, err
	}

	return &dto.ResRegister{
		ID:   user.ID,
		Name: user.Name,
		Role: user.Role,
	}, err
}

func (as *AuthService) LoginAuthService(ctx context.Context, body dto.ReqLogin) (*dto.ResLogin, error) {

	if len(body.Email) == 0 || len(body.Password) <= 8 {
		return nil, ErrFieldEmpty
	}

	user, err := as.ar.GetUserByEmail(ctx, as.db, body.Email)

	if err != nil {
		return nil, err
	}

	pwd, err := as.conf.Compare(body.Password, user.Password)
	if err != nil {
		return nil, err
	}

	body.Password = pwd

	body.Role = user.Role

	jwt := pkg.NewJWTClaims(user.ID, user.Role, 1440)

	token, err := jwt.GeneretToken()

	if err != nil {
		return nil, err
	}

	body.Token = token

	res, err := as.ar.Login(ctx, as.db, body)
	if err != nil {
		return nil, err
	}
	jwt = pkg.NewJWTClaims(res.ID, user.Role, 600)

	token, err = jwt.GeneretToken()
	return &dto.ResLogin{
		Token:          token,
		Name:           user.Name,
		Role:           user.Role,
		DarkPreference: user.DarkPreference,
		Image:          user.Image,
		Email:          user.Email,
	}, nil
}

func (as *AuthService) LogoutAuthService(ctx context.Context, payload *middleware.Payload) error {
	if payload.ID <= 0 {
		return errors.New("invalid auth id")
	}

	key := fmt.Sprintf("Token:%d", payload.ID)
	err := as.ar.Logout(ctx, as.db, payload.ID)
	if err != nil {
		return err
	}

	ttl := time.Duration(payload.ExpiresIn) * time.Second
	if ttl <= 0 {
		ttl = 1 * time.Second
	}

	err = as.redis.Set(ctx, key, "1", ttl).Err()
	if err != nil {
		log.Printf("[Redis.Set] Critical Error setting key %s (TTL: %v): %v\n", key, ttl, err)
	}

	return nil
}

func (as *AuthService) ChangePassword(ctx context.Context, body dto.ReqChangePassword) error {
	tx, err := as.db.Begin(ctx)
	if err != nil {
		return pkg.ParseError(err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Println(err.Error())
		}
	}()

	user, err := as.ar.GetUserByEmail(ctx, as.db, body.Email)
	if err != nil {
		return err
	}

	if len(body.Password) < 8 {
		return errors.New("password minimal 8 character")
	}

	pwd, err := as.conf.HashPassword(body.Password)
	if err != nil {
		return err
	}

	body.ID = user.ID
	body.Password = pwd

	err = as.ar.ChangePassword(ctx, tx, body.Password, body.ID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return pkg.ParseError(err)
	}

	return nil
}
