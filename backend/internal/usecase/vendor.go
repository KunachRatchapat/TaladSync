package usecase

import (
	"context"
	"fmt"
	"time"

	"taladsync/internal/domain"
)

type VendorLeaveRepository interface {
	GetByStallAndDate(ctx context.Context, stallID uint, date time.Time) (*domain.DailyBooking, error)
	TryLeaveStall(ctx context.Context, in domain.LeaveStallInput) (int64, error)
}

type AuditWriter interface {
	Create(ctx context.Context, log domain.AuditLog) error
}

type VendorUsecase struct {
	bookings VendorLeaveRepository
	audit    AuditWriter
}

func NewVendorUsecase(bookings VendorLeaveRepository, audit AuditWriter) *VendorUsecase {
	return &VendorUsecase{bookings: bookings, audit: audit}
}

func (u *VendorUsecase) Leave(ctx context.Context, in domain.LeaveStallInput) error {
	booking, err := u.bookings.GetByStallAndDate(ctx, in.StallID, in.MarketDate)
	if err != nil {
		return err
	}
	if booking.VendorID == nil || *booking.VendorID != in.VendorID {
		return domain.ErrNotOwnStall
	}
	if booking.Status != domain.StatusExpected {
		return domain.ErrLeaveNotAllowed
	}

	rows, err := u.bookings.TryLeaveStall(ctx, in)
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrLeaveNotAllowed
	}

	detail := fmt.Sprintf("stallId=%d", in.StallID)
	if booking.Stall != nil {
		detail = fmt.Sprintf("stall=%s", booking.Stall.Code)
	}
	bookingID := booking.ID
	return u.audit.Create(ctx, domain.AuditLog{
		ActorUserID:    in.VendorID,
		DailyBookingID: &bookingID,
		Action:         domain.AuditLeave,
		Detail:         detail,
	})
}
