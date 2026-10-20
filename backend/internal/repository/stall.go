package repository

import (
	"context"

	"gorm.io/gorm"

	"taladsync/internal/domain"
)

type StallRepository struct {
	db *gorm.DB
}

func NewStallRepository(db *gorm.DB) *StallRepository {
	return &StallRepository{db: db}
}

func (r *StallRepository) ListZones(ctx context.Context) ([]string, error) {
	zones := []string{}
	err := r.db.WithContext(ctx).
		Model(&domain.Stall{}).
		Distinct().
		Order("zone ASC").
		Pluck("zone", &zones).Error
	if err != nil {
		return nil, err
	}
	return zones, nil
}
