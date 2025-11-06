package middleware

import (
	"context"
	"fmt"
	"strings"

	"problum/internal/model"
	"problum/internal/redis"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
)

func Auth(rdb *redis.Redis) fiber.Handler {
	return func(c fiber.Ctx) error {
		access := c.Get("Authorization")
		if access == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		access = strings.TrimPrefix(access, "Bearer ")
		usJSON, err := rdb.Get(c.Context(), fmt.Sprintf("user_sessions:%s", access))
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		session := &model.UserSession{}
		if err := sonic.Unmarshal(usJSON, session); err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals("access_token", access)
		c.Locals("user_session", session)
		c.Locals("user_id", session.UserID)

		ctx := context.WithValue(c.Context(), "user_id", session.UserID)
		c.SetContext(ctx)

		return c.Next()
	}
}
