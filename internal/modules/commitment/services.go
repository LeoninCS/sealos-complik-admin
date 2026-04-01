package commitment

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrCommitmentAlreadyExists = errors.New("commitment already exists")
	ErrCommitmentInvalidInput  = errors.New("namespace, file name, and file url are required")
	ErrCommitmentNotFound      = errors.New("commitment not found")
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

// CreateCommitment creates a new commitment for the given namespace.
func (s *Service) CreateCommitment(ctx context.Context, req CreateCommitmentRequest) error {
	input, err := normalizeCreateCommitmentInput(req.Namespace, req.FileName, req.FileURL)
	if err != nil {
		return err
	}

	if _, err := s.repository.GetCommitmentByNamespace(ctx, input.Namespace); err == nil {
		return ErrCommitmentAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	commitment := &Commitment{
		Namespace: input.Namespace,
		FileName:  input.FileName,
		FileURL:   input.FileURL,
	}

	if err := s.repository.CreateCommitment(ctx, commitment); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// UpdateCommitment updates the commitment record for the given namespace.
func (s *Service) UpdateCommitment(ctx context.Context, namespace string, req UpdateCommitmentRequest) error {
	commitment, err := s.repository.GetCommitmentByNamespace(ctx, namespace)
	if err != nil {
		return translateRepositoryError(err)
	}

	input, err := normalizeUpdateCommitmentInput(namespace, req.FileName, req.FileURL)
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

// DeleteCommitment deletes the commitment record for the given namespace.
func (s *Service) DeleteCommitment(ctx context.Context, namespace string) error {
	if err := s.repository.DeleteCommitmentByNamespace(ctx, namespace); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

// GetCommitment returns the commitment record for the given namespace.
func (s *Service) GetCommitment(ctx context.Context, namespace string) (*CommitmentResponse, error) {
	commitment, err := s.repository.GetCommitmentByNamespace(ctx, namespace)
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
	Namespace string
	FileName  string
	FileURL   string
}

// normalizeCreateCommitmentInput keeps create validation consistent.
func normalizeCreateCommitmentInput(namespace, fileName, fileURL string) (*normalizedCommitmentInput, error) {
	return normalizeCommitmentInput(namespace, fileName, fileURL)
}

// normalizeUpdateCommitmentInput keeps update validation consistent.
func normalizeUpdateCommitmentInput(namespace, fileName, fileURL string) (*normalizedCommitmentInput, error) {
	return normalizeCommitmentInput(namespace, fileName, fileURL)
}

func normalizeCommitmentInput(namespace, fileName, fileURL string) (*normalizedCommitmentInput, error) {
	trimmedNamespace := strings.TrimSpace(namespace)
	trimmedFileName := strings.TrimSpace(fileName)
	trimmedFileURL := strings.TrimSpace(fileURL)

	if trimmedNamespace == "" || trimmedFileName == "" || trimmedFileURL == "" {
		return nil, ErrCommitmentInvalidInput
	}

	return &normalizedCommitmentInput{
		Namespace: trimmedNamespace,
		FileName:  trimmedFileName,
		FileURL:   trimmedFileURL,
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
		Namespace: commitment.Namespace,
		FileName:  commitment.FileName,
		FileURL:   commitment.FileURL,
		CreatedAt: commitment.CreatedAt,
		UpdatedAt: commitment.UpdatedAt,
	}
}
