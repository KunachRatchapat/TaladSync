package usecase

import (
	"context"
	"fmt"
	"time"

	"taladsync/internal/domain"
)

type BookingRepository interface {
	ListAvailable(ctx context.Context, date time.Time) ([]domain.DailyBooking, error)
	TryBookStall(ctx context.Context, in domain.BookStallInput) (int64, error)
	CountVendorBookingToday(ctx context.Context, vendorID uint, date time.Time) (int64, error)
	FindVendorBookedToday(ctx context.Context, vendorID uint, date time.Time) (*domain.DailyBooking, error)
}

type BookingUsecase struct {
	repo  BookingRepository
	audit AuditWriter
}

func NewBookingUsecase(repo BookingRepository, audit AuditWriter) *BookingUsecase {
	return &BookingUsecase{repo: repo, audit: audit}
}

func (u *BookingUsecase) ListAvailable(ctx context.Context, date time.Time) ([]domain.DailyBooking, error) {
	return u.repo.ListAvailable(ctx, date)
}

func (u *BookingUsecase) BookStall(ctx context.Context, in domain.BookStallInput) error {
	cnt, err := u.repo.CountVendorBookingToday(ctx, in.VendorID, in.MarketDate)
	if err != nil {
		return err
	}
	if cnt > 0 {
		return domain.ErrAlreadyHasStall
	}

	rows, err := u.repo.TryBookStall(ctx, in)
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrStallNotAvailable
	}

	booking, err := u.repo.FindVendorBookedToday(ctx, in.VendorID, in.MarketDate)
	if err != nil {
		return err
	}
	detail := fmt.Sprintf("stallId=%d shop=%s", in.StallID, in.ShopName)
	var bookingID *uint
	if booking != nil {
		bookingID = &booking.ID
		if booking.Stall != nil {
			detail = fmt.Sprintf("stall=%s shop=%s", booking.Stall.Code, in.ShopName)
		}
	}
	return u.audit.Create(ctx, domain.AuditLog{
		ActorUserID:    in.VendorID,
		DailyBookingID: bookingID,
		Action:         domain.AuditBook,
		Detail:         detail,
	})
}

func (u *BookingUsecase) GetMyBookingToday(ctx context.Context, vendorID uint, date time.Time) (*domain.DailyBooking, error) {
	return u.repo.FindVendorBookedToday(ctx, vendorID, date)
}
