package service

import (
	"context"
	"fmt"

	"problum/internal/course/service/dto"
	"problum/internal/model"

	enrollmentDTO "problum/internal/enrollment/service/dto"
	lessonDTO "problum/internal/lesson/service/dto"

	"github.com/rs/zerolog/log"
)

type Repository interface {
	List(context.Context) ([]*model.Course, error)
	Get(context.Context, int) (*model.Course, error)
}

type LessonService interface {
	ListByCourseID(ctx context.Context, id int) ([]*lessonDTO.Lesson, error)
}

type EnrollmentService interface {
	Get(context.Context, int, int) (*enrollmentDTO.Enrollment, error)
	GetListByUserID(context.Context, int) ([]*enrollmentDTO.Enrollment, error)
}

type Service struct {
	repo          Repository
	lessonSvc     LessonService
	enrollmentSvc EnrollmentService
}

func New(repo Repository, lessonSvc LessonService, enrollmentSvc EnrollmentService) *Service {
	return &Service{
		repo:          repo,
		lessonSvc:     lessonSvc,
		enrollmentSvc: enrollmentSvc,
	}
}

func (s *Service) List(ctx context.Context) ([]*dto.CourseDTO, error) {
	courses, err := s.repo.List(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get courses list")
		return nil, fmt.Errorf("failed to get courses list: %w", err)
	}

	var enrollments []*enrollmentDTO.Enrollment

	userID, ok := ctx.Value("user_id").(int)
	if ok {
		if enrolls, err := s.enrollmentSvc.GetListByUserID(ctx, userID); err == nil {
			enrollments = enrolls
		}
	} else {
		log.Error().Err(err).Msg("Failed to cast user_id")
	}

	return dto.ToDTOList(courses, enrollments), nil
}

func (s *Service) Get(ctx context.Context, id int) (*dto.CourseDTO, error) {
	course, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Error().Int("course_id", id).Err(err).Msg("Failed to get course")
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	lessons, err := s.lessonSvc.ListByCourseID(ctx, id)
	if err != nil {
		log.Error().Int("course_id", id).Err(err).Msg("Failed to get lessons for the course")
		return nil, fmt.Errorf("failed to get lessons for the course: %w", err)
	}

	enrolled := false

	userID, ok := ctx.Value("user_id").(int)
	if ok {
		enrollment, err := s.enrollmentSvc.Get(ctx, id, userID)
		if err == nil && enrollment != nil && enrollment.CourseID == id && enrollment.UserID == userID {
			enrolled = true
		}
	}

	return dto.ToDTO(course, lessons, enrolled), nil
}
