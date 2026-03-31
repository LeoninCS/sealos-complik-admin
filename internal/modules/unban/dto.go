package unban

import "time"

type UnbanUserIDRequest struct {
	UserID uint64 `uri:"user_id" binding:"required,min=1"`
}

type CreateUnbanRequest struct {
	UserID       uint64 `json:"user_id" binding:"required,min=1"`
	OperatorName string `json:"operator_name" binding:"required,max=100"`
}

type UnbanResponse struct {
	UserID       uint64    `json:"user_id"`
	OperatorName string    `json:"operator_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
