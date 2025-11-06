package middleware

import (
	"context"
	"strconv"

	problemDTO "problum/internal/problem/service/dto"

	"github.com/gofiber/fiber/v3"
)

type ProblemService interface {
	Get(context.Context, int) (*problemDTO.Problem, error)
}

func Problem(problemSvc ProblemService, lessonSvc LessonService) fiber.Handler {
	return func(c fiber.Ctx) error {
		problemID, err := strconv.Atoi(c.Params("problemID"))
		if err != nil {
			return c.SendStatus(fiber.StatusForbidden)
		}

		courseIDString := c.Locals("course_id")
		if courseIDString == nil {
			return c.SendStatus(fiber.StatusForbidden)
		}
		courseID, _ := courseIDString.(int)

		problem, err := problemSvc.Get(c.Context(), problemID)
		if err != nil {
			return c.SendStatus(fiber.StatusForbidden)
		}

		lesson, err := lessonSvc.Get(c.Context(), problem.LessonID)
		if err != nil {
			return c.SendStatus(fiber.StatusForbidden)
		}

		if lesson.CourseID != courseID {
			return c.SendStatus(fiber.StatusForbidden)
		}

		return c.Next()
	}
}
