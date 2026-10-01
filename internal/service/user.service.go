package service

import (
	"backend/EventHub/internal/dto"
	"backend/EventHub/internal/model"
	"backend/EventHub/internal/repo"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	ar    *repo.AuthRepo
	as    *AuthService
	ur    *repo.UserRepo
	redis *redis.Client
}

func NewUserService(db *pgxpool.Pool, ur *repo.UserRepo, ar *repo.AuthRepo, as *AuthService, redis *redis.Client) *UserService {
	return &UserService{
		ar:    ar,
		as:    as,
		ur:    ur,
		redis: redis,
	}
}

func (us *UserService) NewPassword(ctx context.Context, body dto.ReqNewPassword) error {
	auth, err := us.ar.GetAuthUser(ctx, body.ID)
	if err != nil {
		return err
	}

	user, err := us.ar.GetUserById(ctx, auth.User_id)

	if err != nil {
		return err
	}

	pwd, err := us.as.conf.Compare(body.Old_password, user.Password)
	if err != nil {
		return err
	}

	pwd, err = us.as.conf.HashPassword(body.New_password)

	if err != nil {
		return err
	}

	err = us.ar.ChangePassword(ctx, pwd, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) MyProfile(ctx context.Context, ID int32) (*dto.ResMyProfle, error) {
	auth, err := us.ar.GetAuthUser(ctx, ID)
	if err != nil {
		return nil, err
	}
	var user_redis model.UserMDL
	key := fmt.Sprintf("user_id:%d", auth.User_id)
	res, err := us.redis.Get(ctx, key).Result()
	if err == nil {
		println(auth.User_id)
		if err := json.Unmarshal([]byte(res), &user_redis); err != nil {
			log.Printf("Error 222 : %s", err.Error())
		} else {
			log.Println("Data from redis", user_redis)
			return &dto.ResMyProfle{
				ID:             user_redis.ID,
				Email:          user_redis.Email,
				Name:           user_redis.Name,
				Role:           user_redis.Role,
				DarkPreference: user_redis.DarkPreference,
				Address:        user_redis.Address,
				Job:            user_redis.Job,
				Office:         user_redis.Office,
				Image:          user_redis.Image,
				Description:    user_redis.Description,
			}, nil
		}
	} else {
		log.Printf("Redis get %s service myProfile : %s", key, err)

	}

	user, err := us.ar.GetUserById(ctx, auth.User_id)

	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(user)
	if err != nil {
		log.Printf("Marshal service myProfile : %s", err)
	}

	us.redis.Set(ctx, key, string(marshal), 0)

	return &dto.ResMyProfle{
		ID:             user.ID,
		Email:          user.Email,
		Name:           user.Name,
		Role:           user.Role,
		DarkPreference: user.DarkPreference,
		Address:        user.Address,
		Job:            user.Job,
		Office:         user.Office,
		Image:          user.Image,
		Description:    user.Description,
	}, nil
}

func (us *UserService) SetProfile(ctx context.Context, ID int32) (*dto.ResMyProfle, error) {
	auth, err := us.ar.GetAuthUser(ctx, ID)
	if err != nil {
		return nil, err
	}

	user, err := us.ar.GetUserById(ctx, auth.User_id)

	if err != nil {
		return nil, err
	}
	return &dto.ResMyProfle{
		ID:             user.ID,
		Email:          user.Email,
		Name:           user.Name,
		Role:           user.Role,
		DarkPreference: user.DarkPreference,
		Address:        user.Address,
		Job:            user.Job,
		Office:         user.Office,
		Image:          user.Image,
		Description:    user.Description,
	}, nil
}

func (us *UserService) SetProfileService(ctx context.Context, ID int32, body dto.DtoSetProfile) (*dto.DtoSetProfile, error) {
	if body.Name != "" && len(body.Name) < 2 {
		return nil, errors.New("name must be at least 2 characters long")
	}

	auth, err := us.ar.GetAuthUser(ctx, ID)
	if err != nil {
		return nil, err
	}

	user, err := us.ar.GetUserById(ctx, auth.User_id)

	profile, err := us.ur.SetProfiles(ctx, user.ID, body)
	if err != nil {
		return nil, err
	}

	return &dto.DtoSetProfile{
		Name:        profile.Name,
		Address:     profile.Address,
		Job:         profile.Job,
		Office:      profile.Office,
		Image:       profile.Image,
		Description: profile.Description,
	}, nil
}
