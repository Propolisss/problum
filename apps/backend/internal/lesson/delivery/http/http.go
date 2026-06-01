package http

import (
	"context"
	"fmt"
	"strconv"

	"problum/internal/config"
	"problum/internal/lesson/service/dto"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

//go:generate go run go.uber.org/mock/mockgen -source=http.go -destination=mocks/mock_service.go -package=mocks
type Service interface {
	Get(context.Context, int) (*dto.Lesson, error)
	// пока не нужны для апишки
	// List(ctx context.Context) ([]*dto.Lesson, error)
	// ListByCourseID(ctx context.Context, id int) ([]*dto.Lesson, error)
}

type Handler struct {
	cfg *config.Config
	svc Service
}

func New(cfg *config.Config, svc Service) *Handler {
	return &Handler{
		cfg: cfg,
		svc: svc,
	}
}

// Get returns a lesson.
// @Summary Get lesson
// @Description Returns a lesson by ID with problems.
// @Tags lessons
// @Produce json
// @Param courseID path int true "Course ID"
// @Param lessonID path int true "Lesson ID"
// @Success 200 {object} api.LessonGetResponse
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Security BearerAuth
// @Router /courses/{courseID}/lessons/{lessonID} [get]
func (h *Handler) Get(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("lessonID"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse lesson id in query param")
		return fmt.Errorf("failed to parse lesson id in query param")
	}

	lesson, err := h.svc.Get(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get lesson")
		return fmt.Errorf("failed to get lesson")
	}

	return c.JSON(dto.ToAPI(lesson))
}
