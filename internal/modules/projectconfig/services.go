package projectconfig

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrProjectConfigAlreadyExists = errors.New("project config already exists")
	ErrProjectConfigInvalidJSON   = errors.New("project config value must be valid json")
	ErrProjectConfigInvalidInput  = errors.New("config name and config type are required")
	ErrProjectConfigNotFound      = errors.New("project config not found")
	ErrProjectConfigTypeConflict  = errors.New("single type config already exists")
)

type Service struct {
	repository *Repository
}

var singleConfigTypes = map[string]struct{}{
	"model_runtime":                  {},
	"complick_notifications_runtime": {},
	"complik_notifications_runtime":  {},
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

// CreateProjectConfig creates a new project configuration based on the provided request data.
func (s *Service) CreateProjectConfig(ctx context.Context, req CreateProjectConfigRequest) error {
	input, err := normalizeProjectConfigInput(req.ConfigName, req.ConfigType, req.ConfigValue, req.Description)
	if err != nil {
		return err
	}
	if err := s.validateTypeCardinality(ctx, input.ConfigType, ""); err != nil {
		return err
	}
	projectConfig := &ProjectConfig{
		ConfigName:  input.ConfigName,
		ConfigType:  input.ConfigType,
		ConfigValue: input.ConfigValue,
		Description: input.Description,
	}

	// Attempt to create the project configuration in the database.
	if err := s.repository.CreateProjectConfig(ctx, projectConfig); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// UpdateProjectConfig updates an existing project configuration.
func (s *Service) UpdateProjectConfig(ctx context.Context, configName string, req UpdateProjectConfigRequest) error {
	projectConfig, err := s.repository.GetProjectConfigByName(ctx, strings.TrimSpace(configName))
	if err != nil {
		return translateRepositoryError(err)
	}

	input, err := normalizeProjectConfigInput(req.ConfigName, req.ConfigType, req.ConfigValue, req.Description)
	if err != nil {
		return err
	}
	if err := s.validateTypeCardinality(ctx, input.ConfigType, strings.TrimSpace(configName)); err != nil {
		return err
	}

	projectConfig.ConfigName = input.ConfigName
	projectConfig.ConfigType = input.ConfigType
	projectConfig.ConfigValue = input.ConfigValue
	projectConfig.Description = input.Description

	if err := s.repository.UpdateProjectConfig(ctx, projectConfig); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// DeleteProjectConfig deletes a project configuration by config name.
func (s *Service) DeleteProjectConfig(ctx context.Context, configName string) error {
	if err := s.repository.DeleteProjectConfigByName(ctx, strings.TrimSpace(configName)); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// GetProjectConfig returns a project configuration by config name.
func (s *Service) GetProjectConfig(ctx context.Context, configName string) (*ProjectConfigResponse, error) {
	projectConfig, err := s.repository.GetProjectConfigByName(ctx, strings.TrimSpace(configName))
	if err != nil {
		return nil, translateRepositoryError(err)
	}

	return toProjectConfigResponse(projectConfig), nil
}

// ListProjectConfigs returns all project configurations.
func (s *Service) ListProjectConfigs(ctx context.Context) ([]ProjectConfigResponse, error) {
	projectConfigs, err := s.repository.ListProjectConfigs(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]ProjectConfigResponse, 0, len(projectConfigs))
	for i := range projectConfigs {
		responses = append(responses, *toProjectConfigResponse(&projectConfigs[i]))
	}

	return responses, nil
}

// ListProjectConfigsByType returns project configurations filtered by config type.
func (s *Service) ListProjectConfigsByType(ctx context.Context, configType string) ([]ProjectConfigResponse, error) {
	trimmedType := strings.TrimSpace(configType)
	if trimmedType == "" {
		return nil, ErrProjectConfigInvalidInput
	}

	projectConfigs, err := s.repository.ListProjectConfigsByType(ctx, trimmedType)
	if err != nil {
		return nil, err
	}

	responses := make([]ProjectConfigResponse, 0, len(projectConfigs))
	for i := range projectConfigs {
		responses = append(responses, *toProjectConfigResponse(&projectConfigs[i]))
	}

	return responses, nil
}

type normalizedProjectConfigInput struct {
	ConfigName  string
	ConfigType  string
	ConfigValue json.RawMessage
	Description string
}

// normalizeProjectConfigInput keeps create/update validation consistent.
func normalizeProjectConfigInput(configName, configType string, configValue json.RawMessage, description string) (*normalizedProjectConfigInput, error) {
	trimmedConfigName := strings.TrimSpace(configName)
	trimmedConfigType := strings.TrimSpace(configType)
	trimmedDescription := strings.TrimSpace(description)

	if trimmedConfigName == "" || trimmedConfigType == "" {
		return nil, ErrProjectConfigInvalidInput
	}

	if !json.Valid(configValue) {
		return nil, ErrProjectConfigInvalidJSON
	}

	return &normalizedProjectConfigInput{
		ConfigName:  trimmedConfigName,
		ConfigType:  trimmedConfigType,
		ConfigValue: configValue,
		Description: trimmedDescription,
	}, nil
}

// validateTypeCardinality enforces single-entry types while allowing multi-entry types.
func (s *Service) validateTypeCardinality(ctx context.Context, configType, excludeConfigName string) error {
	normalizedType := strings.TrimSpace(configType)
	if _, isSingle := singleConfigTypes[normalizedType]; !isSingle {
		return nil
	}

	count, err := s.repository.CountProjectConfigsByType(ctx, normalizedType, excludeConfigName)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrProjectConfigTypeConflict
	}
	return nil
}

// translateRepositoryError hides storage details from the handler layer.
func translateRepositoryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrProjectConfigNotFound
	}

	var mysqlErr *mysqlDriver.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return ErrProjectConfigAlreadyExists
	}

	return err
}

// toProjectConfigResponse intentionally omits the internal ID from API output.
func toProjectConfigResponse(projectConfig *ProjectConfig) *ProjectConfigResponse {
	return &ProjectConfigResponse{
		ConfigName:  projectConfig.ConfigName,
		ConfigType:  projectConfig.ConfigType,
		ConfigValue: projectConfig.ConfigValue,
		Description: projectConfig.Description,
		CreatedAt:   projectConfig.CreatedAt,
		UpdatedAt:   projectConfig.UpdatedAt,
	}
}
