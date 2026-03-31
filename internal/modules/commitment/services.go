package commitment

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrCommitmentAlreadyExists = errors.New("commitment already exists")
	ErrCommitmentInvalidInput  = errors.New("user id, file name, and file url are required")
	ErrCommitmentNotFound      = errors.New("commitment not found")
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

// CreateCommitment creates a new commitment for the given user.
func (s *Service) CreateCommitment(ctx context.Context, req CreateCommitmentRequest) error {
	input, err := normalizeCreateCommitmentInput(req.UserID, req.FileName, req.FileURL)
	if err != nil {
		return err
	}

	if _, err := s.repository.GetCommitmentByUserID(ctx, input.UserID); err == nil {
		return ErrCommitmentAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	commitment := &Commitment{
		UserID:   input.UserID,
		FileName: input.FileName,
		FileURL:  input.FileURL,
	}

	if err := s.repository.CreateCommitment(ctx, commitment); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// UpdateCommitment updates the commitment record for the given user.
func (s *Service) UpdateCommitment(ctx context.Context, userID uint64, req UpdateCommitmentRequest) error {
	commitment, err := s.repository.GetCommitmentByUserID(ctx, userID)
	if err != nil {
		return translateRepositoryError(err)
	}

	input, err := normalizeUpdateCommitmentInput(userID, req.FileName, req.FileURL)
	if err != nil {
		return err
	}

	commitment.FileName = input.FileName
	commitment.FileURL = input.FileURL

	if err := s.repository.UpdateCommitment(ctx, commitment); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// DeleteCommitment deletes the commitment record for the given user.
func (s *Service) DeleteCommitment(ctx context.Context, userID uint64) error {
	if err := s.repository.DeleteCommitmentByUserID(ctx, userID); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// GetCommitment returns the commitment record for the given user.
func (s *Service) GetCommitment(ctx context.Context, userID uint64) (*CommitmentResponse, error) {
	commitment, err := s.repository.GetCommitmentByUserID(ctx, userID)
	if err != nil {
		return nil, translateRepositoryError(err)
	}

	return toCommitmentResponse(commitment), nil
}

// ListCommitments returns all commitment records.
func (s *Service) ListCommitments(ctx context.Context) ([]CommitmentResponse, error) {
	commitments, err := s.repository.ListCommitments(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]CommitmentResponse, 0, len(commitments))
	for i := range commitments {
		responses = append(responses, *toCommitmentResponse(&commitments[i]))
	}

	return responses, nil
}

type normalizedCommitmentInput struct {
	UserID   uint64
	FileName string
	FileURL  string
}

// normalizeCreateCommitmentInput keeps create validation consistent.
func normalizeCreateCommitmentInput(userID uint64, fileName, fileURL string) (*normalizedCommitmentInput, error) {
	return normalizeCommitmentInput(userID, fileName, fileURL)
}

// normalizeUpdateCommitmentInput keeps update validation consistent.
func normalizeUpdateCommitmentInput(userID uint64, fileName, fileURL string) (*normalizedCommitmentInput, error) {
	return normalizeCommitmentInput(userID, fileName, fileURL)
}

func normalizeCommitmentInput(userID uint64, fileName, fileURL string) (*normalizedCommitmentInput, error) {
	trimmedFileName := strings.TrimSpace(fileName)
	trimmedFileURL := strings.TrimSpace(fileURL)

	if userID == 0 || trimmedFileName == "" || trimmedFileURL == "" {
		return nil, ErrCommitmentInvalidInput
	}

	return &normalizedCommitmentInput{
		UserID:   userID,
		FileName: trimmedFileName,
		FileURL:  trimmedFileURL,
	}, nil
}

// translateRepositoryError hides storage details from the handler layer.
func translateRepositoryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrCommitmentNotFound
	}

	return err
}

func toCommitmentResponse(commitment *Commitment) *CommitmentResponse {
	return &CommitmentResponse{
		UserID:    commitment.UserID,
		FileName:  commitment.FileName,
		FileURL:   commitment.FileURL,
		CreatedAt: commitment.CreatedAt,
		UpdatedAt: commitment.UpdatedAt,
	}
}
