package middleware

import (
	"context"
	"strconv"

	"problum/internal/model"

	enrollmentDTO "problum/internal/enrollment/service/dto"

	"github.com/gofiber/fiber/v3"
)

type EnrollmentService interface {
	Get(ctx context.Context, courseID, userID int) (*enrollmentDTO.Enrollment, error)
}

func Course(courseSvc EnrollmentService) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("courseID"))
		if err != nil {
			return c.SendStatus(fiber.StatusForbidden)
		}

		usString := c.Locals("user_session")
		if usString == nil {
			return c.SendStatus(fiber.StatusForbidden)
		}

		us, ok := usString.(*model.UserSession)
		if !ok {
			return c.SendStatus(fiber.StatusForbidden)
		}

		if _, err := courseSvc.Get(c.Context(), id, us.UserID); err != nil {
			return c.SendStatus(fiber.StatusForbidden)
		}

		c.Locals("course_id", id)

		return c.Next()
	}
}
