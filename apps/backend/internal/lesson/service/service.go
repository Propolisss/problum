package service

import (
	"context"
	"fmt"
	"sync"

	"problum/internal/lesson/service/dto"
	"problum/internal/model"
	problemDTO "problum/internal/problem/service/dto"

	"github.com/rs/zerolog/log"
)

type Repository interface {
	Get(context.Context, int) (*model.Lesson, error)
	List(context.Context) ([]*model.Lesson, error)
	ListByCourseID(context.Context, int) ([]*model.Lesson, error)
}

type ProblemService interface {
	ListByLessonID(context.Context, int) ([]*problemDTO.Problem, error)
}

type Service struct {
	repo       Repository
	problemSvc ProblemService
}

func New(repo Repository, problemSvc ProblemService) *Service {
	return &Service{
		repo:       repo,
		problemSvc: problemSvc,
	}
}

func (s *Service) Get(ctx context.Context, id int) (*dto.Lesson, error) {
	lesson, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get lesson")
		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}

	problems, err := s.problemSvc.ListByLessonID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get problems for lesson")
		return nil, fmt.Errorf("failed to get problems for lesson: %w", err)
	}

	return dto.ToDTO(lesson, problems), nil
}

func (s *Service) List(ctx context.Context) ([]*dto.Lesson, error) {
	lessons, err := s.repo.List(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get lessons")
		return nil, fmt.Errorf("failed to get lessons: %w", err)
	}

	return dto.ToDTOList(lessons, nil), nil
}

func (s *Service) ListByCourseID(ctx context.Context, id int) ([]*dto.Lesson, error) {
	lessons, err := s.repo.ListByCourseID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get lessons by course id")
		return nil, fmt.Errorf("failed to get lessons by course id: %w", err)
	}

	mu := &sync.Mutex{}
	m := make(map[int][]*problemDTO.Problem)
	wg := &sync.WaitGroup{}

	wg.Add(len(lessons))
	for _, lesson := range lessons {
		go func() {
			defer wg.Done()

			problems, err := s.problemSvc.ListByLessonID(ctx, lesson.ID)
			if err != nil {
				return
			}

			mu.Lock()
			m[lesson.ID] = problems
			mu.Unlock()
		}()
	}
	wg.Wait()

	return dto.ToDTOList(lessons, m), nil
}
