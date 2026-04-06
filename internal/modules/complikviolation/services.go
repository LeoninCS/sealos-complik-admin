package complikviolation

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrViolationInvalidInput = errors.New("namespace, detector name, and detected time are required")
	ErrViolationNotFound     = errors.New("complik violation not found")
)

const defaultStatus = "open"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateViolation(ctx context.Context, req CreateViolationRequest) error {
	input, err := normalizeViolationInput(req)
	if err != nil {
		return err
	}

	pathJSON, err := marshalStringSlice(input.Path)
	if err != nil {
		return err
	}
	keywordsJSON, err := marshalStringSlice(input.Keywords)
	if err != nil {
		return err
	}
	rawPayloadJSON, err := marshalRawPayload(input.RawPayload)
	if err != nil {
		return err
	}

	violation := &ComplikViolationEvent{
		Namespace:     input.Namespace,
		Region:        input.Region,
		DiscoveryName: input.DiscoveryName,
		CollectorName: input.CollectorName,
		DetectorName:  input.DetectorName,
		ResourceName:  input.ResourceName,
		Host:          input.Host,
		URL:           input.URL,
		Path:          pathJSON,
		Keywords:      keywordsJSON,
		Description:   input.Description,
		Explanation:   input.Explanation,
		IsIllegal:     input.IsIllegal,
		IsTest:        input.IsTest,
		Status:        input.Status,
		DetectedAt:    input.DetectedAt,
		RawPayload:    rawPayloadJSON,
	}

	if err := s.repository.CreateViolation(ctx, violation); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

func (s *Service) DeleteViolations(ctx context.Context, namespace string) error {
	if err := validateNamespace(namespace); err != nil {
		return err
	}

	if err := s.repository.DeleteViolationsByNamespace(ctx, namespace); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

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

func (s *Service) UpdateViolationStatus(ctx context.Context, id uint64, status string) error {
	trimmedStatus := strings.TrimSpace(status)
	if id == 0 || (trimmedStatus != "open" && trimmedStatus != "reviewing" && trimmedStatus != "closed") {
		return ErrViolationInvalidInput
	}

	if err := s.repository.UpdateViolationStatus(ctx, id, trimmedStatus); err != nil {
		return translateRepositoryError(err)
	}

	return nil
}

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
	Region        string
	DiscoveryName string
	CollectorName string
	DetectorName  string
	ResourceName  string
	Host          string
	URL           string
	Path          []string
	Keywords      []string
	Description   string
	Explanation   string
	IsIllegal     bool
	IsTest        bool
	Status        string
	DetectedAt    time.Time
	RawPayload    json.RawMessage
}

func normalizeViolationInput(req CreateViolationRequest) (*normalizedViolationInput, error) {
	trimmedNamespace := strings.TrimSpace(req.Namespace)
	trimmedDetectorName := strings.TrimSpace(req.DetectorName)
	trimmedStatus := strings.TrimSpace(req.Status)

	if trimmedNamespace == "" || trimmedDetectorName == "" || req.DetectedAt.IsZero() {
		return nil, ErrViolationInvalidInput
	}
	if trimmedStatus == "" {
		trimmedStatus = defaultStatus
	}

	isIllegal := true
	if req.IsIllegal != nil {
		isIllegal = *req.IsIllegal
	}

	return &normalizedViolationInput{
		Namespace:     trimmedNamespace,
		Region:        strings.TrimSpace(req.Region),
		DiscoveryName: strings.TrimSpace(req.DiscoveryName),
		CollectorName: strings.TrimSpace(req.CollectorName),
		DetectorName:  trimmedDetectorName,
		ResourceName:  strings.TrimSpace(req.ResourceName),
		Host:          strings.TrimSpace(req.Host),
		URL:           strings.TrimSpace(req.URL),
		Path:          req.Path,
		Keywords:      req.Keywords,
		Description:   strings.TrimSpace(req.Description),
		Explanation:   strings.TrimSpace(req.Explanation),
		IsIllegal:     isIllegal,
		IsTest:        req.IsTest,
		Status:        trimmedStatus,
		DetectedAt:    req.DetectedAt,
		RawPayload:    req.RawPayload,
	}, nil
}

func validateNamespace(namespace string) error {
	if strings.TrimSpace(namespace) == "" {
		return ErrViolationInvalidInput
	}

	return nil
}

func translateRepositoryError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrViolationNotFound
	}

	return err
}

func marshalStringSlice(values []string) (*string, error) {
	if len(values) == 0 {
		return nil, nil
	}

	data, err := json.Marshal(values)
	if err != nil {
		return nil, ErrViolationInvalidInput
	}

	result := string(data)
	return &result, nil
}

func marshalRawPayload(payload json.RawMessage) (*string, error) {
	if len(payload) == 0 {
		return nil, nil
	}
	if !json.Valid(payload) {
		return nil, ErrViolationInvalidInput
	}

	result := string(payload)
	return &result, nil
}

func parseStringSlice(raw *string) []string {
	if raw == nil || *raw == "" {
		return nil
	}

	var values []string
	if err := json.Unmarshal([]byte(*raw), &values); err != nil {
		return nil
	}

	return values
}

func parseRawPayload(raw *string) json.RawMessage {
	if raw == nil || *raw == "" {
		return nil
	}

	return json.RawMessage(*raw)
}

func toViolationResponse(violation *ComplikViolationEvent) *ViolationResponse {
	return &ViolationResponse{
		ID:            violation.ID,
		Namespace:     violation.Namespace,
		Region:        violation.Region,
		DiscoveryName: violation.DiscoveryName,
		CollectorName: violation.CollectorName,
		DetectorName:  violation.DetectorName,
		ResourceName:  violation.ResourceName,
		Host:          violation.Host,
		URL:           violation.URL,
		Path:          parseStringSlice(violation.Path),
		Keywords:      parseStringSlice(violation.Keywords),
		Description:   violation.Description,
		Explanation:   violation.Explanation,
		IsIllegal:     violation.IsIllegal,
		IsTest:        violation.IsTest,
		Status:        violation.Status,
		DetectedAt:    violation.DetectedAt,
		RawPayload:    parseRawPayload(violation.RawPayload),
		CreatedAt:     violation.CreatedAt,
		UpdatedAt:     violation.UpdatedAt,
	}
}
