package procscanviolation

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrViolationInvalidInput = errors.New("namespace, pid, process name, process command, message, and detected time are required")
	ErrViolationNotFound     = errors.New("procscan violation not found")
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

	rawPayloadJSON, err := marshalRawPayload(input.RawPayload)
	if err != nil {
		return err
	}

	violation := &ProcscanViolationEvent{
		Namespace:         input.Namespace,
		PodName:           input.PodName,
		ContainerID:       input.ContainerID,
		NodeName:          input.NodeName,
		PID:               input.PID,
		ProcessName:       input.ProcessName,
		ProcessCommand:    input.ProcessCommand,
		MatchType:         input.MatchType,
		MatchRule:         input.MatchRule,
		Message:           input.Message,
		LabelActionStatus: input.LabelActionStatus,
		LabelActionResult: input.LabelActionResult,
		Status:            input.Status,
		DetectedAt:        input.DetectedAt,
		RawPayload:        rawPayloadJSON,
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
	Namespace         string
	PodName           string
	ContainerID       string
	NodeName          string
	PID               int
	ProcessName       string
	ProcessCommand    string
	MatchType         string
	MatchRule         string
	Message           string
	LabelActionStatus string
	LabelActionResult string
	Status            string
	DetectedAt        time.Time
	RawPayload        json.RawMessage
}

func normalizeViolationInput(req CreateViolationRequest) (*normalizedViolationInput, error) {
	trimmedNamespace := strings.TrimSpace(req.Namespace)
	trimmedProcessName := strings.TrimSpace(req.ProcessName)
	trimmedProcessCommand := strings.TrimSpace(req.ProcessCommand)
	trimmedMessage := strings.TrimSpace(req.Message)
	trimmedStatus := strings.TrimSpace(req.Status)

	if trimmedNamespace == "" || req.PID <= 0 || trimmedProcessName == "" || trimmedProcessCommand == "" || trimmedMessage == "" || req.DetectedAt.IsZero() {
		return nil, ErrViolationInvalidInput
	}
	if trimmedStatus == "" {
		trimmedStatus = defaultStatus
	}

	return &normalizedViolationInput{
		Namespace:         trimmedNamespace,
		PodName:           strings.TrimSpace(req.PodName),
		ContainerID:       strings.TrimSpace(req.ContainerID),
		NodeName:          strings.TrimSpace(req.NodeName),
		PID:               req.PID,
		ProcessName:       trimmedProcessName,
		ProcessCommand:    trimmedProcessCommand,
		MatchType:         strings.TrimSpace(req.MatchType),
		MatchRule:         strings.TrimSpace(req.MatchRule),
		Message:           trimmedMessage,
		LabelActionStatus: strings.TrimSpace(req.LabelActionStatus),
		LabelActionResult: strings.TrimSpace(req.LabelActionResult),
		Status:            trimmedStatus,
		DetectedAt:        req.DetectedAt,
		RawPayload:        req.RawPayload,
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

func parseRawPayload(raw *string) json.RawMessage {
	if raw == nil || *raw == "" {
		return nil
	}

	return json.RawMessage(*raw)
}

func toViolationResponse(violation *ProcscanViolationEvent) *ViolationResponse {
	return &ViolationResponse{
		ID:                violation.ID,
		Namespace:         violation.Namespace,
		PodName:           violation.PodName,
		ContainerID:       violation.ContainerID,
		NodeName:          violation.NodeName,
		PID:               violation.PID,
		ProcessName:       violation.ProcessName,
		ProcessCommand:    violation.ProcessCommand,
		MatchType:         violation.MatchType,
		MatchRule:         violation.MatchRule,
		Message:           violation.Message,
		LabelActionStatus: violation.LabelActionStatus,
		LabelActionResult: violation.LabelActionResult,
		Status:            violation.Status,
		DetectedAt:        violation.DetectedAt,
		RawPayload:        parseRawPayload(violation.RawPayload),
		CreatedAt:         violation.CreatedAt,
		UpdatedAt:         violation.UpdatedAt,
	}
}
