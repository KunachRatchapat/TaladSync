package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"taladsync/internal/domain"
)

type mockStaffRepo struct {
	booking         *domain.DailyBooking
	findErr         error
	setRows         int64
	setErr          error
	lastAttendance  domain.AttendanceStatus
	lastStaffID     uint
	summary         domain.AttendanceSummary
	listChecklist   []domain.DailyBooking
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

func (m *mockStaffRepo) SetAttendance(ctx context.Context, bookingID, staffID uint, attendance domain.AttendanceStatus, checkedAt time.Time) (int64, error) {
	m.lastStaffID = staffID
	m.lastAttendance = attendance
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
	audit := &mockAuditWriter{}
	uc := NewStaffUsecase(repo, audit)

	_, err := uc.MarkAttendance(context.Background(), domain.MarkAttendanceInput{
		BookingID: 5, StaffID: 3, Present: true,
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if repo.lastAttendance != domain.AttendancePresent {
		t.Fatalf("expected PRESENT, got %s", repo.lastAttendance)
	}
	if len(audit.logs) != 1 || audit.logs[0].Action != domain.AuditCheckin {
		t.Fatalf("expected CHECKIN audit, got %+v", audit.logs)
	}
}

func TestMarkAttendance_NoShow(t *testing.T) {
	repo := &mockStaffRepo{
		booking: &domain.DailyBooking{
			ID: 5, Status: domain.StatusExpected, Attendance: domain.AttendancePending,
		},
		setRows: 1,
	}
	audit := &mockAuditWriter{}
	uc := NewStaffUsecase(repo, audit)

	_, err := uc.MarkAttendance(context.Background(), domain.MarkAttendanceInput{
		BookingID: 5, StaffID: 3, Present: false,
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if audit.logs[0].Action != domain.AuditNoshow {
		t.Fatalf("expected NOSHOW audit, got %s", audit.logs[0].Action)
	}
}

func TestMarkAttendance_Idempotent(t *testing.T) {
	repo := &mockStaffRepo{
		booking: &domain.DailyBooking{
			ID: 5, Status: domain.StatusBooked, Attendance: domain.AttendancePresent,
		},
	}
	audit := &mockAuditWriter{}
	uc := NewStaffUsecase(repo, audit)

	_, err := uc.MarkAttendance(context.Background(), domain.MarkAttendanceInput{
		BookingID: 5, StaffID: 3, Present: true,
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if repo.setRows != 0 && repo.lastAttendance != "" {
		t.Fatal("SetAttendance should not be called when already PRESENT")
	}
	if len(audit.logs) != 0 {
		t.Fatal("expected no audit on idempotent mark")
	}
}

func TestMarkAttendance_InvalidStatus(t *testing.T) {
	repo := &mockStaffRepo{
		booking: &domain.DailyBooking{
			ID: 5, Status: domain.StatusAvailable, Attendance: domain.AttendancePending,
		},
	}
	uc := NewStaffUsecase(repo, &mockAuditWriter{})

	_, err := uc.MarkAttendance(context.Background(), domain.MarkAttendanceInput{
		BookingID: 5, StaffID: 3, Present: true,
	})
	if !errors.Is(err, domain.ErrInvalidAttendanceTarget) {
		t.Fatalf("expected ErrInvalidAttendanceTarget, got %v", err)
	}
}

func TestMarkAttendance_NotFound(t *testing.T) {
	repo := &mockStaffRepo{findErr: domain.ErrBookingNotFound}
	uc := NewStaffUsecase(repo, &mockAuditWriter{})

	_, err := uc.MarkAttendance(context.Background(), domain.MarkAttendanceInput{
		BookingID: 99, StaffID: 3, Present: true,
	})
	if !errors.Is(err, domain.ErrBookingNotFound) {
		t.Fatalf("expected ErrBookingNotFound, got %v", err)
	}
}
