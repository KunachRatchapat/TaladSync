package repository

import (
	"context"

	"gorm.io/gorm"

	"taladsync/internal/domain"
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(ctx context.Context, log domain.AuditLog) error {
	return r.db.WithContext(ctx).Create(&log).Error
}
