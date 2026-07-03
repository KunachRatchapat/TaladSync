package usecase

import (
	"context"
	"strings"
	"time"

	"taladsync/internal/domain"
)

type DirectoryRepository interface {
	SearchDirectory(ctx context.Context, date time.Time, q string) ([]domain.DirectoryEntry, error)
}

type DirectoryUsecase struct {
	bookings DirectoryRepository
}

func NewDirectoryUsecase(bookings DirectoryRepository) *DirectoryUsecase {
	return &DirectoryUsecase{bookings: bookings}
}

func (u *DirectoryUsecase) Search(ctx context.Context, date time.Time, q string) ([]domain.DirectoryEntry, error) {
	return u.bookings.SearchDirectory(ctx, date, strings.TrimSpace(q))
}
