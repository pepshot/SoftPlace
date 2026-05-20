package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
)

func (h *Handler) ListModules(c *gin.Context) {
	items, err := h.module.GetList(c.Request.Context())
	if err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) CreateModule(c *gin.Context) {
	var req dto.ModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	id, err := h.module.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id.String()})
}

func (h *Handler) GetModule(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	item, err := h.module.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) UpdateModule(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.ModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	if err := h.module.Update(c.Request.Context(), id, req); err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteModule(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.module.Delete(c.Request.Context(), id); err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
