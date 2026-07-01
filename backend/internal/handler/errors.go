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
	{domain.ErrBookingNotFound, fiber.StatusNotFound, "booking not found"},
	{domain.ErrForbidden, fiber.StatusForbidden, "forbidden"},
	{domain.ErrNotOwnStall, fiber.StatusForbidden, "not your stall"},
	{domain.ErrStallNotAvailable, fiber.StatusConflict, "stall is not available"},
	{domain.ErrAlreadyHasStall, fiber.StatusConflict, "vendor already booked today"},
	{domain.ErrLeaveNotAllowed, fiber.StatusConflict, "leave not allowed for this stall"},
	{domain.ErrInvalidAttendanceTarget, fiber.StatusConflict, "attendance not allowed for this booking"},
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
