package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"taladsync/internal/domain"
)

type mockVendorRepo struct {
	booking       *domain.DailyBooking
	getErr        error
	tryRows       int64
	tryErr        error
	tryLeaveInput domain.LeaveStallInput
	tryCalled     bool
	lastAuditLog  domain.AuditLog
	owned         []domain.DailyBooking
	ownedErr      error
}

func (m *mockVendorRepo) GetByStallAndDate(ctx context.Context, stallID uint, date time.Time) (*domain.DailyBooking, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.booking, nil
}

func (m *mockVendorRepo) LeaveStallWithAudit(ctx context.Context, in domain.LeaveStallInput, log domain.AuditLog) (int64, error) {
	m.tryCalled = true
	m.tryLeaveInput = in
	m.lastAuditLog = log
	return m.tryRows, m.tryErr
}

func (m *mockVendorRepo) ListOwnedStallBookings(ctx context.Context, vendorID uint, date time.Time) ([]domain.DailyBooking, error) {
	return m.owned, m.ownedErr
}

type mockAuditWriter struct {
	logs    []domain.AuditLog
	createErr error
}

func (m *mockAuditWriter) Create(ctx context.Context, log domain.AuditLog) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.logs = append(m.logs, log)
	return nil
}

func TestGetMyStalls(t *testing.T) {
	repo := &mockVendorRepo{
		owned: []domain.DailyBooking{
			{StallID: 1, Status: domain.StatusExpected, Stall: &domain.Stall{Code: "A-01"}},
			{StallID: 2, Status: domain.StatusAvailable, Stall: &domain.Stall{Code: "A-02"}},
		},
	}
	uc := NewVendorUsecase(repo)

	stalls, err := uc.GetMyStalls(context.Background(), 1, time.Now())
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if len(stalls) != 2 {
		t.Fatalf("expected 2 stalls, got %d", len(stalls))
	}
}

func TestLeave_Success(t *testing.T) {
	vendorID := uint(1)
	repo := &mockVendorRepo{
		booking: &domain.DailyBooking{
			ID:      10,
			StallID: 1,
			Status:  domain.StatusExpected,
			VendorID: &vendorID,
			Stall:   &domain.Stall{Code: "A-01"},
		},
		tryRows: 1,
	}
	uc := NewVendorUsecase(repo)

	err := uc.Leave(context.Background(), domain.LeaveStallInput{
		StallID: 1, VendorID: 1, MarketDate: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if repo.lastAuditLog.Action != domain.AuditLeave {
		t.Fatalf("expected LEAVE audit, got %q", repo.lastAuditLog.Action)
	}
	if repo.lastAuditLog.DailyBookingID == nil || *repo.lastAuditLog.DailyBookingID != 10 {
		t.Fatalf("expected DailyBookingID=10, got %v", repo.lastAuditLog.DailyBookingID)
	}
	if repo.lastAuditLog.Detail != "stall=A-01" {
		t.Fatalf("expected detail stall=A-01, got %q", repo.lastAuditLog.Detail)
	}
}

func TestLeave_RaceLost(t *testing.T) {
	vendorID := uint(1)
	repo := &mockVendorRepo{
		booking: &domain.DailyBooking{
			ID: 10, Status: domain.StatusExpected, VendorID: &vendorID,
		},
		tryRows: 0,
	}
	uc := NewVendorUsecase(repo)

	err := uc.Leave(context.Background(), domain.LeaveStallInput{
		StallID: 1, VendorID: 1, MarketDate: time.Now(),
	})
	if !errors.Is(err, domain.ErrLeaveNotAllowed) {
		t.Fatalf("expected ErrLeaveNotAllowed, got %v", err)
	}
	if !repo.tryCalled {
		t.Fatal("expected LeaveStallWithAudit to be called")
	}
}

func TestLeave_NotOwnStall(t *testing.T) {
	ownerID := uint(1)
	repo := &mockVendorRepo{
		booking: &domain.DailyBooking{
			Status: domain.StatusExpected, VendorID: &ownerID,
		},
	}
	uc := NewVendorUsecase(repo)

	err := uc.Leave(context.Background(), domain.LeaveStallInput{
		StallID: 1, VendorID: 99, MarketDate: time.Now(),
	})
	if !errors.Is(err, domain.ErrNotOwnStall) {
		t.Fatalf("expected ErrNotOwnStall, got %v", err)
	}
}

func TestLeave_WrongStatus(t *testing.T) {
	vendorID := uint(1)
	repo := &mockVendorRepo{
		booking: &domain.DailyBooking{
			Status: domain.StatusAvailable, VendorID: &vendorID,
		},
	}
	uc := NewVendorUsecase(repo)

	err := uc.Leave(context.Background(), domain.LeaveStallInput{
		StallID: 1, VendorID: 1, MarketDate: time.Now(),
	})
	if !errors.Is(err, domain.ErrLeaveNotAllowed) {
		t.Fatalf("expected ErrLeaveNotAllowed, got %v", err)
	}
}
