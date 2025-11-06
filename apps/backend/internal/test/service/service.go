package service

import (
	"context"
	"fmt"

	"problum/internal/model"
	"problum/internal/test/service/dto"

	"github.com/rs/zerolog/log"
)

type Repository interface {
	GetByProblemID(ctx context.Context, problemID int) (*model.Test, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetByProblemID(ctx context.Context, problemID int) (*dto.Test, error) {
	test, err := s.repo.GetByProblemID(ctx, problemID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get test by problem id")
		return nil, fmt.Errorf("failed to get test by problem id")
	}

	return dto.ToDTO(test), nil
}
