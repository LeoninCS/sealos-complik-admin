package unban

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrUnbanInvalidInput = errors.New("user id and operator name are required")
	ErrUnbanNotFound     = errors.New("unban not found")
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

// CreateUnban creates a new unban record.
func (s *Service) CreateUnban(ctx context.Context, req CreateUnbanRequest) error {
	input, err := normalizeUnbanInput(req.UserID, req.OperatorName)
	if err != nil {
		return err
	}

	record := &Unban{
		UserID:       input.UserID,
		OperatorName: input.OperatorName,
	}

	if err := s.repository.CreateUnban(ctx, record); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// DeleteUnbans deletes all unban records for the given user.
func (s *Service) DeleteUnbans(ctx context.Context, userID uint64) error {
	if err := validateUserID(userID); err != nil {
		return err
	}

	if err := s.repository.DeleteUnbansByUserID(ctx, userID); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// GetUnbans returns all unban records for the given user.
func (s *Service) GetUnbans(ctx context.Context, userID uint64) ([]UnbanResponse, error) {
	if err := validateUserID(userID); err != nil {
		return nil, err
	}

	unbans, err := s.repository.GetUnbansByUserID(ctx, userID)
	if err != nil {
		return nil, translateRepositoryError(err)
	}

	responses := make([]UnbanResponse, 0, len(unbans))
	for i := range unbans {
		responses = append(responses, *toUnbanResponse(&unbans[i]))
	}

	return responses, nil
}

// ListUnbans returns all unban records.
func (s *Service) ListUnbans(ctx context.Context) ([]UnbanResponse, error) {
	unbans, err := s.repository.ListUnbans(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]UnbanResponse, 0, len(unbans))
	for i := range unbans {
		responses = append(responses, *toUnbanResponse(&unbans[i]))
	}

	return responses, nil
}

type normalizedUnbanInput struct {
	UserID       uint64
	OperatorName string
}

// normalizeUnbanInput keeps create validation consistent.
func normalizeUnbanInput(userID uint64, operatorName string) (*normalizedUnbanInput, error) {
	trimmedOperatorName := strings.TrimSpace(operatorName)

	if userID == 0 || trimmedOperatorName == "" {
		return nil, ErrUnbanInvalidInput
	}

	return &normalizedUnbanInput{
		UserID:       userID,
		OperatorName: trimmedOperatorName,
	}, nil
}

func validateUserID(userID uint64) error {
	if userID == 0 {
		return ErrUnbanInvalidInput
	}

	return nil
}

// translateRepositoryError hides storage details from the handler layer.
func translateRepositoryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUnbanNotFound
	}

	return err
}

func toUnbanResponse(record *Unban) *UnbanResponse {
	return &UnbanResponse{
		UserID:       record.UserID,
		OperatorName: record.OperatorName,
		CreatedAt:    record.CreatedAt,
		UpdatedAt:    record.UpdatedAt,
	}
}
