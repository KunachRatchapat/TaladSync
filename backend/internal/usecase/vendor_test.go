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
}

func (m *mockVendorRepo) GetByStallAndDate(ctx context.Context, stallID uint, date time.Time) (*domain.DailyBooking, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.booking, nil
}

func (m *mockVendorRepo) TryLeaveStall(ctx context.Context, in domain.LeaveStallInput) (int64, error) {
	m.tryLeaveInput = in
	return m.tryRows, m.tryErr
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
	audit := &mockAuditWriter{}
	uc := NewVendorUsecase(repo, audit)

	err := uc.Leave(context.Background(), domain.LeaveStallInput{
		StallID: 1, VendorID: 1, MarketDate: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if len(audit.logs) != 1 || audit.logs[0].Action != domain.AuditLeave {
		t.Fatalf("expected LEAVE audit, got %+v", audit.logs)
	}
}

func TestLeave_NotOwnStall(t *testing.T) {
	ownerID := uint(1)
	repo := &mockVendorRepo{
		booking: &domain.DailyBooking{
			Status: domain.StatusExpected, VendorID: &ownerID,
		},
	}
	uc := NewVendorUsecase(repo, &mockAuditWriter{})

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
	uc := NewVendorUsecase(repo, &mockAuditWriter{})

	err := uc.Leave(context.Background(), domain.LeaveStallInput{
		StallID: 1, VendorID: 1, MarketDate: time.Now(),
	})
	if !errors.Is(err, domain.ErrLeaveNotAllowed) {
		t.Fatalf("expected ErrLeaveNotAllowed, got %v", err)
	}
}
