package violation

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrViolationInvalidInput = errors.New("namespace, violation type, and violation time are required")
	ErrViolationNotFound     = errors.New("violation not found")
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

// CreateViolation creates a new violation record.
func (s *Service) CreateViolation(ctx context.Context, req CreateViolationRequest) error {
	input, err := normalizeViolationInput(req.Namespace, req.ViolationType, req.ViolationDesc, req.ViolationTime)
	if err != nil {
		return err
	}

	violation := &Violation{
		Namespace:     input.Namespace,
		ViolationType: input.ViolationType,
		ViolationDesc: input.ViolationDesc,
		ViolationTime: input.ViolationTime,
	}

	if err := s.repository.CreateViolation(ctx, violation); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// DeleteViolations deletes all violation records for the given namespace.
func (s *Service) DeleteViolations(ctx context.Context, namespace string) error {
	if err := validateNamespace(namespace); err != nil {
		return err
	}

	if err := s.repository.DeleteViolationsByNamespace(ctx, namespace); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// GetViolations returns all violation records for the given namespace.
func (s *Service) GetViolations(ctx context.Context, namespace string) ([]ViolationResponse, error) {
	if err := validateNamespace(namespace); err != nil {
		return nil, err
	}

	violations, err := s.repository.GetViolationsByNamespace(ctx, namespace)
	if err != nil {
		return nil, translateRepositoryError(err)
	}

	responses := make([]ViolationResponse, 0, len(violations))
	for i := range violations {
		responses = append(responses, *toViolationResponse(&violations[i]))
	}

	return responses, nil
}

// ListViolations returns all violation records.
func (s *Service) ListViolations(ctx context.Context) ([]ViolationResponse, error) {
	violations, err := s.repository.ListViolations(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]ViolationResponse, 0, len(violations))
	for i := range violations {
		responses = append(responses, *toViolationResponse(&violations[i]))
	}

	return responses, nil
}

// GetViolationStatus returns whether the given namespace has any violation records.
func (s *Service) GetViolationStatus(ctx context.Context, namespace string) (*ViolationStatusResponse, error) {
	if err := validateNamespace(namespace); err != nil {
		return nil, err
	}

	violated, err := s.repository.HasViolations(ctx, namespace)
	if err != nil {
		return nil, err
	}

	return &ViolationStatusResponse{Violated: violated}, nil
}

type normalizedViolationInput struct {
	Namespace     string
	ViolationType string
	ViolationDesc string
	ViolationTime time.Time
}

// normalizeViolationInput keeps create validation consistent.
func normalizeViolationInput(namespace, violationType, violationDesc string, violationTime time.Time) (*normalizedViolationInput, error) {
	trimmedNamespace := strings.TrimSpace(namespace)
	trimmedViolationType := strings.TrimSpace(violationType)
	trimmedViolationDesc := strings.TrimSpace(violationDesc)

	if trimmedNamespace == "" || trimmedViolationType == "" || violationTime.IsZero() {
		return nil, ErrViolationInvalidInput
	}

	return &normalizedViolationInput{
		Namespace:     trimmedNamespace,
		ViolationType: trimmedViolationType,
		ViolationDesc: trimmedViolationDesc,
		ViolationTime: violationTime,
	}, nil
}

func validateNamespace(namespace string) error {
	if strings.TrimSpace(namespace) == "" {
		return ErrViolationInvalidInput
	}

	return nil
}

// translateRepositoryError hides storage details from the handler layer.
func translateRepositoryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrViolationNotFound
	}

	return err
}

func toViolationResponse(violation *Violation) *ViolationResponse {
	return &ViolationResponse{
		Namespace:     violation.Namespace,
		ViolationType: violation.ViolationType,
		ViolationDesc: violation.ViolationDesc,
		ViolationTime: violation.ViolationTime,
		CreatedAt:     violation.CreatedAt,
		UpdatedAt:     violation.UpdatedAt,
	}
}
