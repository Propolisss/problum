package middleware

import (
	"context"
	"strconv"

	attemptDTO "problum/internal/attempt/service/dto"

	"github.com/gofiber/fiber/v3"
)

//go:generate go run go.uber.org/mock/mockgen -destination=mocks/mock_service.go -package=mocks problum/internal/middleware AttemptService,EnrollmentService,LessonService,ProblemService
type AttemptService interface {
	Get(context.Context, int) (*attemptDTO.Attempt, error)
}

func Attempt(attemptSvc AttemptService) fiber.Handler {
	return func(c fiber.Ctx) error {
		attemptID, err := strconv.Atoi(c.Params("attemptID"))
		if err != nil {
			return c.SendStatus(fiber.StatusForbidden)
		}

		userID, ok := c.Locals("user_id").(int)
		if !ok {
			return c.SendStatus(fiber.StatusForbidden)
		}

		attempt, err := attemptSvc.Get(c.Context(), attemptID)
		if err != nil {
			return c.SendStatus(fiber.StatusForbidden)
		}

		if attempt.UserID != userID {
			return c.SendStatus(fiber.StatusForbidden)
		}

		return c.Next()
	}
}
