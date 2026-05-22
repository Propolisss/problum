package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"problum/internal/api"
	"problum/internal/course/delivery/http/mocks"
	"problum/internal/course/service/dto"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/mock/gomock"
)

func TestHandlerList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(*mocks.MockService)
		wantStatus int
		wantLen    int
	}{
		{
			name: "success",
			setup: func(svc *mocks.MockService) {
				svc.EXPECT().List(gomock.Any()).Return([]*dto.CourseDTO{{ID: 1, Name: "Go"}}, nil)
			},
			wantStatus: fiber.StatusOK,
			wantLen:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			svc := mocks.NewMockService(ctrl)
			tt.setup(svc)

			app := fiber.New()
			handler := New(nil, svc)
			app.Get("/courses", handler.List)

			resp, err := app.Test(newRequest(http.MethodGet, "/courses", nil))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			var body api.CourseListResponse
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if len(body.Courses) != tt.wantLen {
				t.Fatalf("got %d courses, want %d", len(body.Courses), tt.wantLen)
			}
		})
	}
}

func TestHandlerGet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockService(ctrl)
	svc.EXPECT().Get(gomock.Any(), 3).Return(&dto.CourseDTO{ID: 3, Name: "Algorithms"}, nil)

	app := fiber.New()
	handler := New(nil, svc)
	app.Get("/courses/:courseID", handler.Get)

	resp, err := app.Test(newRequest(http.MethodGet, "/courses/3", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("got status %d, want %d", resp.StatusCode, fiber.StatusOK)
	}

	var body api.CourseGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.ID != 3 || body.Name != "Algorithms" {
		t.Fatalf("unexpected response: %+v", body)
	}
}
