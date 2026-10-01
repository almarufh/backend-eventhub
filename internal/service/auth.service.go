package service

import (
	"backend/EventHub/internal/dto"
	"backend/EventHub/internal/repo"
	"backend/EventHub/pkg"
	"context"
	"errors"
	"log"
)

var (
	ErrFieldEmpty  = errors.New("all fields are required")
	ErrFormatEmail = errors.New("wrong format email")
)

type AuthService struct {
	ur   *repo.AuthRepo
	conf *pkg.ConfigHash
}

func NewAuthService(AuthRepo *repo.AuthRepo) *AuthService {
	return &AuthService{
		ur:   AuthRepo,
		conf: pkg.NewConfigHash(),
	}
}

func (us *AuthService) RegisterAuthService(ctx context.Context, body dto.ReqRegister) (*dto.ResRegister, error) {
	if len(body.Email) == 0 || len(body.Name) == 0 || len(body.Password) <= 8 {
		return nil, ErrFieldEmpty
	}

	hashPassword, err := us.conf.HashPassword(body.Password)
	if err != nil {
		return nil, err
	}

	body.Password = hashPassword

	user, err := us.ur.RegisterAuthRepo(ctx, body)

	if err != nil {
		return nil, err
	}

	return &dto.ResRegister{
		ID:   user.ID,
		Name: user.Name,
		Role: user.Role,
	}, err
}

func (us *AuthService) LoginAuthService(ctx context.Context, body dto.ReqLogin) (*dto.ResLogin, error) {

	if len(body.Email) == 0 || len(body.Password) <= 8 {
		return nil, ErrFieldEmpty
	}

	user, err := us.ur.GetUserByEmail(ctx, body.Email)

	if err != nil {
		return nil, err
	}

	pwd, err := us.conf.Compare(body.Password, user.Password)
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

	res, err := us.ur.Login(ctx, body)
	if err != nil {
		return nil, err
	}
	log.Println(res)
	jwt = pkg.NewJWTClaims(res.ID, user.Role, 600)

	token, err = jwt.GeneretToken()
	return &dto.ResLogin{
		Token: token,
	}, nil
}

func (us *AuthService) LogoutAuthService(ctx context.Context, authID int32) error {
	if authID <= 0 {
		return errors.New("invalid auth id")
	}

	err := us.ur.Logout(ctx, authID)
	if err != nil {
		return err
	}

	return nil
}

func (us *AuthService) ChangePassword(ctx context.Context, body dto.ReqChangePassword) error {
	user, err := us.ur.GetUserByEmail(ctx, body.Email)
	if err != nil {
		return err
	}

	if body.New_password != body.Confirm_password {
		return errors.New("wrong confirm password")
	}

	pwd, err := us.conf.HashPassword(body.New_password)
	if err != nil {
		return err
	}

	body.ID = user.ID
	body.New_password = pwd

	err = us.ur.ChangePassword(ctx, body.New_password, body.ID)
	if err != nil {
		return err
	}

	return nil
}
