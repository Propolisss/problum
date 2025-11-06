package service

import (
	"context"
	"fmt"

	"problum/internal/enrollment/service/dto"
	"problum/internal/model"

	"github.com/rs/zerolog/log"
)

type Repository interface {
	Enroll(context.Context, int, int) error
	Get(ctx context.Context, courseID, userID int) (*model.Enrollment, error)
	GetListByUserID(context.Context, int) ([]*model.Enrollment, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Enroll(ctx context.Context, courseID, userID int) error {
	if enrollment, err := s.repo.Get(ctx, courseID, userID); err == nil && enrollment != nil {
		log.Error().Err(err).Int("course_id", courseID).Int("user_id", userID).Msg("Already enrolled")
		return fmt.Errorf("already enrolled: %w", err)
	}

	if err := s.repo.Enroll(ctx, courseID, userID); err != nil {
		log.Error().Err(err).Int("course_id", courseID).Int("user_id", userID).Msg("Failed to enroll course")
		return fmt.Errorf("failed to enroll course: %w", err)
	}

	return nil
}

func (s *Service) Get(ctx context.Context, courseID, userID int) (*dto.Enrollment, error) {
	enrollment, err := s.repo.Get(ctx, courseID, userID)
	if err != nil {
		log.Error().Err(err).Int("course_id", courseID).Int("user_id", userID).Msg("Failed to get enrollment")
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}

	return dto.ToDTO(enrollment), nil
}

func (s *Service) GetListByUserID(ctx context.Context, userID int) ([]*dto.Enrollment, error) {
	enrollments, err := s.repo.GetListByUserID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Int("user_id", userID).Msg("Failed to get enrollments by user id")
		return nil, fmt.Errorf("failed to get enrollments by user id: %w", err)
	}

	return dto.ToDTOList(enrollments), nil
}
