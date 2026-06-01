package http

import (
	"context"
	"fmt"

	"problum/internal/config"

	"problum/internal/user/service/dto"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

//go:generate go run go.uber.org/mock/mockgen -source=http.go -destination=mocks/mock_service.go -package=mocks
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

// Get returns the current user.
// @Summary Get profile
// @Description Returns the authenticated user's profile.
// @Tags profile
// @Produce json
// @Success 200 {object} api.UserGetResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Security BearerAuth
// @Router /profile [get]
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
