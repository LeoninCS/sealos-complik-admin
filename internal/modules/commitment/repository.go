package commitment

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

// CreateCommitment creates a new commitment record.
func (r *Repository) CreateCommitment(ctx context.Context, commitment *Commitment) error {
	return r.db.WithContext(ctx).Create(commitment).Error
}

// GetCommitmentByUserID returns a commitment by user ID.
func (r *Repository) GetCommitmentByUserID(ctx context.Context, userID uint64) (*Commitment, error) {
	var commitment Commitment
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&commitment).Error; err != nil {
		return nil, err
	}

	return &commitment, nil
}

// ListCommitments returns all commitment records.
func (r *Repository) ListCommitments(ctx context.Context) ([]Commitment, error) {
	var commitments []Commitment
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&commitments).Error; err != nil {
		return nil, err
	}

	return commitments, nil
}

// UpdateCommitment updates an existing commitment record.
func (r *Repository) UpdateCommitment(ctx context.Context, commitment *Commitment) error {
	return r.db.WithContext(ctx).Save(commitment).Error
}

// DeleteCommitmentByUserID deletes commitment records for the given user ID.
func (r *Repository) DeleteCommitmentByUserID(ctx context.Context, userID uint64) error {
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&Commitment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
