package projectconfig

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CreateProjectConfig creates a new project configuration in the database.
func (r *Repository) CreateProjectConfig(ctx context.Context, projectConfig *ProjectConfig) error {
	return r.db.WithContext(ctx).Create(projectConfig).Error
}

// GetProjectConfigByName returns a project configuration by its config name.
func (r *Repository) GetProjectConfigByName(ctx context.Context, configName string) (*ProjectConfig, error) {
	var projectConfig ProjectConfig
	if err := r.db.WithContext(ctx).Where("config_name = ?", configName).First(&projectConfig).Error; err != nil {
		return nil, err
	}

	return &projectConfig, nil
}

// ListProjectConfigs returns all project configurations.
func (r *Repository) ListProjectConfigs(ctx context.Context) ([]ProjectConfig, error) {
	var projectConfigs []ProjectConfig
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&projectConfigs).Error; err != nil {
		return nil, err
	}

	return projectConfigs, nil
}

// ListProjectConfigsByType returns project configurations filtered by config type.
func (r *Repository) ListProjectConfigsByType(ctx context.Context, configType string) ([]ProjectConfig, error) {
	var projectConfigs []ProjectConfig
	if err := r.db.WithContext(ctx).
		Where("config_type = ?", configType).
		Order("id ASC").
		Find(&projectConfigs).Error; err != nil {
		return nil, err
	}

	return projectConfigs, nil
}

// CountProjectConfigsByType counts configs for a type, optionally excluding one config name.
func (r *Repository) CountProjectConfigsByType(ctx context.Context, configType, excludeConfigName string) (int64, error) {
	query := r.db.WithContext(ctx).Model(&ProjectConfig{}).Where("config_type = ?", configType)
	if strings.TrimSpace(excludeConfigName) != "" {
		query = query.Where("config_name <> ?", strings.TrimSpace(excludeConfigName))
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateProjectConfig updates an existing project configuration in the database.
func (r *Repository) UpdateProjectConfig(ctx context.Context, projectConfig *ProjectConfig) error {
	return r.db.WithContext(ctx).Save(projectConfig).Error
}

// DeleteProjectConfigByName deletes a project configuration by its config name.
func (r *Repository) DeleteProjectConfigByName(ctx context.Context, configName string) error {
	result := r.db.WithContext(ctx).Where("config_name = ?", configName).Delete(&ProjectConfig{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
