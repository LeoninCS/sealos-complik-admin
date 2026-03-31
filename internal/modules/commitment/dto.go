package commitment

import "time"

type CommitmentUserIDRequest struct {
	UserID uint64 `uri:"user_id" binding:"required,min=1"`
}

type CreateCommitmentRequest struct {
	UserID   uint64 `json:"user_id" binding:"required,min=1"`
	FileName string `json:"file_name" binding:"required,max=255"`
	FileURL  string `json:"file_url" binding:"required,max=512"`
}

type UpdateCommitmentRequest struct {
	FileName string `json:"file_name" binding:"required,max=255"`
	FileURL  string `json:"file_url" binding:"required,max=512"`
}

type CommitmentResponse struct {
	UserID    uint64    `json:"user_id"`
	FileName  string    `json:"file_name"`
	FileURL   string    `json:"file_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
