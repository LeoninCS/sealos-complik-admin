package ban

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrBanInvalidInput = errors.New("namespace, ban start time, and operator name are required")
	ErrBanNotFound     = errors.New("ban not found")
)

type Service struct {
	repository *Repository
	now        func() time.Time
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
		now:        time.Now,
	}
}

// CreateBan creates a new ban record.
func (s *Service) CreateBan(ctx context.Context, req CreateBanRequest) error {
	input, err := normalizeBanInput(req.Namespace, req.Reason, req.BanStartTime, req.BanEndTime, req.OperatorName)
	if err != nil {
		return err
	}

	ban := &Ban{
		Namespace:    input.Namespace,
		Reason:       input.Reason,
		BanStartTime: input.BanStartTime,
		BanEndTime:   input.BanEndTime,
		OperatorName: input.OperatorName,
	}

	if err := s.repository.CreateBan(ctx, ban); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// DeleteBans deletes all ban records for the given namespace.
func (s *Service) DeleteBans(ctx context.Context, namespace string) error {
	if err := validateNamespace(namespace); err != nil {
		return err
	}

	if err := s.repository.DeleteBansByNamespace(ctx, namespace); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// GetBans returns all ban records for the given namespace.
func (s *Service) GetBans(ctx context.Context, namespace string) ([]BanResponse, error) {
	if err := validateNamespace(namespace); err != nil {
		return nil, err
	}

	bans, err := s.repository.GetBansByNamespace(ctx, namespace)
	if err != nil {
		return nil, translateRepositoryError(err)
	}

	responses := make([]BanResponse, 0, len(bans))
	for i := range bans {
		responses = append(responses, *toBanResponse(&bans[i]))
	}

	return responses, nil
}

// ListBans returns all ban records.
func (s *Service) ListBans(ctx context.Context) ([]BanResponse, error) {
	bans, err := s.repository.ListBans(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]BanResponse, 0, len(bans))
	for i := range bans {
		responses = append(responses, *toBanResponse(&bans[i]))
	}

	return responses, nil
}

// GetBanStatus returns whether the given namespace is currently banned.
func (s *Service) GetBanStatus(ctx context.Context, namespace string) (*BanStatusResponse, error) {
	if err := validateNamespace(namespace); err != nil {
		return nil, err
	}

	banned, err := s.repository.HasActiveBan(ctx, namespace, s.now())
	if err != nil {
		return nil, err
	}

	return &BanStatusResponse{Banned: banned}, nil
}

type normalizedBanInput struct {
	Namespace    string
	Reason       string
	BanStartTime time.Time
	BanEndTime   *time.Time
	OperatorName string
}

// normalizeBanInput keeps create validation consistent.
func normalizeBanInput(namespace, reason string, banStartTime time.Time, banEndTime *time.Time, operatorName string) (*normalizedBanInput, error) {
	trimmedNamespace := strings.TrimSpace(namespace)
	trimmedReason := strings.TrimSpace(reason)
	trimmedOperatorName := strings.TrimSpace(operatorName)

	if trimmedNamespace == "" || banStartTime.IsZero() || trimmedOperatorName == "" {
		return nil, ErrBanInvalidInput
	}
	if banEndTime != nil && banEndTime.Before(banStartTime) {
		return nil, ErrBanInvalidInput
	}

	return &normalizedBanInput{
		Namespace:    trimmedNamespace,
		Reason:       trimmedReason,
		BanStartTime: banStartTime,
		BanEndTime:   banEndTime,
		OperatorName: trimmedOperatorName,
	}, nil
}

func validateNamespace(namespace string) error {
	if strings.TrimSpace(namespace) == "" {
		return ErrBanInvalidInput
	}

	return nil
}

// translateRepositoryError hides storage details from the handler layer.
func translateRepositoryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrBanNotFound
	}

	return err
}

func toBanResponse(ban *Ban) *BanResponse {
	return &BanResponse{
		Namespace:    ban.Namespace,
		Reason:       ban.Reason,
		BanStartTime: ban.BanStartTime,
		BanEndTime:   ban.BanEndTime,
		OperatorName: ban.OperatorName,
		CreatedAt:    ban.CreatedAt,
		UpdatedAt:    ban.UpdatedAt,
	}
}
