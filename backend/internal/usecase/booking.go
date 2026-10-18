package usecase

import (
	"context"
	"fmt"
	"time"

	"taladsync/internal/domain"
)

type BookingRepository interface {
	ListAvailable(ctx context.Context, date time.Time) ([]domain.DailyBooking, error)
	BookStallWithAudit(ctx context.Context, in domain.BookStallInput, log domain.AuditLog) (int64, error)
	CountVendorBookingToday(ctx context.Context, vendorID uint, date time.Time) (int64, error)
	FindVendorBookedToday(ctx context.Context, vendorID uint, date time.Time) (*domain.DailyBooking, error)
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
		existing, err := u.repo.FindVendorBookedToday(ctx, in.VendorID, in.MarketDate)
		if err != nil {
			return err
		}
		if existing != nil && existing.StallID == in.StallID {
			return nil
		}
		return domain.ErrAlreadyHasStall
	}

	log := domain.AuditLog{
		ActorUserID: in.VendorID,
		Action:      domain.AuditBook,
		Detail:      fmt.Sprintf("stallId=%d shop=%s", in.StallID, in.ShopName),
	}

	rows, err := u.repo.BookStallWithAudit(ctx, in, log)
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrStallNotAvailable
	}
	return nil
}

func (u *BookingUsecase) GetMyBookingToday(ctx context.Context, vendorID uint, date time.Time) (*domain.DailyBooking, error) {
	return u.repo.FindVendorBookedToday(ctx, vendorID, date)
}
