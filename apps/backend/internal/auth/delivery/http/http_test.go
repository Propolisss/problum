package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"problum/internal/api"
	"problum/internal/auth/delivery/http/mocks"
	"problum/internal/auth/service/dto"
	"problum/internal/config"
	"problum/internal/utils"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/mock/gomock"
)

func TestHandlerLogin(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockService(ctrl)
	svc.EXPECT().Login(gomock.Any(), "user", "pass").Return(&dto.LoginDTO{
		AccessToken:  "access",
		RefreshToken: "refresh",
		ExpiresAt:    time.Minute,
	}, nil)

	app := fiber.New()
	handler := New(testConfig(), svc)
	app.Post("/login", handler.Login)

	resp, err := app.Test(utils.NewTestRequest(http.MethodPost, "/login", strings.NewReader(`{"login":"user","password":"pass"}`)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("got status %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
	if got := refreshCookie(resp); got != "refresh" {
		t.Fatalf("got cookie %q, want %q", got, "refresh")
	}

	var body api.LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.AccessToken != "access" || body.ExpiresAt != time.Minute {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestHandlerRegister(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockService(ctrl)
	svc.EXPECT().Register(gomock.Any(), "user", "pass", "pass").Return(&dto.RegisterDTO{
		AccessToken:  "access",
		RefreshToken: "refresh",
		ExpiresAt:    time.Minute,
	}, nil)

	app := fiber.New()
	handler := New(testConfig(), svc)
	app.Post("/register", handler.Register)

	resp, err := app.Test(utils.NewTestRequest(http.MethodPost, "/register", strings.NewReader(`{"login":"user","password":"pass","repeated_password":"pass"}`)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("got status %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
	if got := refreshCookie(resp); got != "refresh" {
		t.Fatalf("got cookie %q, want %q", got, "refresh")
	}
}

func TestHandlerRefresh(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockService(ctrl)
	svc.EXPECT().Refresh(gomock.Any(), "old-refresh").Return(&dto.RefreshDTO{
		AccessToken:  "new-access",
		RefreshToken: "new-refresh",
		ExpiresAt:    time.Minute,
	}, nil)

	app := fiber.New()
	handler := New(testConfig(), svc)
	app.Post("/refresh", handler.Refresh)

	req := utils.NewTestRequest(http.MethodPost, "/refresh", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "old-refresh"})

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("got status %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
	if got := refreshCookie(resp); got != "new-refresh" {
		t.Fatalf("got cookie %q, want %q", got, "new-refresh")
	}

	var body api.RefreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.AccessToken != "new-access" {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestHandlerLogout(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockService(ctrl)
	svc.EXPECT().Logout(gomock.Any(), "access", "refresh").Return(nil)

	app := fiber.New()
	handler := New(testConfig(), svc)
	app.Post("/logout", handler.Logout)

	req := utils.NewTestRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "refresh"})
	req.Header.Set("Authorization", "Bearer access")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("got status %d, want %d", resp.StatusCode, fiber.StatusNoContent)
	}
}

func testConfig() *config.Config {
	return &config.Config{Server: &config.Server{}}
}

func refreshCookie(resp *http.Response) string {
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "refresh_token" {
			return cookie.Value
		}
	}
	return ""
}
