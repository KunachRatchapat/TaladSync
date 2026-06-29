package domain

import (
	"errors"
	"time"
)

var (
	ErrStallNotAvailable = errors.New("stall is not available")
	ErrAlreadyHasStall   = errors.New("vendor already booked today")
	ErrBookingNotFound   = errors.New("booking not found")
	ErrInvalidRole       = errors.New("invalid role")
	ErrForbidden         = errors.New("forbidden")
)

type BookStallInput struct {
	StallID     uint
	VendorID    uint
	MarketDate  time.Time
	ShopName    string
	ProductType string
}
