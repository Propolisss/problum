package http

import (
	"context"
	"fmt"

	"problum/internal/api"
	"problum/internal/config"
	"problum/internal/model"
	"problum/internal/redis"

	"github.com/gofiber/fiber/v3"
)

//go:generate go run go.uber.org/mock/mockgen -source=http.go -destination=mocks/mock_service.go -package=mocks
type Service interface {
	Enroll(context.Context, int, int) error
}

type Handler struct {
	cfg *config.Config
	svc Service
	rdb *redis.Redis
}

func New(cfg *config.Config, svc Service) *Handler {
	return &Handler{
		cfg: cfg,
		svc: svc,
	}
}

// Enroll enrolls the current user.
// @Summary Enroll course
// @Description Enrolls the authenticated user into a course.
// @Tags enrollments
// @Accept json
// @Param request body api.EnrollRequest true "Enroll request"
// @Success 200 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Security BearerAuth
// @Router /enrollments [post]
func (h *Handler) Enroll(c fiber.Ctx) error {
	enrollReq := &api.EnrollRequest{}
	if err := c.Bind().JSON(enrollReq); err != nil {
		return err
	}

	us, ok := c.Locals("user_session").(*model.UserSession)
	if !ok {
		return fmt.Errorf("failed to cast user session")
	}

	if err := h.svc.Enroll(c.Context(), enrollReq.CourseID, us.UserID); err != nil {
		return fmt.Errorf("failed to enroll course")
	}

	return nil
}
