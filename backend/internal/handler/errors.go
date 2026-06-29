package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"taladsync/internal/domain"
)

type errMapping struct {
	err     error
	status  int
	message string
}

var errMappings = []errMapping{
	{domain.ErrInvalidRole, fiber.StatusBadRequest, "invalid role"},
	{domain.ErrUserNotFound, fiber.StatusNotFound, "user not found"},
	{domain.ErrForbidden, fiber.StatusForbidden, "forbidden"},
	{errUnauthorized, fiber.StatusUnauthorized, "unauthorized"},
}

func RespondError(c *fiber.Ctx, err error) error {
	for _, m := range errMappings {
		if errors.Is(err, m.err) {
			return c.Status(m.status).JSON(fiber.Map{"error": m.message})
		}
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
}

var errUnauthorized = errors.New("unauthorized")

func badRequest(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
}
