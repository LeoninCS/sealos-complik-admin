package ban

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrBanInvalidInput = errors.New("user id, ban start time, and operator name are required")
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
	input, err := normalizeBanInput(req.UserID, req.Reason, req.BanStartTime, req.BanEndTime, req.OperatorName)
	if err != nil {
		return err
	}

	ban := &Ban{
		UserID:       input.UserID,
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

// DeleteBans deletes all ban records for the given user.
func (s *Service) DeleteBans(ctx context.Context, userID uint64) error {
	if err := validateUserID(userID); err != nil {
		return err
	}

	if err := s.repository.DeleteBansByUserID(ctx, userID); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// GetBans returns all ban records for the given user.
func (s *Service) GetBans(ctx context.Context, userID uint64) ([]BanResponse, error) {
	if err := validateUserID(userID); err != nil {
		return nil, err
	}

	bans, err := s.repository.GetBansByUserID(ctx, userID)
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

// GetBanStatus returns whether the given user is currently banned.
func (s *Service) GetBanStatus(ctx context.Context, userID uint64) (*BanStatusResponse, error) {
	if err := validateUserID(userID); err != nil {
		return nil, err
	}

	banned, err := s.repository.HasActiveBan(ctx, userID, s.now())
	if err != nil {
		return nil, err
	}

	return &BanStatusResponse{Banned: banned}, nil
}

type normalizedBanInput struct {
	UserID       uint64
	Reason       string
	BanStartTime time.Time
	BanEndTime   *time.Time
	OperatorName string
}

// normalizeBanInput keeps create validation consistent.
func normalizeBanInput(userID uint64, reason string, banStartTime time.Time, banEndTime *time.Time, operatorName string) (*normalizedBanInput, error) {
	trimmedReason := strings.TrimSpace(reason)
	trimmedOperatorName := strings.TrimSpace(operatorName)

	if userID == 0 || banStartTime.IsZero() || trimmedOperatorName == "" {
		return nil, ErrBanInvalidInput
	}
	if banEndTime != nil && banEndTime.Before(banStartTime) {
		return nil, ErrBanInvalidInput
	}

	return &normalizedBanInput{
		UserID:       userID,
		Reason:       trimmedReason,
		BanStartTime: banStartTime,
		BanEndTime:   banEndTime,
		OperatorName: trimmedOperatorName,
	}, nil
}

func validateUserID(userID uint64) error {
	if userID == 0 {
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
		UserID:       ban.UserID,
		Reason:       ban.Reason,
		BanStartTime: ban.BanStartTime,
		BanEndTime:   ban.BanEndTime,
		OperatorName: ban.OperatorName,
		CreatedAt:    ban.CreatedAt,
		UpdatedAt:    ban.UpdatedAt,
	}
}
