package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
)

func (h *Handler) ListSupplies(c *gin.Context) {
	items, err := h.supply.GetList(c.Request.Context())
	if err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) CreateSupply(c *gin.Context) {
	supplierID, ok := getSupplierID(c)
	if !ok {
		return
	}

	var req dto.SupplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	id, err := h.supply.Create(c.Request.Context(), supplierID, req)
	if err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id.String()})
}

func (h *Handler) GetSupply(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	item, err := h.supply.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) UpdateSupply(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.SupplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	if err := h.supply.Update(c.Request.Context(), id, req); err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteSupply(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.supply.Delete(c.Request.Context(), id); err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
