package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"problum/internal/api"
	"problum/internal/user/delivery/http/mocks"
	"problum/internal/user/service/dto"
	"problum/internal/utils"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/mock/gomock"
)

func TestHandlerGet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockService(ctrl)
	svc.EXPECT().Get(gomock.Any(), 2).Return(&dto.User{ID: 2, Login: "user"}, nil)

	app := fiber.New()
	app.Use(func(c fiber.Ctx) error {
		c.Locals("user_id", 2)
		return c.Next()
	})
	handler := New(nil, svc)
	app.Get("/me", handler.Get)

	resp, err := app.Test(utils.NewTestRequest(http.MethodGet, "/me", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("got status %d, want %d", resp.StatusCode, fiber.StatusOK)
	}

	var body api.UserGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.ID != 2 || body.Login != "user" {
		t.Fatalf("unexpected response: %+v", body)
	}
}
