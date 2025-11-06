package middleware

import (
	"context"
	"strconv"

	lessonDTO "problum/internal/lesson/service/dto"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

type LessonService interface {
	Get(ctx context.Context, id int) (*lessonDTO.Lesson, error)
}

func Lesson(lessonSvc LessonService) fiber.Handler {
	return func(c fiber.Ctx) error {
		lessonID, err := strconv.Atoi(c.Params("lessonID"))
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse lesson id")
			return c.SendStatus(fiber.StatusForbidden)
		}

		courseIDString := c.Locals("course_id")
		if courseIDString == nil {
			log.Error().Msg("Failed to get course id")
			return c.SendStatus(fiber.StatusForbidden)
		}
		courseID, _ := courseIDString.(int)

		lesson, err := lessonSvc.Get(c.Context(), lessonID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get lesson")
			return c.SendStatus(fiber.StatusForbidden)
		}

		if lesson.CourseID != courseID {
			log.Error().Msg("Mismatch course id")
			return c.SendStatus(fiber.StatusForbidden)
		}

		return c.Next()
	}
}
