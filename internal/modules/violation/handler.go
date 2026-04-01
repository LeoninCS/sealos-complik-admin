package violation

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

// CreateViolation handles the creation of a new violation record.
func (h *Handler) CreateViolation(c *gin.Context) {
	var req CreateViolationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.CreateViolation(c.Request.Context(), req); err != nil {
		h.respondWithServiceError(c, err, "failed to create violation")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "violation created successfully",
	})
}

// DeleteViolations handles deleting all violation records for a namespace.
func (h *Handler) DeleteViolations(c *gin.Context) {
	namespace, ok := bindViolationNamespace(c)
	if !ok {
		return
	}

	if err := h.service.DeleteViolations(c.Request.Context(), namespace); err != nil {
		h.respondWithServiceError(c, err, "failed to delete violations")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "violations deleted successfully",
	})
}

// GetViolations handles retrieving all violation records for a namespace.
func (h *Handler) GetViolations(c *gin.Context) {
	namespace, ok := bindViolationNamespace(c)
	if !ok {
		return
	}

	resp, err := h.service.GetViolations(c.Request.Context(), namespace)
	if err != nil {
		h.respondWithServiceError(c, err, "failed to get violations")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListViolations handles listing all violation records.
func (h *Handler) ListViolations(c *gin.Context) {
	resp, err := h.service.ListViolations(c.Request.Context())
	if err != nil {
		h.respondWithServiceError(c, err, "failed to list violations")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetViolationStatus handles checking whether a namespace has any violation records.
func (h *Handler) GetViolationStatus(c *gin.Context) {
	namespace, ok := bindViolationNamespace(c)
	if !ok {
		return
	}

	resp, err := h.service.GetViolationStatus(c.Request.Context(), namespace)
	if err != nil {
		h.respondWithServiceError(c, err, "failed to get violation status")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// bindViolationNamespace extracts the namespace from the URI and validates it.
func bindViolationNamespace(c *gin.Context) (string, bool) {
	var req ViolationNamespaceRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request path",
			"error":   err.Error(),
		})
		return "", false
	}

	return req.Namespace, true
}

// respondWithServiceError handles responding with appropriate error messages based on the service error.
func (h *Handler) respondWithServiceError(c *gin.Context, err error, fallbackMessage string) {
	switch {
	case errors.Is(err, ErrViolationInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	case errors.Is(err, ErrViolationNotFound):
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
