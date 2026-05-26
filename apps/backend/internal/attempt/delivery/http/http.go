package http

import (
	"context"
	"strconv"

	"problum/internal/api"
	"problum/internal/attempt/service/dto"
	"problum/internal/config"

	"github.com/gofiber/fiber/v3"
)

//go:generate go run go.uber.org/mock/mockgen -source=http.go -destination=mocks/mock_service.go -package=mocks
type Service interface {
	ListByProblemID(context.Context, int, int) ([]*dto.Attempt, error)
	ListByUserID(context.Context, int) ([]*dto.Attempt, error)
	Get(context.Context, int) (*dto.Attempt, error)
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

// ListByProblemID returns attempts for a problem.
// @Summary List problem attempts
// @Description Returns the authenticated user's attempts for a problem.
// @Tags attempts
// @Produce json
// @Param courseID path int true "Course ID"
// @Param problemID path int true "Problem ID"
// @Success 200 {object} api.AttemptListResponse
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Security BearerAuth
// @Router /courses/{courseID}/problems/{problemID}/attempts [get]
func (h *Handler) ListByProblemID(c fiber.Ctx) error {
	problemID, err := strconv.Atoi(c.Params("problemID"))
	if err != nil {
		return c.SendStatus(fiber.StatusForbidden)
	}

	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}

	attempts, err := h.svc.ListByProblemID(c.Context(), userID, problemID)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(api.AttemptListResponse{
		Attempts: dto.ToAPIList(attempts),
	})
}

// ListByUserID returns attempts for the current user.
// @Summary List user attempts
// @Description Returns all attempts for the authenticated user.
// @Tags attempts
// @Produce json
// @Success 200 {object} api.AttemptListResponse
// @Failure 401 {string} string
// @Failure 404 {string} string
// @Security BearerAuth
// @Router /attempts [get]
func (h *Handler) ListByUserID(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}

	attempts, err := h.svc.ListByUserID(c.Context(), userID)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(api.AttemptListResponse{
		Attempts: dto.ToAPIList(attempts),
	})
}

// Get returns an attempt.
// @Summary Get attempt
// @Description Returns an attempt by ID.
// @Tags attempts
// @Produce json
// @Param attemptID path int true "Attempt ID"
// @Success 200 {object} api.AttemptGetResponse
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Failure 404 {string} string
// @Security BearerAuth
// @Router /attempts/{attemptID} [get]
func (h *Handler) Get(c fiber.Ctx) error {
	attemptID, err := strconv.Atoi(c.Params("attemptID"))
	if err != nil {
		return c.SendStatus(fiber.StatusForbidden)
	}

	attempt, err := h.svc.Get(c.Context(), attemptID)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(dto.ToAPI(attempt))
}
