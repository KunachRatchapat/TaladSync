package handler

import (
	"github.com/gofiber/fiber/v2"

	"taladsync/internal/domain"
	"taladsync/internal/middleware"
	"taladsync/internal/usecase"
)

type VendorHandler struct {
	uc *usecase.VendorUsecase
}

func NewVendorHandler(uc *usecase.VendorUsecase) *VendorHandler {
	return &VendorHandler{uc: uc}
}

type leaveRequest struct {
	StallID uint   `json:"stallId"`
	Date    string `json:"date,omitempty"`
}

func (h *VendorHandler) Leave(c *fiber.Ctx) error {
	claims, err := middleware.GetClaims(c)
	if err != nil {
		return RespondError(c, errUnauthorized)
	}

	var req leaveRequest
	if err := c.BodyParser(&req); err != nil || req.StallID == 0 {
		return badRequest(c)
	}

	marketDate, err := parseMarketDate(req.Date)
	if err != nil {
		return badRequest(c)
	}

	err = h.uc.Leave(c.Context(), domain.LeaveStallInput{
		StallID:    req.StallID,
		VendorID:   claims.UserID,
		MarketDate: marketDate,
	})
	if err != nil {
		return RespondError(c, err)
	}

	return c.JSON(fiber.Map{"status": "available", "stallId": req.StallID})
}
