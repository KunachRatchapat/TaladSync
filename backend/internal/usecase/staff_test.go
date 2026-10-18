package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"taladsync/internal/domain"
)

type mockStaffRepo struct {
	booking        *domain.DailyBooking
	findErr        error
	setRows        int64
	setErr         error
	setCalled      bool
	lastAttendance domain.AttendanceStatus
	lastStaffID    uint
	lastAuditLog   domain.AuditLog
	summary        domain.AttendanceSummary
	listChecklist  []domain.DailyBooking
}

func (m *mockStaffRepo) FindByID(ctx context.Context, id uint) (*domain.DailyBooking, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.booking, nil
}

func (m *mockStaffRepo) ListChecklist(ctx context.Context, date time.Time) ([]domain.DailyBooking, error) {
	return m.listChecklist, nil
}

func (m *mockStaffRepo) SetAttendanceWithAudit(ctx context.Context, bookingID, staffID uint, attendance domain.AttendanceStatus, checkedAt time.Time, log domain.AuditLog) (int64, error) {
	m.setCalled = true
	m.lastStaffID = staffID
	m.lastAttendance = attendance
	m.lastAuditLog = log
	return m.setRows, m.setErr
}

func (m *mockStaffRepo) AttendanceSummary(ctx context.Context, date time.Time) (domain.AttendanceSummary, error) {
	return m.summary, nil
}

func TestMarkAttendance_Present(t *testing.T) {
	repo := &mockStaffRepo{
		booking: &domain.DailyBooking{
			ID: 5, Status: domain.StatusBooked, Attendance: domain.AttendancePending,
		},
		setRows: 1,
	}
	uc := NewStaffUsecase(repo)

	_, err := uc.MarkAttendance(context.Background(), domain.MarkAttendanceInput{
		BookingID: 5, StaffID: 3, Present: true,
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if repo.lastAttendance != domain.AttendancePresent {
		t.Fatalf("expected PRESENT, got %s", repo.lastAttendance)
	}
	if repo.lastAuditLog.Action != domain.AuditCheckin {
		t.Fatalf("expected CHECKIN audit, got %q", repo.lastAuditLog.Action)
	}
	if repo.lastAuditLog.ActorUserID != 3 {
		t.Fatalf("expected actor=3, got %d", repo.lastAuditLog.ActorUserID)
	}
	if repo.lastAuditLog.DailyBookingID == nil || *repo.lastAuditLog.DailyBookingID != 5 {
		t.Fatalf("expected DailyBookingID=5, got %v", repo.lastAuditLog.DailyBookingID)
	}
}

func TestMarkAttendance_NoShow(t *testing.T) {
	repo := &mockStaffRepo{
		booking: &domain.DailyBooking{
			ID: 5, Status: domain.StatusExpected, Attendance: domain.AttendancePending,
		},
		setRows: 1,
	}
	uc := NewStaffUsecase(repo)

	_, err := uc.MarkAttendance(context.Background(), domain.MarkAttendanceInput{
		BookingID: 5, StaffID: 3, Present: false,
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if repo.lastAuditLog.Action != domain.AuditNoshow {
		t.Fatalf("expected NOSHOW audit, got %q", repo.lastAuditLog.Action)
	}
}

func TestMarkAttendance_Idempotent(t *testing.T) {
	repo := &mockStaffRepo{
		booking: &domain.DailyBooking{
			ID: 5, Status: domain.StatusBooked, Attendance: domain.AttendancePresent,
		},
	}
	uc := NewStaffUsecase(repo)

	_, err := uc.MarkAttendance(context.Background(), domain.MarkAttendanceInput{
		BookingID: 5, StaffID: 3, Present: true,
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if repo.setCalled {
		t.Fatal("SetAttendanceWithAudit should not be called when already PRESENT")
	}
}

func TestMarkAttendance_InvalidStatus(t *testing.T) {
	repo := &mockStaffRepo{
		booking: &domain.DailyBooking{
			ID: 5, Status: domain.StatusAvailable, Attendance: domain.AttendancePending,
		},
	}
	uc := NewStaffUsecase(repo)

	_, err := uc.MarkAttendance(context.Background(), domain.MarkAttendanceInput{
		BookingID: 5, StaffID: 3, Present: true,
	})
	if !errors.Is(err, domain.ErrInvalidAttendanceTarget) {
		t.Fatalf("expected ErrInvalidAttendanceTarget, got %v", err)
	}
}

func TestMarkAttendance_NotFound(t *testing.T) {
	repo := &mockStaffRepo{findErr: domain.ErrBookingNotFound}
	uc := NewStaffUsecase(repo)

	_, err := uc.MarkAttendance(context.Background(), domain.MarkAttendanceInput{
		BookingID: 99, StaffID: 3, Present: true,
	})
	if !errors.Is(err, domain.ErrBookingNotFound) {
		t.Fatalf("expected ErrBookingNotFound, got %v", err)
	}
}
