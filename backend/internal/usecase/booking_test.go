package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"taladsync/internal/domain"
)

type mockBookingRepo struct {
	listAvailable     []domain.DailyBooking
	listAvailableErr  error
	tryRows           int64
	tryErr            error
	countToday        int64
	countTodayErr     error
	tryBookCalled     bool
	countTodayCalled  bool
	lastTryBookInput  domain.BookStallInput
	lastCountVendorID uint
	lastCountDate     time.Time
	findBooking       *domain.DailyBooking
	findBookingErr    error
	lastAuditLog      domain.AuditLog
}

func (m *mockBookingRepo) ListAvailable(ctx context.Context, date time.Time) ([]domain.DailyBooking, error) {
	return m.listAvailable, m.listAvailableErr
}

func (m *mockBookingRepo) BookStallWithAudit(ctx context.Context, in domain.BookStallInput, log domain.AuditLog) (int64, error) {
	m.tryBookCalled = true
	m.lastTryBookInput = in
	m.lastAuditLog = log
	return m.tryRows, m.tryErr
}

func (m *mockBookingRepo) CountVendorBookingToday(ctx context.Context, vendorID uint, date time.Time) (int64, error) {
	m.countTodayCalled = true
	m.lastCountVendorID = vendorID
	m.lastCountDate = date
	return m.countToday, m.countTodayErr
}

func (m *mockBookingRepo) FindVendorBookedToday(ctx context.Context, vendorID uint, date time.Time) (*domain.DailyBooking, error) {
	return m.findBooking, m.findBookingErr
}

func TestBookStall_Success(t *testing.T) {
	repo := &mockBookingRepo{countToday: 0, tryRows: 1}
	uc := NewBookingUsecase(repo)

	in := domain.BookStallInput{
		StallID: 13, VendorID: 2,
		MarketDate: time.Date(2026, 6, 29, 0, 0, 0, 0, time.UTC),
		ShopName:   "ร้านทด", ProductType: "ผลไม้",
	}

	if err := uc.BookStall(context.Background(), in); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if !repo.countTodayCalled || !repo.tryBookCalled {
		t.Fatal("expected CountVendorBookingToday and BookStallWithAudit to be called")
	}
	if repo.lastTryBookInput.StallID != 13 || repo.lastTryBookInput.VendorID != 2 {
		t.Fatalf("unexpected BookStallWithAudit input: %+v", repo.lastTryBookInput)
	}
	if repo.lastAuditLog.Action != domain.AuditBook {
		t.Fatalf("expected BOOK audit, got %q", repo.lastAuditLog.Action)
	}
	if repo.lastAuditLog.ActorUserID != 2 {
		t.Fatalf("expected actor=2, got %d", repo.lastAuditLog.ActorUserID)
	}
}

func TestBookStall_RetrySameStall(t *testing.T) {
	repo := &mockBookingRepo{
		countToday:  1,
		findBooking: &domain.DailyBooking{StallID: 13, Status: domain.StatusBooked},
	}
	uc := NewBookingUsecase(repo)

	err := uc.BookStall(context.Background(), domain.BookStallInput{
		StallID: 13, VendorID: 2, MarketDate: time.Now(),
	})
	if err != nil {
		t.Fatalf("expected nil on retry of the same stall, got %v", err)
	}
	if repo.tryBookCalled {
		t.Fatal("BookStallWithAudit should not be called on retry")
	}
}

func TestBookStall_AlreadyHasStall(t *testing.T) {
	repo := &mockBookingRepo{
		countToday:  1,
		findBooking: &domain.DailyBooking{StallID: 14, Status: domain.StatusBooked},
	}
	uc := NewBookingUsecase(repo)

	err := uc.BookStall(context.Background(), domain.BookStallInput{
		StallID: 13, VendorID: 2, MarketDate: time.Now(),
	})
	if !errors.Is(err, domain.ErrAlreadyHasStall) {
		t.Fatalf("expected ErrAlreadyHasStall, got %v", err)
	}
	if repo.tryBookCalled {
		t.Fatal("BookStallWithAudit should not be called when vendor already has a stall")
	}
}

func TestBookStall_StallNotAvailable(t *testing.T) {
	repo := &mockBookingRepo{countToday: 0, tryRows: 0}
	uc := NewBookingUsecase(repo)

	err := uc.BookStall(context.Background(), domain.BookStallInput{
		StallID: 13, VendorID: 2, MarketDate: time.Now(),
	})
	if !errors.Is(err, domain.ErrStallNotAvailable) {
		t.Fatalf("expected ErrStallNotAvailable, got %v", err)
	}
	if !repo.tryBookCalled {
		t.Fatal("expected BookStallWithAudit to be called")
	}
}

func TestBookStall_AlreadyHasStall_UniqueViolation(t *testing.T) {
	repo := &mockBookingRepo{countToday: 0, tryErr: domain.ErrAlreadyHasStall}
	uc := NewBookingUsecase(repo)

	err := uc.BookStall(context.Background(), domain.BookStallInput{
		StallID: 13, VendorID: 2, MarketDate: time.Now(),
	})
	if !errors.Is(err, domain.ErrAlreadyHasStall) {
		t.Fatalf("expected ErrAlreadyHasStall, got %v", err)
	}
	if !repo.tryBookCalled {
		t.Fatal("expected BookStallWithAudit to be called")
	}
}
