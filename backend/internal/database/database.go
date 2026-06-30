package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"taladsync/internal/domain"
)

var bangkok *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(fmt.Sprintf("load Asia/Bangkok: %v", err))
	}
	bangkok = loc
}

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		NowFunc: func() time.Time {
			return time.Now().In(bangkok)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	log.Println("database connect successful")

	return db, nil
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Stall{},
		&domain.DailyBooking{},
		&domain.AuditLog{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	// กัน vendor จองมากกว่า 1 แผง/วัน แม้ request พร้อมกัน (race หลัง conditional UPDATE)
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_vendor_one_booked_per_day
		ON daily_bookings (vendor_id, market_date)
		WHERE status = 'BOOKED' AND vendor_id IS NOT NULL
	`).Error; err != nil {
		return fmt.Errorf("create partial unique index: %w", err)
	}
	log.Println("database migrated")
	return nil
}
