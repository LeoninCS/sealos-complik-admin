package violation

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Violation struct {
	// ID remains the internal primary key for joins and stable references.
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64    `gorm:"not null;index" json:"user_id"`
	ViolationType string    `gorm:"size:50;not null" json:"violation_type"`
	ViolationDesc string    `gorm:"type:text" json:"violation_desc,omitempty"`
	ViolationTime time.Time `gorm:"not null" json:"violation_time"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// AutoMigrate creates or updates the violations table schema.
func AutoMigrate(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("violation automigrate: database is nil")
	}

	if err := db.AutoMigrate(&Violation{}); err != nil {
		return fmt.Errorf("violation automigrate: %w", err)
	}

	return nil
}
