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
	GetLanguagesByProblemID(context.Context, int) ([]string, error)
}

type Service struct {
	repo        Repository
	js          jetstream.JetStream
	attemptSvc  AttemptService
	templateSvc TemplateService
}

type options struct {
	language      string
	withLanguage  bool
	withTemplate  bool
	withLanguages bool
}

type Option func(*options)

func WithTemplate() Option {
	return func(opts *options) {
		opts.withTemplate = true
	}
}

func WithLanguage(language string) Option {
	return func(opts *options) {
		opts.language = language
		opts.withLanguage = true
	}
}

func WithLanguages() Option {
	return func(opts *options) {
		opts.withLanguages = true
	}
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

	return dto.ToDTO(problem, nil, nil), nil
}

func (s *Service) GetWithOptions(ctx context.Context, id int, opts ...Option) (*dto.Problem, error) {
	problem, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get problem by id")
		return nil, fmt.Errorf("failed to get problem by id: %w", err)
	}

	options := &options{
		withTemplate:  false,
		withLanguages: false,
		language:      "",
	}
	for _, opt := range opts {
		opt(options)
	}

	if options.withLanguage && options.language == "" {
		log.Warn().Msg("Empty language")
		return nil, fmt.Errorf("empty language")
	}

	var template *templateDTO.Template
	if options.withTemplate {
		template, err = s.templateSvc.GetByProblemIDAndLanguage(ctx, id, options.language)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get template")
			return nil, fmt.Errorf("failed to get template: %w", err)
		}
	}

	var languages []string
	if options.withLanguages {
		languages, err = s.templateSvc.GetLanguagesByProblemID(ctx, id)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get languages")
			return nil, fmt.Errorf("failed to get languages: %w", err)
		}
	}

	return dto.ToDTO(problem, template, languages), nil
}

func (s *Service) ListByLessonID(ctx context.Context, lessonID int) ([]*dto.Problem, error) {
	problems, err := s.repo.ListByLessonID(ctx, lessonID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get problems by lesson id")
		return nil, fmt.Errorf("failed to get problems by lesson id: %w", err)
	}

	problemsDTO := dto.ToDTOList(problems)

	for _, problem := range problemsDTO {
		languages, err := s.templateSvc.GetLanguagesByProblemID(ctx, problem.ID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get languages")
			continue
		}

		problem.Languages = languages
	}

	return problemsDTO, nil
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
