package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"taladsync/internal/auth"
)

const claimsLocalKey = "claims"

func Auth(tokens auth.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, err := bearerToken(c.Get("Authorization"))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		claims, err := tokens.Parse(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		c.Locals(claimsLocalKey, claims)
		return c.Next()
	}
}

func GetClaims(c *fiber.Ctx) (*auth.Claims, error) {
	v := c.Locals(claimsLocalKey)
	if v == nil {
		return nil, errors.New("missing claims")
	}
	claims, ok := v.(*auth.Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}

func bearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("missing authorization header")
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", errors.New("invalid authorization header")
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", errors.New("empty token")
	}
	return token, nil
}
