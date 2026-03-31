package ban

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Ban struct {
	// ID remains the internal primary key for joins and stable references.
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint64     `gorm:"not null;index" json:"user_id"`
	Reason       string     `gorm:"size:500" json:"reason,omitempty"`
	BanStartTime time.Time  `gorm:"not null" json:"ban_start_time"`
	BanEndTime   *time.Time `json:"ban_end_time,omitempty"`
	OperatorName string     `gorm:"size:100;not null" json:"operator_name"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// AutoMigrate creates or updates the bans table schema.
func AutoMigrate(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("ban automigrate: database is nil")
	}

	if err := db.AutoMigrate(&Ban{}); err != nil {
		return fmt.Errorf("ban automigrate: %w", err)
	}

	return nil
}
