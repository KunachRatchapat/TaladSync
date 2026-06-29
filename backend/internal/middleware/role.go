package middleware

import (
	"github.com/gofiber/fiber/v2"

	"taladsync/internal/domain"
)

func RequireRoles(roles ...domain.Role) fiber.Handler {
	allowed := make(map[domain.Role]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		claims, err := GetClaims(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		if _, ok := allowed[claims.Role]; !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		return c.Next()
	}
}
