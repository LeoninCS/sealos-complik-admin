package complikviolation

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
		h.respondWithServiceError(c, err, "failed to create complik violation")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "complik violation created successfully",
	})
}

func (h *Handler) DeleteViolations(c *gin.Context) {
	namespace, ok := bindNamespace(c)
	if !ok {
		return
	}

	if err := h.service.DeleteViolations(c.Request.Context(), namespace); err != nil {
		h.respondWithServiceError(c, err, "failed to delete complik violations")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "complik violations deleted successfully",
	})
}

func (h *Handler) GetViolations(c *gin.Context) {
	namespace, ok := bindNamespace(c)
	if !ok {
		return
	}

	resp, err := h.service.GetViolations(c.Request.Context(), namespace)
	if err != nil {
		h.respondWithServiceError(c, err, "failed to get complik violations")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListViolations(c *gin.Context) {
	resp, err := h.service.ListViolations(c.Request.Context())
	if err != nil {
		h.respondWithServiceError(c, err, "failed to list complik violations")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateViolationStatus(c *gin.Context) {
	var uriReq ViolationIDRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request path",
			"error":   err.Error(),
		})
		return
	}

	var req UpdateViolationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.UpdateViolationStatus(c.Request.Context(), uriReq.ID, req.Status); err != nil {
		h.respondWithServiceError(c, err, "failed to update complik violation status")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "complik violation status updated successfully",
	})
}

func (h *Handler) GetViolationStatus(c *gin.Context) {
	namespace, ok := bindNamespace(c)
	if !ok {
		return
	}

	resp, err := h.service.GetViolationStatus(c.Request.Context(), namespace)
	if err != nil {
		h.respondWithServiceError(c, err, "failed to get complik violation status")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func bindNamespace(c *gin.Context) (string, bool) {
	var req NamespaceRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request path",
			"error":   err.Error(),
		})
		return "", false
	}

	return req.Namespace, true
}

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
