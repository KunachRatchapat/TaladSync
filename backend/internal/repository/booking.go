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

func (r *BookingRepository) GetByStallAndDate(ctx context.Context, stallID uint, date time.Time) (*domain.DailyBooking, error) {
	date = timeutil.MarketDate(date)
	var booking domain.DailyBooking
	err := r.db.WithContext(ctx).
		Preload("Stall").
		Where("stall_id = ? AND market_date = ?", stallID, date).
		First(&booking).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) TryLeaveStall(ctx context.Context, in domain.LeaveStallInput) (int64, error) {
	date := timeutil.MarketDate(in.MarketDate)
	res := r.db.WithContext(ctx).
		Model(&domain.DailyBooking{}).
		Where("stall_id = ? AND market_date = ? AND status = ? AND vendor_id = ?",
			in.StallID, date, domain.StatusExpected, in.VendorID).
		Updates(map[string]interface{}{
			"status":       domain.StatusAvailable,
			"vendor_id":    nil,
			"shop_name":    "",
			"product_type": "",
		})
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

//เอาไว้หา booking ตาม id
func (r *BookingRepository) FindByID(ctx context.Context, id uint) (*domain.DailyBooking, error) {
	var booking domain.DailyBooking
	err := r.db.WithContext(ctx).Preload("Stall").First(&booking, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) ListChecklist(ctx context.Context, date time.Time) ([]domain.DailyBooking, error) {
	date = timeutil.MarketDate(date)
	var bookings []domain.DailyBooking
	err := r.db.WithContext(ctx).
		Preload("Stall").
		Where("market_date = ? AND status IN ?",
			date, []domain.BookingStatus{domain.StatusExpected, domain.StatusBooked}).
		Order("stall_id ASC").
		Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

//เช็คว่ามาไหม
func (r *BookingRepository) SetAttendance(
	ctx context.Context,
	bookingID, staffID uint,
	attendance domain.AttendanceStatus,
	checkedAt time.Time,
) (int64, error) {
	res := r.db.WithContext(ctx).
		Model(&domain.DailyBooking{}).
		Where("id = ? AND attendance != ?", bookingID, attendance).
		Updates(map[string]interface{}{
			"attendance":           attendance,
			"checked_by_staff_id": staffID,
			"checked_at":           checkedAt,
		})
	return res.RowsAffected, res.Error
}

func (r *BookingRepository) SearchDirectory(ctx context.Context, date time.Time, q string) ([]domain.DirectoryEntry, error) {
	date = timeutil.MarketDate(date)
	statuses := []domain.BookingStatus{domain.StatusExpected, domain.StatusBooked}

	query := r.db.WithContext(ctx).
		Preload("Stall").
		Where("market_date = ? AND status IN ? AND shop_name != ''", date, statuses)

	if q != "" {
		pattern := "%" + q + "%"
		query = query.Where("shop_name ILIKE ? OR product_type ILIKE ?", pattern, pattern)
	}

	var bookings []domain.DailyBooking
	if err := query.Order("stall_id ASC").Find(&bookings).Error; err != nil {
		return nil, err
	}

	entries := make([]domain.DirectoryEntry, 0, len(bookings))
	for _, b := range bookings {
		if b.Stall == nil {
			continue
		}
		entries = append(entries, domain.DirectoryEntry{
			StallCode:   b.Stall.Code,
			Zone:        b.Stall.Zone,
			ShopName:    b.ShopName,
			ProductType: b.ProductType,
			PosX:        b.Stall.PosX,
			PosY:        b.Stall.PosY,
		})
	}
	return entries, nil
}

func (r *BookingRepository) AttendanceSummary(ctx context.Context, date time.Time) (domain.AttendanceSummary, error) {
	date = timeutil.MarketDate(date)
	statuses := []domain.BookingStatus{domain.StatusExpected, domain.StatusBooked}

	type row struct {
		Attendance domain.AttendanceStatus
		Count      int64
	}
	var rows []row
	err := r.db.WithContext(ctx).
		Model(&domain.DailyBooking{}).
		Select("attendance, count(*) as count").
		Where("market_date = ? AND status IN ?", date, statuses).
		Group("attendance").
		Scan(&rows).Error
	if err != nil {
		return domain.AttendanceSummary{}, err
	}

	var summary domain.AttendanceSummary
	for _, row := range rows {
		summary.Total += row.Count
		switch row.Attendance {
		case domain.AttendancePresent:
			summary.Present = row.Count
		case domain.AttendanceNoShow:
			summary.NoShow = row.Count
		case domain.AttendancePending:
			summary.Pending = row.Count
		}
	}
	return summary, nil
}

func isDuplicateKey(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
