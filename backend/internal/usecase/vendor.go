package usecase

import (
	"context"
	"fmt"
	"time"

	"taladsync/internal/domain"
)

type VendorLeaveRepository interface {
	GetByStallAndDate(ctx context.Context, stallID uint, date time.Time) (*domain.DailyBooking, error)
	LeaveStallWithAudit(ctx context.Context, in domain.LeaveStallInput, log domain.AuditLog) (int64, error)
	ListOwnedStallBookings(ctx context.Context, vendorID uint, date time.Time) ([]domain.DailyBooking, error)
}

type VendorUsecase struct {
	bookings VendorLeaveRepository
}

func NewVendorUsecase(bookings VendorLeaveRepository) *VendorUsecase {
	return &VendorUsecase{bookings: bookings}
}

func (u *VendorUsecase) GetMyStalls(ctx context.Context, vendorID uint, date time.Time) ([]domain.DailyBooking, error) {
	return u.bookings.ListOwnedStallBookings(ctx, vendorID, date)
}

func (u *VendorUsecase) Leave(ctx context.Context, in domain.LeaveStallInput) error {
	booking, err := u.bookings.GetByStallAndDate(ctx, in.StallID, in.MarketDate)
	if err != nil {
		return err
	}
	if !ownsStall(booking, in.VendorID) {
		return domain.ErrNotOwnStall
	}
	if booking.Status != domain.StatusExpected {
		return nil
	}

	detail := fmt.Sprintf("stallId=%d", in.StallID)
	if booking.Stall != nil {
		detail = fmt.Sprintf("stall=%s", booking.Stall.Code)
	}
	bookingID := booking.ID
	log := domain.AuditLog{
		ActorUserID:    in.VendorID,
		DailyBookingID: &bookingID,
		Action:         domain.AuditLeave,
		Detail:         detail,
	}

	rows, err := u.bookings.LeaveStallWithAudit(ctx, in, log)
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrLeaveNotAllowed
	}
	return nil
}

func ownsStall(booking *domain.DailyBooking, vendorID uint) bool {
	return booking.Stall != nil &&
		booking.Stall.OwnerVendorID != nil &&
		*booking.Stall.OwnerVendorID == vendorID
}
