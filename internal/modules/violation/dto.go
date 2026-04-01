package violation

import "time"

type ViolationNamespaceRequest struct {
	Namespace string `uri:"namespace" binding:"required,max=255"`
}

type CreateViolationRequest struct {
	Namespace     string    `json:"namespace" binding:"required,max=255"`
	ViolationType string    `json:"violation_type" binding:"required,max=50"`
	ViolationDesc string    `json:"violation_desc" binding:"omitempty"`
	ViolationTime time.Time `json:"violation_time" binding:"required"`
}

type ViolationResponse struct {
	Namespace     string    `json:"namespace"`
	ViolationType string    `json:"violation_type"`
	ViolationDesc string    `json:"violation_desc,omitempty"`
	ViolationTime time.Time `json:"violation_time"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ViolationStatusResponse struct {
	Violated bool `json:"violated"`
}
