package ban

import "time"

type BanNamespaceRequest struct {
	Namespace string `uri:"namespace" binding:"required,max=255"`
}

type CreateBanRequest struct {
	Namespace    string     `json:"namespace" binding:"required,max=255"`
	Reason       string     `json:"reason" binding:"omitempty,max=500"`
	BanStartTime time.Time  `json:"ban_start_time" binding:"required"`
	BanEndTime   *time.Time `json:"ban_end_time"`
	OperatorName string     `json:"operator_name" binding:"required,max=100"`
}

type BanResponse struct {
	Namespace    string     `json:"namespace"`
	Reason       string     `json:"reason,omitempty"`
	BanStartTime time.Time  `json:"ban_start_time"`
	BanEndTime   *time.Time `json:"ban_end_time,omitempty"`
	OperatorName string     `json:"operator_name"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type BanStatusResponse struct {
	Banned bool `json:"banned"`
}
