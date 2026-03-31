package ban

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

// CreateBan handles the creation of a new ban record.
func (h *Handler) CreateBan(c *gin.Context) {
	var req CreateBanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.CreateBan(c.Request.Context(), req); err != nil {
		h.respondWithServiceError(c, err, "failed to create ban")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "ban created successfully",
	})
}

// DeleteBans handles deleting all ban records for a user.
func (h *Handler) DeleteBans(c *gin.Context) {
	userID, ok := bindBanUserID(c)
	if !ok {
		return
	}

	if err := h.service.DeleteBans(c.Request.Context(), userID); err != nil {
		h.respondWithServiceError(c, err, "failed to delete bans")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "bans deleted successfully",
	})
}

// GetBans handles retrieving all ban records for a user.
func (h *Handler) GetBans(c *gin.Context) {
	userID, ok := bindBanUserID(c)
	if !ok {
		return
	}

	resp, err := h.service.GetBans(c.Request.Context(), userID)
	if err != nil {
		h.respondWithServiceError(c, err, "failed to get bans")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListBans handles listing all ban records.
func (h *Handler) ListBans(c *gin.Context) {
	resp, err := h.service.ListBans(c.Request.Context())
	if err != nil {
		h.respondWithServiceError(c, err, "failed to list bans")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetBanStatus handles checking whether a user is currently banned.
func (h *Handler) GetBanStatus(c *gin.Context) {
	userID, ok := bindBanUserID(c)
	if !ok {
		return
	}

	resp, err := h.service.GetBanStatus(c.Request.Context(), userID)
	if err != nil {
		h.respondWithServiceError(c, err, "failed to get ban status")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// bindBanUserID extracts the user ID from the URI and validates it.
func bindBanUserID(c *gin.Context) (uint64, bool) {
	var req BanUserIDRequest
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
	case errors.Is(err, ErrBanInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	case errors.Is(err, ErrBanNotFound):
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
