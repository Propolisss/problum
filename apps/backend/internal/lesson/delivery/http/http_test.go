package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"problum/internal/api"
	"problum/internal/lesson/delivery/http/mocks"
	"problum/internal/lesson/service/dto"
	"problum/internal/utils"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/mock/gomock"
)

func TestHandlerGet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockService(ctrl)
	svc.EXPECT().Get(gomock.Any(), 3).Return(&dto.Lesson{ID: 3, CourseID: 1, Name: "Intro"}, nil)

	app := fiber.New()
	handler := New(nil, svc)
	app.Get("/lessons/:lessonID", handler.Get)

	resp, err := app.Test(utils.NewTestRequest(http.MethodGet, "/lessons/3", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("got status %d, want %d", resp.StatusCode, fiber.StatusOK)
	}

	var body api.LessonGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.ID != 3 || body.Name != "Intro" {
		t.Fatalf("unexpected response: %+v", body)
	}
}
