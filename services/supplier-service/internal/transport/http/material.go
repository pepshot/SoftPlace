package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
)

func (h *Handler) ListMaterials(c *gin.Context) {
	items, err := h.material.GetList(c.Request.Context())
	if err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) CreateMaterial(c *gin.Context) {
	var req dto.MaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	id, err := h.material.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id.String()})
}

func (h *Handler) GetMaterial(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	item, err := h.material.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) UpdateMaterial(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.MaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	if err := h.material.Update(c.Request.Context(), id, req); err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteMaterial(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.material.Delete(c.Request.Context(), id); err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
