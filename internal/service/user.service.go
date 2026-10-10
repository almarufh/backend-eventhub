package service

import (
	"backend/EventHub/internal/dto"
	msgerr "backend/EventHub/internal/message"
	"backend/EventHub/internal/model"
	"backend/EventHub/internal/repo"
	"backend/EventHub/pkg"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	db    *pgxpool.Pool
	ar    *repo.AuthRepo
	er    *repo.EventsRepo
	as    *AuthService
	ur    *repo.UserRepo
	redis *redis.Client
}

func NewUserService(db *pgxpool.Pool, ur *repo.UserRepo, ar *repo.AuthRepo, as *AuthService, redis *redis.Client, er *repo.EventsRepo) *UserService {
	return &UserService{
		db:    db,
		ar:    ar,
		er:    er,
		as:    as,
		ur:    ur,
		redis: redis,
	}
}

func (us *UserService) NewPassword(ctx context.Context, body dto.ReqNewPassword) error {
	tx, err := us.db.Begin(ctx)
	if err != nil {
		return pkg.ParseError(err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Println(err.Error())
		}
	}()
	auth, err := us.ar.GetAuthUser(ctx, us.db, body.ID)
	if err != nil {
		return err
	}

	user, err := us.ar.GetUserById(ctx, us.db, auth.User_id)

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

	err = us.ar.ChangePassword(ctx, tx, pwd, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) MyProfile(ctx context.Context, ID int32) (*dto.ResMyProfle, error) {
	auth, err := us.ar.GetAuthUser(ctx, us.db, ID)
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

	user, err := us.ar.GetUserById(ctx, us.db, auth.User_id)

	if err != nil {
		return nil, err
	}

	if user.Image != nil && *user.Image != "" {
		host := os.Getenv("HOST_STATIC")
		imageUrl := fmt.Sprintf("http://%s/public/images/profile/%s", host, *user.Image)
		user.Image = &imageUrl
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

func (us *UserService) SetProfileService(ctx context.Context, ID int32, body dto.ReqSetProfile) (*dto.ResSetProfile, error) {

	auth, err := us.ar.GetAuthUser(ctx, us.db, ID)
	if err != nil {
		return nil, msgerr.AuthNotFound
	}

	var filename string

	if body.Image != nil {
		file := body.Image

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
			return nil, msgerr.InvalidImageExt
		}

		if file.Size > 2*1024*1024 {
			return nil, msgerr.ImageTooLarge
		}

		filename = fmt.Sprintf("%d.jpg", auth.User_id)

		dirPath := path.Join("public", "photo_profiles")
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			log.Printf("SetProfileService.path.Join : %s", err.Error())
			return nil, msgerr.WriteFile
		}

		targetFilePath := path.Join(dirPath, filename)

		src, err := file.Open()
		if err != nil {
			return nil, msgerr.OpenFile
		}
		defer src.Close()

		data, err := io.ReadAll(src)
		if err != nil {
			log.Printf("SetProfileService.io.ReadAll : %s", err.Error())
			return nil, msgerr.ReadFile
		}

		if err := os.WriteFile(targetFilePath, data, 0644); err != nil {
			log.Printf("SetProfileService.os.Write : %s", err.Error())
			return nil, msgerr.WriteFile
		}
	}

	user, err := us.ar.GetUserById(ctx, us.db, auth.User_id)
	if err != nil || user == nil {
		return nil, msgerr.UserNotFound
	}

	key := fmt.Sprintf("user_id:%d", auth.User_id)

	profile, err := us.ur.SetProfiles(ctx, user.ID, filename, body)
	if err != nil {
		return nil, msgerr.UpdateProfile
	}

	if profile.Image != nil && *profile.Image != "" {
		host := os.Getenv("HOST_STATIC")
		imageUrl := fmt.Sprintf("http://%s/public/images/profile/%s", host, *profile.Image)
		profile.Image = &imageUrl
	}

	var data *model.UserMDL
	data = &model.UserMDL{
		ID:             user.ID,
		Email:          user.Email,
		Name:           *profile.Name,
		Role:           user.Role,
		DarkPreference: user.DarkPreference,
		Address:        profile.Address,
		Job:            profile.Job,
		Office:         profile.Office,
		Image:          profile.Image,
		Description:    profile.Description,
	}

	fmt.Println(data.ID)

	marshal, err := json.Marshal(data)
	if err != nil {
		log.Printf("Marshal service myProfile : %s", err)
	}

	us.redis.Set(ctx, key, string(marshal), 0)

	return &dto.ResSetProfile{
		Name:        profile.Name,
		Address:     profile.Address,
		Job:         profile.Job,
		Office:      profile.Office,
		Image:       profile.Image,
		Description: profile.Description,
	}, nil
}

func (us *UserService) JoinedEvents(ctx context.Context, ID int32) (dto.ResJoinedEvents, error) {
	auth, err := us.ar.GetAuthUser(ctx, us.db, ID)
	if err != nil {
		return dto.ResJoinedEvents{}, msgerr.AuthNotFound
	}

	eventsJoined, err := us.ur.JoinedEvents(ctx, auth.User_id)
	if err != nil {
		return dto.ResJoinedEvents{}, fmt.Errorf("[UserService.JoinedEvents] %w", err)
	}

	res := make([]dto.ResDetailEvent, 0, len(eventsJoined))

	for _, e := range eventsJoined {
		organizer, err := us.ar.GetUserById(ctx, us.db, e.Organizer.ID)
		if err != nil {
			return dto.ResJoinedEvents{}, err
		}

		speakers, err := us.er.GetSpeakersByEventID(ctx, e.ID)
		if err != nil {
			return dto.ResJoinedEvents{}, err
		}

		speakersDTO := make([]dto.Speaker, 0, len(speakers))
		for _, s := range speakers {
			speakersDTO = append(speakersDTO, dto.Speaker{
				Name:   s.Name,
				Job:    s.Job,
				Office: s.Office,
			})
		}

		categories, err := us.er.GetCategoryEvents(ctx, e.Organizer.ID)
		if err != nil {
			return dto.ResJoinedEvents{}, err
		}
		res = append(res, dto.ResDetailEvent{
			ID:        e.ID,
			Community: e.Community,
			Organizer: dto.Organizer{
				Name:   organizer.Name,
				Job:    *organizer.Job,
				Office: *organizer.Office,
				Image:  *organizer.Image,
			},
			Speakers:    speakersDTO,
			Categorys:   categories,
			Attendee:    e.Attendee,
			Title:       e.Title,
			Location:    e.Location,
			Description: e.Description,
			Image:       e.Image,
			Capacity:    e.Capacity,
			StartTime:   e.StartTime,
			EndTime:     e.EndTime,
		})
	}

	total := uint32(len(res))

	if total < 1 {
		return dto.ResJoinedEvents{}, msgerr.JoinedEventsNotFound
	}

	return dto.ResJoinedEvents{
		Total: total,
		Data:  res,
	}, nil
}

func (us *UserService) SavedEvents(ctx context.Context, ID int32) (dto.ResSavedEvents, error) {
	auth, err := us.ar.GetAuthUser(ctx, us.db, ID)
	if err != nil {
		return dto.ResSavedEvents{}, msgerr.AuthNotFound
	}

	eventsJoined, err := us.ur.SavedEvents(ctx, auth.User_id)
	if err != nil {
		return dto.ResSavedEvents{}, fmt.Errorf("[UserService.JoinedEvents] %w", err)
	}

	res := make([]dto.ResDetailEvent, 0, len(eventsJoined))

	for _, e := range eventsJoined {
		organizer, err := us.ar.GetUserById(ctx, us.db, e.Organizer.ID)
		if err != nil {
			return dto.ResSavedEvents{}, err
		}

		speakers, err := us.er.GetSpeakersByEventID(ctx, e.ID)
		if err != nil {
			return dto.ResSavedEvents{}, err
		}

		speakersDTO := make([]dto.Speaker, 0, len(speakers))
		for _, s := range speakers {
			speakersDTO = append(speakersDTO, dto.Speaker{
				Name:   s.Name,
				Job:    s.Job,
				Office: s.Office,
			})
		}

		categories, err := us.er.GetCategoryEvents(ctx, e.Organizer.ID)
		if err != nil {
			return dto.ResSavedEvents{}, err
		}
		res = append(res, dto.ResDetailEvent{
			ID:        e.ID,
			Community: e.Community,
			Organizer: dto.Organizer{
				Name:   organizer.Name,
				Job:    *organizer.Job,
				Office: *organizer.Office,
				Image:  *organizer.Image,
			},
			Speakers:    speakersDTO,
			Categorys:   categories,
			Attendee:    e.Attendee,
			Title:       e.Title,
			Location:    e.Location,
			Description: e.Description,
			Image:       e.Image,
			Capacity:    e.Capacity,
			StartTime:   e.StartTime,
			EndTime:     e.EndTime,
		})
	}

	total := uint32(len(res))

	if total < 1 {
		return dto.ResSavedEvents{}, msgerr.SavedEventsUserNotFound
	}

	return dto.ResSavedEvents{
		Total: total,
		Data:  res,
	}, nil
}

func (us *UserService) JoinedCommunities(ctx context.Context, ID int32) (dto.ResJoinedCommunities, error) {
	auth, err := us.ar.GetAuthUser(ctx, us.db, ID)
	if err != nil {
		return dto.ResJoinedCommunities{}, msgerr.AuthNotFound
	}

	communitiesJoined, err := us.ur.JoinedCommunities(ctx, auth.User_id)
	if err != nil {
		return dto.ResJoinedCommunities{}, fmt.Errorf("[UserService.JoinedEvents] %w", err)
	}

	res := make([]dto.ResDetailCommunitiy, 0, len(communitiesJoined))

	for _, e := range communitiesJoined {
		res = append(res, dto.ResDetailCommunitiy{
			ID:          e.ID,
			Title:       e.Title,
			Description: e.Description,
			Image:       e.Image,
			Status:      e.Status,
			Categories:  e.Categories,
			Members:     e.Members,
		})
	}

	total := uint32(len(res))

	if total < 1 {
		return dto.ResJoinedCommunities{}, msgerr.JoinedEventsNotFound
	}

	return dto.ResJoinedCommunities{
		Total: total,
		Data:  res,
	}, nil
}
