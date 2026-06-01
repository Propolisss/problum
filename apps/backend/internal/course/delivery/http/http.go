package http

import (
	"context"
	"strconv"

	"problum/internal/api"
	"problum/internal/config"
	"problum/internal/course/service/dto"

	"github.com/gofiber/fiber/v3"
)

//go:generate go run go.uber.org/mock/mockgen -source=http.go -destination=mocks/mock_service.go -package=mocks
type Service interface {
	List(context.Context) ([]*dto.CourseDTO, error)
	Get(context.Context, int) (*dto.CourseDTO, error)
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

// List returns courses.
// @Summary List courses
// @Description Returns all courses for the authenticated user.
// @Tags courses
// @Produce json
// @Success 200 {object} api.CourseListResponse
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Security BearerAuth
// @Router /courses [get]
func (h *Handler) List(c fiber.Ctx) error {
	resp, err := h.svc.List(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(api.CourseListResponse{
		Courses: dto.ListToAPI(resp),
	})
}

// Get returns a course.
// @Summary Get course
// @Description Returns a course by ID with lessons.
// @Tags courses
// @Produce json
// @Param courseID path int true "Course ID"
// @Success 200 {object} api.CourseGetResponse
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 500 {string} string
// @Security BearerAuth
// @Router /courses/{courseID} [get]
func (h *Handler) Get(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("courseID"))
	if err != nil {
		return err
	}

	resp, err := h.svc.Get(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(dto.ToAPI(resp))
}
