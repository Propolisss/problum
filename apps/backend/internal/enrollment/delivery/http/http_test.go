package http

import (
	"net/http"
	"strings"
	"testing"

	"problum/internal/enrollment/delivery/http/mocks"
	"problum/internal/model"
	"problum/internal/utils"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/mock/gomock"
)

func TestHandlerEnroll(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockService(ctrl)
	svc.EXPECT().Enroll(gomock.Any(), 1, 2).Return(nil)

	app := fiber.New()
	app.Use(func(c fiber.Ctx) error {
		c.Locals("user_session", &model.UserSession{UserID: 2})
		return c.Next()
	})
	handler := New(nil, svc)
	app.Post("/enrollments", handler.Enroll)

	resp, err := app.Test(utils.NewTestRequest(http.MethodPost, "/enrollments", strings.NewReader(`{"course_id":1}`)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("got status %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
}
