package violation

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CreateViolation creates a new violation record.
func (r *Repository) CreateViolation(ctx context.Context, violation *Violation) error {
	return r.db.WithContext(ctx).Create(violation).Error
}

// GetViolationsByNamespace returns all violation records for the given namespace.
func (r *Repository) GetViolationsByNamespace(ctx context.Context, namespace string) ([]Violation, error) {
	var violations []Violation
	if err := r.db.WithContext(ctx).Where("namespace = ?", namespace).Order("violation_time DESC, id DESC").Find(&violations).Error; err != nil {
		return nil, err
	}
	if len(violations) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return violations, nil
}

// ListViolations returns all violation records.
func (r *Repository) ListViolations(ctx context.Context) ([]Violation, error) {
	var violations []Violation
	if err := r.db.WithContext(ctx).Order("violation_time DESC, id DESC").Find(&violations).Error; err != nil {
		return nil, err
	}

	return violations, nil
}

// DeleteViolationsByNamespace deletes all violation records for the given namespace.
func (r *Repository) DeleteViolationsByNamespace(ctx context.Context, namespace string) error {
	result := r.db.WithContext(ctx).Where("namespace = ?", namespace).Delete(&Violation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// HasViolations reports whether the given namespace has any violation records.
func (r *Repository) HasViolations(ctx context.Context, namespace string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&Violation{}).Where("namespace = ?", namespace).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
