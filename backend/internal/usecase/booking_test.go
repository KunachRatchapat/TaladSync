package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"taladsync/internal/domain"
)

type mockBookingRepo struct {
	listAvailable    []domain.DailyBooking
	listAvailableErr error
	tryRows          int64
	tryErr           error
	tryBookCalled    bool
	lastTryBookInput domain.BookStallInput
	findBooking      *domain.DailyBooking
	findBookingErr   error
	findCalled       bool
	lastFindVendorID uint
	lastAuditLog     domain.AuditLog
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

func (m *mockBookingRepo) FindVendorBookedToday(ctx context.Context, vendorID uint, date time.Time) (*domain.DailyBooking, error) {
	m.findCalled = true
	m.lastFindVendorID = vendorID
	return m.findBooking, m.findBookingErr
}

func TestBookStall_Success(t *testing.T) {
	repo := &mockBookingRepo{tryRows: 1}
	uc := NewBookingUsecase(repo)

	in := domain.BookStallInput{
		StallID: 13, VendorID: 2,
		MarketDate: time.Date(2026, 6, 29, 0, 0, 0, 0, time.UTC),
		ShopName:   "ร้านทด", ProductType: "ผลไม้",
	}

	if err := uc.BookStall(context.Background(), in); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if !repo.findCalled || !repo.tryBookCalled {
		t.Fatal("expected FindVendorBookedToday and BookStallWithAudit to be called")
	}
	if repo.lastFindVendorID != 2 {
		t.Fatalf("expected FindVendorBookedToday for vendor 2, got %d", repo.lastFindVendorID)
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
	repo := &mockBookingRepo{tryRows: 0}
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
	repo := &mockBookingRepo{tryErr: domain.ErrAlreadyHasStall}
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

func TestBookStall_FindError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mockBookingRepo{findBookingErr: dbErr}
	uc := NewBookingUsecase(repo)

	err := uc.BookStall(context.Background(), domain.BookStallInput{
		StallID: 13, VendorID: 2, MarketDate: time.Now(),
	})
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
	if repo.tryBookCalled {
		t.Fatal("BookStallWithAudit should not be called when the lookup fails")
	}
}
