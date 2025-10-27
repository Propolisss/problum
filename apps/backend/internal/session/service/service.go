package service

import (
	"context"

	"problum/internal/model"
)

type Repository interface {
	GetByRefreshHash(context.Context, string) (*model.UserSession, error)
	Create(context.Context, *model.UserSession) (*model.UserSession, error)
	Update(context.Context, *model.UserSession) (*model.UserSession, error)
	GetByPreviousRefreshHash(context.Context, string) (*model.UserSession, error)
	LogoutAll(context.Context, int) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, session *model.UserSession) (*model.UserSession, error) {
	return s.repo.Create(ctx, session)
}

func (s *Service) GetByRefreshHash(ctx context.Context, refresh string) (*model.UserSession, error) {
	return s.repo.GetByRefreshHash(ctx, refresh)
}

func (s *Service) GetByPreviousRefreshHash(ctx context.Context, refresh string) (*model.UserSession, error) {
	return s.repo.GetByPreviousRefreshHash(ctx, refresh)
}

func (s *Service) Update(ctx context.Context, session *model.UserSession) (*model.UserSession, error) {
	return s.repo.Update(ctx, session)
}

func (s *Service) LogoutAll(ctx context.Context, id int) error {
	return s.repo.LogoutAll(ctx, id)
}
