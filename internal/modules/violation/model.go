package violation

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

const legacyNamespaceColumn = "user" + "_id"

type Violation struct {
	// ID remains the internal primary key for joins and stable references.
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Namespace     string    `gorm:"size:255;not null;index" json:"namespace"`
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

	if db.Migrator().HasColumn(&Violation{}, legacyNamespaceColumn) && !db.Migrator().HasColumn(&Violation{}, "namespace") {
		if err := db.Migrator().RenameColumn(&Violation{}, legacyNamespaceColumn, "namespace"); err != nil {
			return fmt.Errorf("violation rename legacy namespace column: %w", err)
		}
	}

	if err := db.AutoMigrate(&Violation{}); err != nil {
		return fmt.Errorf("violation automigrate: %w", err)
	}

	return nil
}
