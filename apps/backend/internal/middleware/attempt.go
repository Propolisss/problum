package middleware

import (
	"context"
	"strconv"

	attemptDTO "problum/internal/attempt/service/dto"

	"github.com/gofiber/fiber/v3"
)

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
