package unban

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Unban struct {
	// ID remains the internal primary key for joins and stable references.
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint64    `gorm:"not null;index" json:"user_id"`
	OperatorName string    `gorm:"size:100;not null" json:"operator_name"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// AutoMigrate creates or updates the unbans table schema.
func AutoMigrate(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("unban automigrate: database is nil")
	}

	if err := db.AutoMigrate(&Unban{}); err != nil {
		return fmt.Errorf("unban automigrate: %w", err)
	}

	return nil
}
