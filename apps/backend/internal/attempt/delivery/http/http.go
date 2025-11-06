package http

import (
	"context"
	"strconv"

	"problum/internal/api"
	"problum/internal/attempt/service/dto"
	"problum/internal/config"

	"github.com/gofiber/fiber/v3"
)

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
