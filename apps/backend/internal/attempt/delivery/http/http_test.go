package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"problum/internal/api"
	"problum/internal/attempt/delivery/http/mocks"
	"problum/internal/attempt/service/dto"
	"problum/internal/utils"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/mock/gomock"
)

func TestHandlerListByProblemID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockService(ctrl)
	svc.EXPECT().ListByProblemID(gomock.Any(), 2, 3).Return([]*dto.Attempt{{ID: 1, UserID: 2, ProblemID: 3}}, nil)

	app := fiber.New()
	app.Use(func(c fiber.Ctx) error {
		c.Locals("user_id", 2)
		return c.Next()
	})
	handler := New(nil, svc)
	app.Get("/problems/:problemID/attempts", handler.ListByProblemID)

	resp, err := app.Test(utils.NewTestRequest(http.MethodGet, "/problems/3/attempts", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("got status %d, want %d", resp.StatusCode, fiber.StatusOK)
	}

	var body api.AttemptListResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Attempts) != 1 || body.Attempts[0].ProblemID != 3 {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestHandlerListByUserID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockService(ctrl)
	svc.EXPECT().ListByUserID(gomock.Any(), 2).Return([]*dto.Attempt{{ID: 1, UserID: 2}}, nil)

	app := fiber.New()
	app.Use(func(c fiber.Ctx) error {
		c.Locals("user_id", 2)
		return c.Next()
	})
	handler := New(nil, svc)
	app.Get("/attempts", handler.ListByUserID)

	resp, err := app.Test(utils.NewTestRequest(http.MethodGet, "/attempts", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("got status %d, want %d", resp.StatusCode, fiber.StatusOK)
	}

	var body api.AttemptListResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Attempts) != 1 || body.Attempts[0].UserID != 2 {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestHandlerGet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockService(ctrl)
	svc.EXPECT().Get(gomock.Any(), 4).Return(&dto.Attempt{ID: 4, UserID: 2, ProblemID: 3}, nil)

	app := fiber.New()
	handler := New(nil, svc)
	app.Get("/attempts/:attemptID", handler.Get)

	resp, err := app.Test(utils.NewTestRequest(http.MethodGet, "/attempts/4", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("got status %d, want %d", resp.StatusCode, fiber.StatusOK)
	}

	var body api.AttemptGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.ID != 4 || body.ProblemID != 3 {
		t.Fatalf("unexpected response: %+v", body)
	}
}
