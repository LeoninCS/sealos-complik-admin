package unban

import "time"

type UnbanNamespaceRequest struct {
	Namespace string `uri:"namespace" binding:"required,max=255"`
}

type CreateUnbanRequest struct {
	Namespace    string `json:"namespace" binding:"required,max=255"`
	OperatorName string `json:"operator_name" binding:"required,max=100"`
}

type UnbanResponse struct {
	Namespace    string    `json:"namespace"`
	OperatorName string    `json:"operator_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
