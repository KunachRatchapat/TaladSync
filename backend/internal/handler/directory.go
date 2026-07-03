package handler

import (
	"github.com/gofiber/fiber/v2"

	"taladsync/internal/usecase"
)

type DirectoryHandler struct {
	uc *usecase.DirectoryUsecase
}

func NewDirectoryHandler(uc *usecase.DirectoryUsecase) *DirectoryHandler {
	return &DirectoryHandler{uc: uc}
}

func (h *DirectoryHandler) Search(c *fiber.Ctx) error {
	date, err := parseMarketDate(c.Query("date"))
	if err != nil {
		return badRequest(c)
	}

	results, err := h.uc.Search(c.Context(), date, c.Query("q"))
	if err != nil {
		return RespondError(c, err)
	}

	return c.JSON(fiber.Map{"results": results})
}
