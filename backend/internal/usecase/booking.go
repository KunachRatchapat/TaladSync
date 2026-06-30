package usecase

import (
	"context"
	"time"

	"taladsync/internal/domain"
)

type BookingRepository interface {
	ListAvailable(ctx context.Context, date time.Time) ([]domain.DailyBooking, error)
	TryBookStall(ctx context.Context, in domain.BookStallInput) (int64, error)
	CountVendorBookingToday(ctx context.Context, vendorID uint, date time.Time) (int64, error)
}

type BookingUsecase struct {
	repo BookingRepository
}

func NewBookingUsecase(repo BookingRepository) *BookingUsecase {
	return &BookingUsecase{repo: repo}
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
	return nil
}
