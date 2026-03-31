package violation

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrViolationInvalidInput = errors.New("user id, violation type, and violation time are required")
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
	input, err := normalizeViolationInput(req.UserID, req.ViolationType, req.ViolationDesc, req.ViolationTime)
	if err != nil {
		return err
	}

	violation := &Violation{
		UserID:        input.UserID,
		ViolationType: input.ViolationType,
		ViolationDesc: input.ViolationDesc,
		ViolationTime: input.ViolationTime,
	}

	if err := s.repository.CreateViolation(ctx, violation); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// DeleteViolations deletes all violation records for the given user.
func (s *Service) DeleteViolations(ctx context.Context, userID uint64) error {
	if err := validateUserID(userID); err != nil {
		return err
	}

	if err := s.repository.DeleteViolationsByUserID(ctx, userID); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// GetViolations returns all violation records for the given user.
func (s *Service) GetViolations(ctx context.Context, userID uint64) ([]ViolationResponse, error) {
	if err := validateUserID(userID); err != nil {
		return nil, err
	}

	violations, err := s.repository.GetViolationsByUserID(ctx, userID)
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

// GetViolationStatus returns whether the given user has any violation records.
func (s *Service) GetViolationStatus(ctx context.Context, userID uint64) (*ViolationStatusResponse, error) {
	if err := validateUserID(userID); err != nil {
		return nil, err
	}

	violated, err := s.repository.HasViolations(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &ViolationStatusResponse{Violated: violated}, nil
}

type normalizedViolationInput struct {
	UserID        uint64
	ViolationType string
	ViolationDesc string
	ViolationTime time.Time
}

// normalizeViolationInput keeps create validation consistent.
func normalizeViolationInput(userID uint64, violationType, violationDesc string, violationTime time.Time) (*normalizedViolationInput, error) {
	trimmedViolationType := strings.TrimSpace(violationType)
	trimmedViolationDesc := strings.TrimSpace(violationDesc)

	if userID == 0 || trimmedViolationType == "" || violationTime.IsZero() {
		return nil, ErrViolationInvalidInput
	}

	return &normalizedViolationInput{
		UserID:        userID,
		ViolationType: trimmedViolationType,
		ViolationDesc: trimmedViolationDesc,
		ViolationTime: violationTime,
	}, nil
}

func validateUserID(userID uint64) error {
	if userID == 0 {
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
		UserID:        violation.UserID,
		ViolationType: violation.ViolationType,
		ViolationDesc: violation.ViolationDesc,
		ViolationTime: violation.ViolationTime,
		CreatedAt:     violation.CreatedAt,
		UpdatedAt:     violation.UpdatedAt,
	}
}
