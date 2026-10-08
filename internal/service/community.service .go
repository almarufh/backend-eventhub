package service

import (
	"backend/EventHub/internal/dto"
	"backend/EventHub/internal/middleware"
	"backend/EventHub/internal/repo"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommunitieService struct {
	ar *repo.AuthRepo
	cr *repo.CommuntyRepo
	db *pgxpool.Pool
}

func NewCommunitieService(db *pgxpool.Pool, cr *repo.CommuntyRepo, ar *repo.AuthRepo) *CommunitieService {
	return &CommunitieService{
		ar: ar,
		cr: cr,
		db: db,
	}
}

func (cs *CommunitieService) ToggleJoinCommunity(ctx context.Context, param int32, payload *middleware.Payload) (*dto.ResJoinCommunity, error) {
	auth, err := cs.ar.GetAuthUser(ctx, cs.db, payload.ID)
	if err != nil {
		return nil, err
	}

	joined, err := cs.cr.ToggleJoinCommunityRepo(ctx, auth.User_id, param)
	if err != nil {
		return nil, err
	}

	return &dto.ResJoinCommunity{
		ID:     joined.ID,
		Title:  joined.Title,
		Status: joined.Status,
	}, nil
}

func (cs *CommunitieService) DetailCommunity(ctx context.Context, param int32) (*dto.ResDetailCommunitiy, error) {
	communitie, err := cs.cr.GetDetailCommunity(ctx, param)
	if err != nil {
		return nil, err
	}

	return &dto.ResDetailCommunitiy{
		ID:          communitie.ID,
		Title:       communitie.Title,
		Description: communitie.Description,
		Image:       communitie.Image,
		Status:      communitie.Status,
		Categories:  communitie.Categories,
		Members:     communitie.Members,
	}, err
}

func (cs *CommunitieService) GetMembersCommunity(ctx context.Context, param int32) ([]dto.ResMembersCommunity, error) {
	members, err := cs.cr.GetMembersCommunity(ctx, param)
	if err != nil {
		return nil, err
	}

	var membersCommunities []dto.ResMembersCommunity
	for _, user := range *members {
		membersCommunities = append(membersCommunities, dto.ResMembersCommunity{
			Name:   user.Name,
			Job:    user.Job,
			Office: user.Office,
			Image:  user.Image,
		})
	}

	return membersCommunities, err
}
