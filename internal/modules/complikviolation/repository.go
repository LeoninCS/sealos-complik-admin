package complikviolation

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

func (r *Repository) CreateViolation(ctx context.Context, violation *ComplikViolationEvent) error {
	return r.db.WithContext(ctx).Create(violation).Error
}

func (r *Repository) GetViolationsByNamespace(ctx context.Context, namespace string, includeAll bool) ([]ComplikViolationEvent, error) {
	var violations []ComplikViolationEvent
	query := r.db.WithContext(ctx).Where("namespace = ?", namespace)
	if !includeAll {
		query = query.Where("is_illegal = ?", true)
	}

	if err := query.
		Order("detected_at DESC, id DESC").
		Find(&violations).Error; err != nil {
		return nil, err
	}
	if len(violations) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return violations, nil
}

func (r *Repository) ListViolations(ctx context.Context, includeAll bool) ([]ComplikViolationEvent, error) {
	var violations []ComplikViolationEvent
	query := r.db.WithContext(ctx)
	if !includeAll {
		query = query.Where("is_illegal = ?", true)
	}

	if err := query.Order("detected_at DESC, id DESC").Find(&violations).Error; err != nil {
		return nil, err
	}

	return violations, nil
}

func (r *Repository) DeleteViolationByID(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&ComplikViolationEvent{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *Repository) DeleteViolationsByNamespace(ctx context.Context, namespace string) error {
	result := r.db.WithContext(ctx).Where("namespace = ?", namespace).Delete(&ComplikViolationEvent{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *Repository) HasViolations(ctx context.Context, namespace string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&ComplikViolationEvent{}).
		Where("namespace = ?", namespace).
		Where("is_illegal = ?", true).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
