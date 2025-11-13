package http

import (
	"context"
	"fmt"

	"problum/internal/config"

	"problum/internal/user/service/dto"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Get(context.Context, int) (*dto.User, error)
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

func (h *Handler) Get(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, err := h.svc.Get(c.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user")
		return fmt.Errorf("failed to get use")
	}

	return c.JSON(dto.ToAPI(user))
}
