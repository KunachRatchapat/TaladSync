//go:build integration

package usecase

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"

	"taladsync/internal/database"
	"taladsync/internal/domain"
	"taladsync/internal/pkg/timeutil"
	"taladsync/internal/repository"
)

const concurrentVendors = 50

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("docker not available — start Docker Desktop then rerun: %v", r)
		}
	}()
	return startPostgresContainer(t)
}

func startPostgresContainer(t *testing.T) *gorm.DB {
	t.Helper()
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("taladsync_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("terminate container: %v", err)
		}
	})

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}

	db, err := database.Connect(connStr)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

type concurrentFixture struct {
	stallID    uint
	vendorIDs  []uint
	marketDate time.Time
}

func seedConcurrentFixture(t *testing.T, db *gorm.DB) concurrentFixture {
	t.Helper()
	marketDate := timeutil.Today()

	stall := domain.Stall{
		Code: "P-TEST",
		Zone: "Pool",
		PosX: 99,
		PosY: 2,
	}
	if err := db.Create(&stall).Error; err != nil {
		t.Fatalf("create stall: %v", err)
	}

	booking := domain.DailyBooking{
		StallID:    stall.ID,
		MarketDate: marketDate,
		Status:     domain.StatusAvailable,
		Attendance: domain.AttendancePending,
	}
	if err := db.Create(&booking).Error; err != nil {
		t.Fatalf("create booking: %v", err)
	}

	vendorIDs := make([]uint, 0, concurrentVendors)
	for i := 1; i <= concurrentVendors; i++ {
		user := domain.User{
			NationalID: fmt.Sprintf("1199900%06d", i),
			Name:       fmt.Sprintf("ขาจรเทส %02d", i),
			Phone:      fmt.Sprintf("089%07d", i),
			Role:       domain.RoleWalkin,
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("create user %d: %v", i, err)
		}
		vendorIDs = append(vendorIDs, user.ID)
	}

	return concurrentFixture{
		stallID:    stall.ID,
		vendorIDs:  vendorIDs,
		marketDate: marketDate,
	}
}

func TestBookStall_Concurrent_OnlyOneWins(t *testing.T) {
	db := setupTestDB(t)
	fx := seedConcurrentFixture(t, db)

	uc := NewBookingUsecase(repository.NewBookingRepository(db), repository.NewAuditRepository(db))
	ctx := context.Background()

	var wg sync.WaitGroup
	var successCount int64
	var conflictCount int64

	start := make(chan struct{})
	for _, vid := range fx.vendorIDs {
		wg.Add(1)
		go func(vendorID uint) {
			defer wg.Done()
			<-start
			err := uc.BookStall(ctx, domain.BookStallInput{
				StallID:     fx.stallID,
				VendorID:    vendorID,
				MarketDate:  fx.marketDate,
				ShopName:    "ร้านทดสอบ",
				ProductType: "ของกิน",
			})
			if err == nil {
				atomic.AddInt64(&successCount, 1)
			} else if errors.Is(err, domain.ErrStallNotAvailable) {
				atomic.AddInt64(&conflictCount, 1)
			} else {
				t.Errorf("vendor %d: unexpected error: %v", vendorID, err)
			}
		}(vid)
	}
	close(start)
	wg.Wait()

	if successCount != 1 {
		t.Fatalf("expected 1 success, got %d", successCount)
	}
	if conflictCount != concurrentVendors-1 {
		t.Fatalf("expected %d conflicts, got %d", concurrentVendors-1, conflictCount)
	}

	var booked int64
	if err := db.Model(&domain.DailyBooking{}).
		Where("stall_id = ? AND market_date = ? AND status = ?",
			fx.stallID, fx.marketDate, domain.StatusBooked).
		Count(&booked).Error; err != nil {
		t.Fatalf("count booked: %v", err)
	}
	if booked != 1 {
		t.Fatalf("expected 1 BOOKED row, got %d", booked)
	}
}
