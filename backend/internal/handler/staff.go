package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"taladsync/internal/domain"
	"taladsync/internal/middleware"
	"taladsync/internal/usecase"
)

type StaffHandler struct {
	uc *usecase.StaffUsecase
}

func NewStaffHandler(uc *usecase.StaffUsecase) *StaffHandler {
	return &StaffHandler{uc: uc}
}

func (h *StaffHandler) Checklist(c *fiber.Ctx) error {
	date, err := parseMarketDate(c.Query("date"))
	if err != nil {
		return badRequest(c)
	}

	items, err := h.uc.GetChecklist(c.Context(), date)
	if err != nil {
		return RespondError(c, err)
	}

	return c.JSON(fiber.Map{"items": items})
}

func (h *StaffHandler) Summary(c *fiber.Ctx) error {
	date, err := parseMarketDate(c.Query("date"))
	if err != nil {
		return badRequest(c)
	}

	summary, err := h.uc.GetSummary(c.Context(), date)
	if err != nil {
		return RespondError(c, err)
	}

	return c.JSON(summary)
}

type attendanceRequest struct {
	Present bool `json:"present"`
}

func (h *StaffHandler) MarkAttendance(c *fiber.Ctx) error {
	claims, err := middleware.GetClaims(c)
	if err != nil {
		return RespondError(c, errUnauthorized)
	}

	bookingID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || bookingID == 0 {
		return badRequest(c)
	}

	var req attendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return badRequest(c)
	}

	booking, err := h.uc.MarkAttendance(c.Context(), domain.MarkAttendanceInput{
		BookingID: uint(bookingID),
		StaffID:   claims.UserID,
		Present:   req.Present,
	})
	if err != nil {
		return RespondError(c, err)
	}

	return c.JSON(booking)
}
