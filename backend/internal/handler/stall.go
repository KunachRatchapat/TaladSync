package handler

import (
	"github.com/gofiber/fiber/v2"

	"taladsync/internal/usecase"
)

type StallHandler struct {
	uc *usecase.StallUsecase
}

func NewStallHandler(uc *usecase.StallUsecase) *StallHandler {
	return &StallHandler{uc: uc}
}

func (h *StallHandler) ListZones(c *fiber.Ctx) error {
	zones, err := h.uc.ListZones(c.Context())
	if err != nil {
		return RespondError(c, err)
	}
	return c.JSON(fiber.Map{"zones": zones})
}
