package api

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"-"`
	ExpiresAt    time.Duration `json:"expires_at"`
}

type RegisterRequest struct {
	Login            string `json:"login"`
	Password         string `json:"password"`
	RepeatedPassword string `json:"repeated_password"`
}

type RegisterResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"-"`
	ExpiresAt    time.Duration `json:"expires_at"`
}

type RefreshResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"-"`
	ExpiresAt    time.Duration `json:"expires_at"`
}

type AuthAPI interface {
	Login(fiber.Ctx) error
	Register(fiber.Ctx) error
	Refresh(fiber.Ctx) error
	Logout(fiber.Ctx) error
}
