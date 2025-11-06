package api

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

type LessonGetResponse struct {
	ID          int       `json:"id"`
	CourseID    int       `json:"course_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Position    int       `json:"position"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Problems []ProblemGetResponse `json:"problems,omitempty"`
}

type LessonAPI interface {
	Get(fiber.Ctx) error
}
