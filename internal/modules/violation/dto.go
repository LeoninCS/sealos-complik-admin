package violation

import "time"

type ViolationUserIDRequest struct {
	UserID uint64 `uri:"user_id" binding:"required,min=1"`
}

type CreateViolationRequest struct {
	UserID        uint64    `json:"user_id" binding:"required,min=1"`
	ViolationType string    `json:"violation_type" binding:"required,max=50"`
	ViolationDesc string    `json:"violation_desc" binding:"omitempty"`
	ViolationTime time.Time `json:"violation_time" binding:"required"`
}

type ViolationResponse struct {
	UserID        uint64    `json:"user_id"`
	ViolationType string    `json:"violation_type"`
	ViolationDesc string    `json:"violation_desc,omitempty"`
	ViolationTime time.Time `json:"violation_time"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ViolationStatusResponse struct {
	Violated bool `json:"violated"`
}
