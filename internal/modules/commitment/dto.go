package commitment

import "time"

type CommitmentNamespaceRequest struct {
	Namespace string `uri:"namespace" binding:"required,max=255"`
}

type CreateCommitmentRequest struct {
	Namespace string `json:"namespace" binding:"required,max=255"`
	FileName  string `json:"file_name" binding:"required,max=255"`
	FileURL   string `json:"file_url" binding:"required,max=512"`
}

type UpdateCommitmentRequest struct {
	FileName string `json:"file_name" binding:"required,max=255"`
	FileURL  string `json:"file_url" binding:"required,max=512"`
}

type CommitmentResponse struct {
	Namespace string    `json:"namespace"`
	FileName  string    `json:"file_name"`
	FileURL   string    `json:"file_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
