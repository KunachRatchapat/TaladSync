package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"taladsync/internal/domain"
	"taladsync/internal/pkg/timeutil"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) ListAvailable(ctx context.Context, date time.Time) ([]domain.DailyBooking, error) {
	date = timeutil.MarketDate(date)
	var bookings []domain.DailyBooking
	err := r.db.WithContext(ctx).
		Preload("Stall").
		Where("market_date = ? AND status = ?", date, domain.StatusAvailable).
		Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *BookingRepository) TryBookStall(ctx context.Context, in domain.BookStallInput) (int64, error) {
	date := timeutil.MarketDate(in.MarketDate)
	res := r.db.WithContext(ctx).
		Model(&domain.DailyBooking{}).
		Where("stall_id = ? AND market_date = ? AND status = ?",
			in.StallID, date, domain.StatusAvailable).
		Updates(map[string]interface{}{
			"status":       domain.StatusBooked,
			"vendor_id":    in.VendorID,
			"shop_name":    in.ShopName,
			"product_type": in.ProductType,
			"booked_at":    timeutil.Now(),
		})
	if res.Error != nil {
		if isDuplicateKey(res.Error) {
			return 0, domain.ErrAlreadyHasStall
		}
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func (r *BookingRepository) CountVendorBookingToday(ctx context.Context, vendorID uint, date time.Time) (int64, error) {
	date = timeutil.MarketDate(date)
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.DailyBooking{}).
		Where("vendor_id = ? AND market_date = ? AND status = ?",
			vendorID, date, domain.StatusBooked).
		Count(&count).Error
	return count, err
}

func isDuplicateKey(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
