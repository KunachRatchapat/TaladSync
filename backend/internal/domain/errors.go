package domain

import (
	"errors"
	"time"
)

var (
	ErrStallNotAvailable         = errors.New("stall is not available")
	ErrAlreadyHasStall           = errors.New("vendor already booked today")
	ErrBookingNotFound           = errors.New("booking not found")
	ErrInvalidRole               = errors.New("invalid role")
	ErrForbidden                 = errors.New("forbidden")
	ErrUserNotFound              = errors.New("user not found")
	ErrLeaveNotAllowed           = errors.New("leave not allowed for this stall")
	ErrNotOwnStall               = errors.New("not your stall")
	ErrInvalidAttendanceTarget   = errors.New("attendance not allowed for this booking")
)

type BookStallInput struct {
	StallID     uint
	VendorID    uint
	MarketDate  time.Time
	ShopName    string
	ProductType string
}

type LeaveStallInput struct {
	StallID    uint
	VendorID   uint
	MarketDate time.Time
}

type MarkAttendanceInput struct {
	BookingID uint
	StaffID   uint
	Present   bool
}

type SetAttendanceInput struct {
	BookingID  uint
	StaffID    uint
	Attendance AttendanceStatus
	CheckedAt  time.Time
}

type AttendanceSummary struct {
	Present int64
	NoShow  int64
	Pending int64
	Total   int64
}

type DirectoryEntry struct {
	StallCode   string  `json:"stallCode"`
	Zone        string  `json:"zone"`
	ShopName    string  `json:"shopName"`
	ProductType string  `json:"productType"`
	PosX        float64 `json:"posX"`
	PosY        float64 `json:"posY"`
}
