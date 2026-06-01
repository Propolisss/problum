package http

import (
	"context"
	"fmt"
	"strings"
	"time"

	"problum/internal/api"
	"problum/internal/auth/service/dto"
	"problum/internal/config"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

//go:generate go run go.uber.org/mock/mockgen -source=http.go -destination=mocks/mock_service.go -package=mocks
type Service interface {
	Login(context.Context, string, string) (*dto.LoginDTO, error)
	Register(context.Context, string, string, string) (*dto.RegisterDTO, error)
	Refresh(context.Context, string) (*dto.RefreshDTO, error)
	Logout(context.Context, string, string) error
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

// Login authenticates a user.
// @Summary Login
// @Description Authenticates a user and returns an access token. The refresh token is also set as an HTTP-only cookie.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body api.LoginRequest true "Login request"
// @Success 200 {object} api.LoginResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /auth/login [post]
func (h *Handler) Login(c fiber.Ctx) error {
	loginReq := &api.LoginRequest{}
	if err := c.Bind().JSON(loginReq); err != nil {
		return err
	}

	resp, err := h.svc.Login(c.Context(), loginReq.Login, loginReq.Password)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		Expires:  time.Now().Add(14 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   h.cfg.IsProduction(),
		SameSite: "Lax",
	})

	return c.JSON(api.LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    resp.ExpiresAt,
	})
}

// Register creates a user account.
// @Summary Register
// @Description Creates a new user account, logs the user in, and sets the refresh token as an HTTP-only cookie.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body api.RegisterRequest true "Register request"
// @Success 200 {object} api.RegisterResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /auth/register [post]
func (h *Handler) Register(c fiber.Ctx) error {
	registerReq := &api.RegisterRequest{}
	if err := c.Bind().JSON(registerReq); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal request")
		return fmt.Errorf("failed to unmarshal request: %w", err)
	}

	resp, err := h.svc.Register(c.Context(), registerReq.Login, registerReq.Password, registerReq.RepeatedPassword)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register user")
		return fmt.Errorf("failed to register user: %w", err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		Expires:  time.Now().Add(14 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   h.cfg.IsProduction(),
		SameSite: "Lax",
	})

	return c.JSON(api.RegisterResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    resp.ExpiresAt,
	})
}

// Refresh rotates tokens.
// @Summary Refresh tokens
// @Description Refreshes the access token using the refresh_token cookie.
// @Tags auth
// @Produce json
// @Success 200 {object} api.RefreshResponse
// @Failure 401 {string} string
// @Router /auth/refresh [post]
func (h *Handler) Refresh(c fiber.Ctx) error {
	refresh := c.Cookies("refresh_token")
	if refresh == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	resp, err := h.svc.Refresh(c.Context(), refresh)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		Expires:  time.Now().Add(14 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   h.cfg.IsProduction(),
		SameSite: "Lax",
	})

	return c.JSON(api.RefreshResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    resp.ExpiresAt,
	})
}

// Logout revokes the current session.
// @Summary Logout
// @Description Revokes the current refresh token and clears the refresh_token cookie.
// @Tags auth
// @Success 204 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Security BearerAuth
// @Router /auth/logout [post]
func (h *Handler) Logout(c fiber.Ctx) error {
	refresh := c.Cookies("refresh_token")
	if refresh == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	access := c.Get("Authorization")
	if access == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	access = strings.TrimPrefix(access, "Bearer ")

	if err := h.svc.Logout(c.Context(), access, refresh); err != nil {
		log.Error().Err(err).Msg("Failed to logout")
		return fmt.Errorf("failed to logout")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-100 * time.Hour),
		HTTPOnly: true,
		Secure:   h.cfg.IsProduction(),
		SameSite: "Lax",
	})

	return c.SendStatus(fiber.StatusNoContent)
}
