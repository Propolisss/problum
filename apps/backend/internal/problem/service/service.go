package service

import (
	"context"
	"fmt"

	attemptDTO "problum/internal/attempt/service/dto"
	"problum/internal/model"
	"problum/internal/problem/service/dto"
	templateDTO "problum/internal/template/service/dto"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
)

type Repository interface {
	Get(context.Context, int) (*model.Problem, error)
	ListByLessonID(context.Context, int) ([]*model.Problem, error)
}

type AttemptService interface {
	Submit(context.Context, *attemptDTO.Attempt) (int, error)
}

type TemplateService interface {
	GetByProblemIDAndLanguage(context.Context, int, string) (*templateDTO.Template, error)
}

type Service struct {
	repo        Repository
	js          jetstream.JetStream
	attemptSvc  AttemptService
	templateSvc TemplateService
}

func New(repo Repository, js jetstream.JetStream, attemptSvc AttemptService, templateSvc TemplateService) *Service {
	return &Service{
		repo:        repo,
		js:          js,
		attemptSvc:  attemptSvc,
		templateSvc: templateSvc,
	}
}

func (s *Service) Get(ctx context.Context, id int) (*dto.Problem, error) {
	problem, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get problem by id")
		return nil, fmt.Errorf("failed to get problem by id: %w", err)
	}

	return dto.ToDTO(problem, nil), nil
}

func (s *Service) GetWithTemplate(ctx context.Context, id int, language string) (*dto.Problem, error) {
	problem, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get problem by id")
		return nil, fmt.Errorf("failed to get problem by id: %w", err)
	}

	template, err := s.templateSvc.GetByProblemIDAndLanguage(ctx, id, language)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get template")
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	return dto.ToDTO(problem, template), nil
}

func (s *Service) ListByLessonID(ctx context.Context, lessonID int) ([]*dto.Problem, error) {
	problems, err := s.repo.ListByLessonID(ctx, lessonID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get problems by lesson id")
		return nil, fmt.Errorf("failed to get problems by lesson id: %w", err)
	}

	return dto.ToDTOList(problems), nil
}

func (s *Service) Submit(ctx context.Context, submit *dto.ProblemSubmit) (int, error) {
	id, err := s.attemptSvc.Submit(ctx, &attemptDTO.Attempt{
		ProblemID: submit.ProblemID,
		UserID:    submit.UserID,
		Language:  submit.Language,
		Code:      submit.Code,
		Status:    "pending",
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to submit")
		return 0, fmt.Errorf("failed to submit: %w", err)
	}
	submit.ID = id

	payload, err := sonic.Marshal(submit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal payload")
		return 0, fmt.Errorf("failed to create or update stream: %w", err)
	}

	if _, err := s.js.PublishAsync("ATTEMPTS.new", payload); err != nil {
		log.Error().Err(err).Msg("Failed to publish")
		return 0, fmt.Errorf("failed to publish: %w", err)
	}

	return id, nil
}
