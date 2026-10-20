package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"taladsync/internal/domain"
	"taladsync/internal/pkg/timeutil"
)

// Reseed ล้างทุกตารางแล้วใส่ข้อมูล dev สดใหม่ (สำหรับ flag -reseed / เดโม)
func Reseed(db *gorm.DB) error {
	if err := db.Exec(
		"TRUNCATE audit_logs, daily_bookings, stalls, users RESTART IDENTITY CASCADE",
	).Error; err != nil {
		return fmt.Errorf("truncate: %w", err)
	}
	log.Println("reseed: tables truncated")
	if err := SeedIfEmpty(db); err != nil {
		return err
	}
	return SeedDailyBookingsForToday(db)
}

// SeedIfEmpty ใส่ข้อมูล dev เมื่อ users ว่างอยู่
func SeedIfEmpty(db *gorm.DB) error {
	var count int64
	if err := db.Model(&domain.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count users: %w", err)
	}
	if count > 0 {
		log.Println("seed skipped (users already exist)")
		return nil
	}

	today := timeutil.Today()

	const (
		zoneAStalls = 15 // โซน A — แผงแม่ค้าประจำ (1 คน : 1 แผง)
		zoneBStalls = 15 // โซน B — แผงกลาง ว่างให้ขาจรจอง
	)

	return db.Transaction(func(tx *gorm.DB) error {
		
		regulars := make([]*domain.User, zoneAStalls)
		for i := 0; i < zoneAStalls; i++ {
			regulars[i] = &domain.User{
				NationalID: fmt.Sprintf("11007000000%02d", i+1),
				Name:       fmt.Sprintf("แม่ค้าประจำ%02d", i+1),
				Phone:      fmt.Sprintf("063%07d", i+1),
				Role:       domain.RoleRegular,
			}
		}
		walkin := domain.User{
			NationalID: "1100700009002",
			Name:       "ป้าแดง ขายผัก",
			Phone:      "0810000002",
			Role:       domain.RoleWalkin,
		}
		staff := domain.User{
			NationalID: "1100700009003",
			Name:       "เจ้าหน้าที่สมศรี",
			Phone:      "0810000003",
			Role:       domain.RoleStaff,
		}
		users := append([]*domain.User{}, regulars...)
		users = append(users, &walkin, &staff)
		if err := tx.Create(&users).Error; err != nil {
			return fmt.Errorf("create users: %w", err)
		}

		stalls := make([]domain.Stall, 0, zoneAStalls+zoneBStalls)
		for i := 1; i <= zoneAStalls; i++ {
			ownerID := regulars[i-1].ID // แต่ละแผงเจ้าของคนละคน
			stalls = append(stalls, domain.Stall{
				Code:          fmt.Sprintf("A-%02d", i),
				Zone:          "โซน A",
				OwnerVendorID: &ownerID,
				PosX:          float64(i),
				PosY:          1,
			})
		}
		// โซน B ไม่มีเจ้าของ → daily_booking เป็น AVAILABLE ให้ขาจรจอง
		for i := 1; i <= zoneBStalls; i++ {
			stalls = append(stalls, domain.Stall{
				Code:  fmt.Sprintf("B-%02d", i),
				Zone:  "โซน B",
				PosX:  float64(i),
				PosY:  2,
			})
		}
		if err := tx.Create(&stalls).Error; err != nil {
			return fmt.Errorf("create stalls: %w", err)
		}

		bookings := make([]domain.DailyBooking, 0, len(stalls))
		for _, s := range stalls {
			b := domain.DailyBooking{
				StallID:    s.ID,
				MarketDate: today,
				Attendance: domain.AttendancePending,
			}
			if s.OwnerVendorID != nil {
				b.Status = domain.StatusExpected
				b.VendorID = s.OwnerVendorID
				b.ShopName = "ร้าน" + s.Code
				b.ProductType = "ของสด"
			} else {
				b.Status = domain.StatusAvailable
			}
			bookings = append(bookings, b)
		}
		if err := tx.Create(&bookings).Error; err != nil {
			return fmt.Errorf("create daily bookings: %w", err)
		}

		log.Printf("seed complete: %d users, %d stalls, %d daily_bookings for %s",
			len(users), len(stalls), len(bookings), today.Format("2006-01-02"))
		return nil
	})
}

func SeedDailyBookingsForToday(db *gorm.DB) error {
	today := timeutil.Today()

	var stalls []domain.Stall
	if err := db.Find(&stalls).Error; err != nil {
		return fmt.Errorf("load stalls: %w", err)
	}
	if len(stalls) == 0 {
		return nil
	}

	var existing []domain.DailyBooking
	if err := db.Where("market_date = ?", today).Find(&existing).Error; err != nil {
		return fmt.Errorf("load today bookings: %w", err)
	}
	have := make(map[uint]bool, len(existing))
	for _, b := range existing {
		have[b.StallID] = true
	}

	toCreate := make([]domain.DailyBooking, 0)
	for _, s := range stalls {
		if have[s.ID] {
			continue
		}
		b := domain.DailyBooking{
			StallID:    s.ID,
			MarketDate: today,
			Attendance: domain.AttendancePending,
		}
		if s.OwnerVendorID != nil {
			b.Status = domain.StatusExpected
			b.VendorID = s.OwnerVendorID
			b.ShopName = "ร้าน" + s.Code
			b.ProductType = "ของสด"
		} else {
			b.Status = domain.StatusAvailable
		}
		toCreate = append(toCreate, b)
	}

	if len(toCreate) == 0 {
		log.Printf("daily_bookings for %s already present", today.Format("2006-01-02"))
		return nil
	}
	if err := db.Create(&toCreate).Error; err != nil {
		return fmt.Errorf("create today bookings: %w", err)
	}
	log.Printf("seeded %d daily_bookings for %s", len(toCreate), today.Format("2006-01-02"))
	return nil
}
