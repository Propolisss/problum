package service

import (
	"context"

	"problum/internal/model"
)

type Repository interface {
	FindByLogin(context.Context, string) (*model.User, error)
	Create(context.Context, *model.User) (*model.User, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) FindByLogin(ctx context.Context, login string) (*model.User, error) {
	return s.repo.FindByLogin(ctx, login)
}

func (s *Service) Create(ctx context.Context, user *model.User) (*model.User, error) {
	return s.repo.Create(ctx, user)
}
