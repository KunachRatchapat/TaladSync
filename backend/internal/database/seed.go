package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"taladsync/internal/domain"
	"taladsync/internal/pkg/timeutil"
)

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

	return db.Transaction(func(tx *gorm.DB) error {
		regular := domain.User{
			NationalID: "1100700000001",
			Name:       "แม่ค้าประจำ01",
			Phone:      "0639999999",
			Role:       domain.RoleRegular,
		}
		walkin := domain.User{
			NationalID: "1100700000002",
			Name:       "ขาจรทดสอบ",
			Phone:      "0810000002",
			Role:       domain.RoleWalkin,
		}
		staff := domain.User{
			NationalID: "1100700000003",
			Name:       "เจ้าหน้าที่สมศรี",
			Phone:      "0810000003",
			Role:       domain.RoleStaff,
		}
		if err := tx.Create([]*domain.User{&regular, &walkin, &staff}).Error; err != nil {
			return fmt.Errorf("create users: %w", err)
		}

		stalls := make([]domain.Stall, 0, 20)
		for i := 1; i <= 12; i++ {
			ownerID := regular.ID
			stalls = append(stalls, domain.Stall{
				Code:          fmt.Sprintf("A-%02d", i),
				Zone:          "Zone A",
				OwnerVendorID: &ownerID,
				PosX:          float64(i),
				PosY:          1,
			})
		}
		for i := 1; i <= 8; i++ {
			stalls = append(stalls, domain.Stall{
				Code:  fmt.Sprintf("P-%02d", i),
				Zone:  "Pool",
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
			3, len(stalls), len(bookings), today.Format("2006-01-02"))
		return nil
	})
}
