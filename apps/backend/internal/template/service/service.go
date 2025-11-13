package service

import (
	"context"
	"fmt"

	"problum/internal/model"
	"problum/internal/template/service/dto"

	"github.com/rs/zerolog/log"
)

type Repository interface {
	GetByProblemIDAndLanguage(context.Context, int, string) (*model.Template, error)
	GetLanguagesByProblemID(context.Context, int) ([]string, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetByProblemIDAndLanguage(
	ctx context.Context,
	problemID int,
	language string,
) (*dto.Template, error) {
	template, err := s.repo.GetByProblemIDAndLanguage(ctx, problemID, language)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get template")
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	return dto.ToDTO(template), nil
}

func (s *Service) GetLanguagesByProblemID(ctx context.Context, problemID int) ([]string, error) {
	languages, err := s.repo.GetLanguagesByProblemID(ctx, problemID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get languages")
		return nil, fmt.Errorf("failed to get languages: %w", err)
	}

	return languages, nil
}
