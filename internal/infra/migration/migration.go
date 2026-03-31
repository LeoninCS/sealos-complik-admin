package migration

import (
	"fmt"

	"sealos-complik-admin/internal/modules/commitment"
	"sealos-complik-admin/internal/modules/projectconfig"
	"sealos-complik-admin/internal/modules/violation"

	"gorm.io/gorm"
)

// AutoMigrate runs all module migrations in one place.
func AutoMigrate(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("migration automigrate: database is nil")
	}

	migrations := []struct {
		name string
		run  func(*gorm.DB) error
	}{
		{name: "project config", run: projectconfig.AutoMigrate},
		{name: "commitment", run: commitment.AutoMigrate},
		{name: "violation", run: violation.AutoMigrate},
	}

	for _, migration := range migrations {
		if err := migration.run(db); err != nil {
			return fmt.Errorf("auto migrate %s: %w", migration.name, err)
		}
	}

	return nil
}
