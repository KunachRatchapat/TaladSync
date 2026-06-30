package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"taladsync/internal/domain"
	"taladsync/internal/middleware"
	"taladsync/internal/pkg/timeutil"
	"taladsync/internal/usecase"
)

type BookingHandler struct {
	uc *usecase.BookingUsecase
}

func NewBookingHandler(uc *usecase.BookingUsecase) *BookingHandler {
	return &BookingHandler{uc: uc}
}

type bookRequest struct {
	StallID     uint   `json:"stallId"`
	ShopName    string `json:"shopName"`
	ProductType string `json:"productType"`
}

func (h *BookingHandler) ListAvailable(c *fiber.Ctx) error {
	date, err := parseMarketDate(c.Query("date"))
	if err != nil {
		return badRequest(c)
	}

	bookings, err := h.uc.ListAvailable(c.Context(), date)
	if err != nil {
		return RespondError(c, err)
	}

	return c.JSON(fiber.Map{"stalls": bookings})
}

func (h *BookingHandler) Book(c *fiber.Ctx) error {
	claims, err := middleware.GetClaims(c)
	if err != nil {
		return RespondError(c, errUnauthorized)
	}

	var req bookRequest
	if err := c.BodyParser(&req); err != nil ||
		req.StallID == 0 || req.ShopName == "" || req.ProductType == "" {
		return badRequest(c)
	}

	err = h.uc.BookStall(c.Context(), domain.BookStallInput{
		StallID:     req.StallID,
		VendorID:    claims.UserID,
		MarketDate:  timeutil.Today(),
		ShopName:    req.ShopName,
		ProductType: req.ProductType,
	})
	if err != nil {
		return RespondError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "booked"})
}

func parseMarketDate(raw string) (time.Time, error) {
	if raw == "" {
		return timeutil.Today(), nil
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Time{}, err
	}
	return timeutil.MarketDate(t), nil
}
