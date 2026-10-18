package usecase

import (
	"context"
	"fmt"
	"time"

	"taladsync/internal/domain"
	"taladsync/internal/pkg/timeutil"
)

type StaffBookingRepository interface {
	FindByID(ctx context.Context, id uint) (*domain.DailyBooking, error)
	ListChecklist(ctx context.Context, date time.Time) ([]domain.DailyBooking, error)
	SetAttendanceWithAudit(ctx context.Context, bookingID, staffID uint, attendance domain.AttendanceStatus, checkedAt time.Time, log domain.AuditLog) (int64, error)
	AttendanceSummary(ctx context.Context, date time.Time) (domain.AttendanceSummary, error)
}

type StaffUsecase struct {
	bookings StaffBookingRepository
}

func NewStaffUsecase(bookings StaffBookingRepository) *StaffUsecase {
	return &StaffUsecase{bookings: bookings}
}

func (u *StaffUsecase) GetChecklist(ctx context.Context, date time.Time) ([]domain.DailyBooking, error) {
	return u.bookings.ListChecklist(ctx, date)
}

func (u *StaffUsecase) GetSummary(ctx context.Context, date time.Time) (domain.AttendanceSummary, error) {
	return u.bookings.AttendanceSummary(ctx, date)
}

func (u *StaffUsecase) MarkAttendance(ctx context.Context, in domain.MarkAttendanceInput) (*domain.DailyBooking, error) {
	booking, err := u.bookings.FindByID(ctx, in.BookingID)
	if err != nil {
		return nil, err
	}
	if booking.Status != domain.StatusExpected && booking.Status != domain.StatusBooked {
		return nil, domain.ErrInvalidAttendanceTarget
	}

	target := domain.AttendanceNoShow
	action := domain.AuditNoshow
	if in.Present {
		target = domain.AttendancePresent
		action = domain.AuditCheckin
	}

	if booking.Attendance == target {
		return booking, nil
	}

	bookingID := in.BookingID
	log := domain.AuditLog{
		ActorUserID:    in.StaffID,
		DailyBookingID: &bookingID,
		Action:         action,
		Detail:         fmt.Sprintf("attendance=%s", target),
	}

	checkedAt := timeutil.Now()
	rows, err := u.bookings.SetAttendanceWithAudit(ctx, in.BookingID, in.StaffID, target, checkedAt, log)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		updated, err := u.bookings.FindByID(ctx, in.BookingID)
		if err != nil {
			return nil, err
		}
		if updated.Attendance == target {
			return updated, nil
		}
		return nil, domain.ErrInvalidAttendanceTarget
	}

	return u.bookings.FindByID(ctx, in.BookingID)
}
