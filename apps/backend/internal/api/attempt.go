package api

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

type AttemptGetResponse struct {
	ID           int           `json:"id"`
	UserID       int           `json:"user_id"`
	ProblemID    int           `json:"problem_id"`
	Duration     time.Duration `json:"duration" swaggertype:"integer"`
	MemoryUsage  int64         `json:"memory_usage"`
	Language     string        `json:"language"`
	Code         string        `json:"code"`
	Status       string        `json:"status"`
	ErrorMessage *string       `json:"error_message"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type AttemptListResponse struct {
	Attempts []AttemptGetResponse `json:"attempts"`
}

//go:generate go run go.uber.org/mock/mockgen -destination=mocks/mock_service.go -package=mocks problum/internal/api AttemptAPI,AuthAPI,CourseAPI,EnrollmentAPI,LessonAPI,ProblemAPI
type AttemptAPI interface {
	ListByProblemID(fiber.Ctx) error
	ListByUserID(fiber.Ctx) error
}
