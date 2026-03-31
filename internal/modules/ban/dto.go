package ban

import "time"

type BanUserIDRequest struct {
	UserID uint64 `uri:"user_id" binding:"required,min=1"`
}

type CreateBanRequest struct {
	UserID       uint64     `json:"user_id" binding:"required,min=1"`
	Reason       string     `json:"reason" binding:"omitempty,max=500"`
	BanStartTime time.Time  `json:"ban_start_time" binding:"required"`
	BanEndTime   *time.Time `json:"ban_end_time"`
	OperatorName string     `json:"operator_name" binding:"required,max=100"`
}

type BanResponse struct {
	UserID       uint64     `json:"user_id"`
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
