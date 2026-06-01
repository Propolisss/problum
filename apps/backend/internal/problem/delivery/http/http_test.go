package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"problum/internal/api"
	"problum/internal/problem/delivery/http/mocks"
	"problum/internal/problem/service/dto"
	"problum/internal/utils"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/mock/gomock"
)

func TestHandlerGet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockService(ctrl)
	svc.EXPECT().GetWithOptions(gomock.Any(), 3, gomock.Any(), gomock.Any()).Return(&dto.Problem{ID: 3, LessonID: 1, Name: "Two Sum"}, nil)

	app := fiber.New()
	handler := New(nil, svc)
	app.Get("/problems/:problemID", handler.Get)

	resp, err := app.Test(utils.NewTestRequest(http.MethodGet, "/problems/3?language=go", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("got status %d, want %d", resp.StatusCode, fiber.StatusOK)
	}

	var body api.ProblemGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.ID != 3 || body.Name != "Two Sum" {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestHandlerSubmit(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockService(ctrl)
	svc.EXPECT().Submit(gomock.Any(), gomock.Any()).DoAndReturn(func(_ any, submit *dto.ProblemSubmit) (int, error) {
		if submit.ProblemID != 3 || submit.UserID != 2 || submit.Language != "go" || submit.Code != "package main" {
			t.Fatalf("unexpected submit: %+v", submit)
		}
		return 7, nil
	})

	app := fiber.New()
	app.Use(func(c fiber.Ctx) error {
		c.Locals("user_id", 2)
		return c.Next()
	})
	handler := New(nil, svc)
	app.Post("/problems/:problemID/submit", handler.Submit)

	resp, err := app.Test(utils.NewTestRequest(http.MethodPost, "/problems/3/submit", strings.NewReader(`{"language":"go","code":"package main"}`)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("got status %d, want %d", resp.StatusCode, fiber.StatusOK)
	}

	var body api.ProblemSubmitResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.AttemptID != 7 {
		t.Fatalf("unexpected response: %+v", body)
	}
}
