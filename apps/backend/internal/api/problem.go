package api

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

type ProblemGetResponse struct {
	ID          int                 `json:"id"`
	LessonID    int                 `json:"lesson_id"`
	Name        string              `json:"name"`
	Statement   string              `json:"statement"`
	Difficulty  string              `json:"difficulty"`
	TimeLimit   time.Duration       `json:"time_limit"`
	MemoryLimit int64               `json:"memory_limit"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Template    TemplateGetResponse `json:"template,omitempty"`
}

type ProblemSubmitRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type ProblemSubmitResponse struct {
	AttemptID int `json:"attempt_id"`
}

type ProblemAPI interface {
	// надо ли???
	// List(fiber.Ctx) error
	Get(fiber.Ctx) error
	Submit(fiber.Ctx) error
}
