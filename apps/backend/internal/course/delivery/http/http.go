package http

import (
	"context"
	"strconv"

	"problum/internal/api"
	"problum/internal/config"
	"problum/internal/course/service/dto"

	"github.com/gofiber/fiber/v3"
)

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

func (h *Handler) List(c fiber.Ctx) error {
	resp, err := h.svc.List(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(api.CourseListResponse{
		Courses: dto.ListToAPI(resp),
	})
}

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
