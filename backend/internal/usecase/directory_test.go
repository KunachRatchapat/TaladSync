package usecase

import (
	"context"
	"testing"
	"time"

	"taladsync/internal/domain"
)

type mockDirectoryRepo struct {
	results []domain.DirectoryEntry
	err     error
	lastQ   string
}

func (m *mockDirectoryRepo) SearchDirectory(ctx context.Context, date time.Time, q string) ([]domain.DirectoryEntry, error) {
	m.lastQ = q
	if m.err != nil {
		return nil, m.err
	}
	return m.results, nil
}

func TestSearch_TrimsQuery(t *testing.T) {
	repo := &mockDirectoryRepo{results: []domain.DirectoryEntry{}}
	uc := NewDirectoryUsecase(repo)
	date := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)

	_, err := uc.Search(context.Background(), date, "  ของสด  ")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if repo.lastQ != "ของสด" {
		t.Fatalf("expected trimmed q, got %q", repo.lastQ)
	}
}

func TestSearch_EmptyQueryListsAll(t *testing.T) {
	expected := []domain.DirectoryEntry{
		{StallCode: "A-01", Zone: "Zone A", ShopName: "ร้านA-01", ProductType: "ของสด"},
		{StallCode: "A-02", Zone: "Zone A", ShopName: "ร้านA-02", ProductType: "ของสด"},
	}
	repo := &mockDirectoryRepo{results: expected}
	uc := NewDirectoryUsecase(repo)

	results, err := uc.Search(context.Background(), time.Now(), "")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if repo.lastQ != "" {
		t.Fatalf("expected empty q, got %q", repo.lastQ)
	}
}
