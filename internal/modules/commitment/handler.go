package commitment

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateCommitment handles the creation of a new commitment.
func (h *Handler) CreateCommitment(c *gin.Context) {
	var req CreateCommitmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.CreateCommitment(c.Request.Context(), req); err != nil {
		h.respondWithServiceError(c, err, "failed to create commitment")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "commitment created successfully",
	})
}

// UpdateCommitment handles updating a commitment.
func (h *Handler) UpdateCommitment(c *gin.Context) {
	userID, ok := bindCommitmentUserID(c)
	if !ok {
		return
	}

	var req UpdateCommitmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.UpdateCommitment(c.Request.Context(), userID, req); err != nil {
		h.respondWithServiceError(c, err, "failed to update commitment")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "commitment updated successfully",
	})
}

// DeleteCommitment handles deleting a commitment.
func (h *Handler) DeleteCommitment(c *gin.Context) {
	userID, ok := bindCommitmentUserID(c)
	if !ok {
		return
	}

	if err := h.service.DeleteCommitment(c.Request.Context(), userID); err != nil {
		h.respondWithServiceError(c, err, "failed to delete commitment")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "commitment deleted successfully",
	})
}

// GetCommitment handles retrieving a commitment by user ID.
func (h *Handler) GetCommitment(c *gin.Context) {
	userID, ok := bindCommitmentUserID(c)
	if !ok {
		return
	}

	resp, err := h.service.GetCommitment(c.Request.Context(), userID)
	if err != nil {
		h.respondWithServiceError(c, err, "failed to get commitment")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListCommitments handles listing all commitments.
func (h *Handler) ListCommitments(c *gin.Context) {
	resp, err := h.service.ListCommitments(c.Request.Context())
	if err != nil {
		h.respondWithServiceError(c, err, "failed to list commitments")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// bindCommitmentUserID extracts the user ID from the URI and validates it.
func bindCommitmentUserID(c *gin.Context) (uint64, bool) {
	var req CommitmentUserIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request path",
			"error":   err.Error(),
		})
		return 0, false
	}

	return req.UserID, true
}

// respondWithServiceError handles responding with appropriate error messages based on the service error.
func (h *Handler) respondWithServiceError(c *gin.Context, err error, fallbackMessage string) {
	switch {
	case errors.Is(err, ErrCommitmentInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	case errors.Is(err, ErrCommitmentAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{
			"message": err.Error(),
		})
	case errors.Is(err, ErrCommitmentNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fallbackMessage,
			"error":   err.Error(),
		})
	}
}
