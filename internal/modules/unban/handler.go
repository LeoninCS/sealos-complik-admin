package unban

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

// CreateUnban handles the creation of a new unban record.
func (h *Handler) CreateUnban(c *gin.Context) {
	var req CreateUnbanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.CreateUnban(c.Request.Context(), req); err != nil {
		h.respondWithServiceError(c, err, "failed to create unban")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "unban created successfully",
	})
}

// DeleteUnbans handles deleting all unban records for a user.
func (h *Handler) DeleteUnbans(c *gin.Context) {
	userID, ok := bindUnbanUserID(c)
	if !ok {
		return
	}

	if err := h.service.DeleteUnbans(c.Request.Context(), userID); err != nil {
		h.respondWithServiceError(c, err, "failed to delete unbans")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "unbans deleted successfully",
	})
}

// GetUnbans handles retrieving all unban records for a user.
func (h *Handler) GetUnbans(c *gin.Context) {
	userID, ok := bindUnbanUserID(c)
	if !ok {
		return
	}

	resp, err := h.service.GetUnbans(c.Request.Context(), userID)
	if err != nil {
		h.respondWithServiceError(c, err, "failed to get unbans")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListUnbans handles listing all unban records.
func (h *Handler) ListUnbans(c *gin.Context) {
	resp, err := h.service.ListUnbans(c.Request.Context())
	if err != nil {
		h.respondWithServiceError(c, err, "failed to list unbans")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// bindUnbanUserID extracts the user ID from the URI and validates it.
func bindUnbanUserID(c *gin.Context) (uint64, bool) {
	var req UnbanUserIDRequest
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
	case errors.Is(err, ErrUnbanInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	case errors.Is(err, ErrUnbanNotFound):
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
