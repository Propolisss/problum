package middleware

import (
	"fmt"
	"strings"

	"problum/internal/redis"

	"github.com/gofiber/fiber/v3"
)

func Auth(rdb *redis.Redis) fiber.Handler {
	return func(c fiber.Ctx) error {
		access := c.Get("Authorization")
		if access == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		access = strings.TrimPrefix(access, "Bearer ")
		if _, err := rdb.Get(c.Context(), fmt.Sprintf("user_sessions:%s", access)); err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}
