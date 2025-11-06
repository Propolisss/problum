package api

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

type CourseListResponse struct {
	Courses []CourseGetResponse `json:"courses"`
}

type CourseGetResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Enrolled    bool      `json:"enrolled"`

	Lessons []LessonGetResponse `json:"lessons,omitempty"`
}

type CourseAPI interface {
	List(fiber.Ctx) error
	Get(fiber.Ctx) error
	// Favorite(fiber.Ctx) error
}
