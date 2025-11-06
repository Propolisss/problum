package service

import (
	"context"
	"fmt"

	"problum/internal/attempt/service/dto"
	"problum/internal/model"

	"github.com/rs/zerolog/log"
)

type Repository interface {
	ListByProblemID(context.Context, int, int) ([]*model.Attempt, error)
	ListByUserID(context.Context, int) ([]*model.Attempt, error)
	Submit(context.Context, *model.Attempt) (int, error)
	Update(context.Context, *model.Attempt) error
	Get(context.Context, int) (*model.Attempt, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ListByProblemID(ctx context.Context, userID, problemID int) ([]*dto.Attempt, error) {
	attempts, err := s.repo.ListByProblemID(ctx, userID, problemID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get attempts by problem id")
		return nil, fmt.Errorf("failed to get attempts by problem id: %w", err)
	}

	return dto.ToDTOList(attempts), nil
}

func (s *Service) ListByUserID(ctx context.Context, userID int) ([]*dto.Attempt, error) {
	attempts, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get attempts by user id")
		return nil, fmt.Errorf("failed to get attempts by user id: %w", err)
	}

	return dto.ToDTOList(attempts), nil
}

func (s *Service) Submit(ctx context.Context, attempt *dto.Attempt) (int, error) {
	id, err := s.repo.Submit(ctx, dto.ToModel(attempt))
	if err != nil {
		log.Error().Err(err).Msg("Failed to submit")
		return 0, fmt.Errorf("failed to submit: %w", err)
	}

	return id, nil
}

func (s *Service) Update(ctx context.Context, attempt *dto.Attempt) error {
	if err := s.repo.Update(ctx, dto.ToModel(attempt)); err != nil {
		log.Error().Err(err).Msg("Failed to update attempt")
		return fmt.Errorf("failed to update attempt: %w", err)
	}

	return nil
}

func (s *Service) Get(ctx context.Context, attemptID int) (*dto.Attempt, error) {
	attempt, err := s.repo.Get(ctx, attemptID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get attempt")
		return nil, fmt.Errorf("failed to get attempts: %w", err)
	}

	return dto.ToDTO(attempt), nil
}
