package commitment

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Commitment struct {
	// ID remains the internal primary key for joins and stable references.
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"not null;index" json:"user_id"`
	FileName  string    `gorm:"size:255;not null" json:"file_name"`
	FileURL   string    `gorm:"size:512;not null" json:"file_url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// AutoMigrate creates or updates the commitments table schema.
func AutoMigrate(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("commitment automigrate: database is nil")
	}

	if err := db.AutoMigrate(&Commitment{}); err != nil {
		return fmt.Errorf("commitment automigrate: %w", err)
	}

	return nil
}
