package service

import (
	"context"
	"fmt"

	"problum/internal/model"
	"problum/internal/user/service/dto"

	"github.com/rs/zerolog/log"
)

type Repository interface {
	FindByLogin(context.Context, string) (*model.User, error)
	Create(context.Context, *model.User) (*model.User, error)
	Get(context.Context, int) (*model.User, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Get(ctx context.Context, userID int) (*dto.User, error) {
	user, err := s.repo.Get(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find user by user_id")
		return nil, fmt.Errorf("failed to find user by user_id: %w", err)
	}

	return dto.ToDTO(user), nil
}

func (s *Service) FindByLogin(ctx context.Context, login string) (*dto.User, error) {
	user, err := s.repo.FindByLogin(ctx, login)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find user by login")
		return nil, fmt.Errorf("failed to find user by login: %w", err)
	}

	return dto.ToDTO(user), nil
}

func (s *Service) Create(ctx context.Context, user *dto.User) (*dto.User, error) {
	u, err := s.repo.Create(ctx, dto.ToModel(user))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return dto.ToDTO(u), nil
}
