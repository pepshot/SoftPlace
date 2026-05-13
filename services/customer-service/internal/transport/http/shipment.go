package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
)

// ListShipments godoc
// @Summary List shipments
// @Tags Shipment
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.ShipmentResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /shipments [get]
func (h *Handler) ListShipments(c *gin.Context) {
	items, err := h.shipment.GetList(c.Request.Context())
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetShipment godoc
// @Summary Get shipment by ID
// @Tags Shipment
// @Security BearerAuth
// @Produce json
// @Param id path string true "Shipment ID"
// @Success 200 {object} dto.ShipmentResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /shipments/{id} [get]
func (h *Handler) GetShipment(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	item, err := h.shipment.GetByID(c.Request.Context(), id)
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.JSON(http.StatusOK, item)
}

// CreateShipment godoc
// @Summary Create shipment
// @Tags Shipment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.ShipmentRequest true "Shipment payload"
// @Success 201 {object} createResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /shipments [post]
func (h *Handler) CreateShipment(c *gin.Context) {
	customerID, ok := getUserID(c)
	if !ok {
		return
	}

	var req dto.ShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	id, err := h.shipment.Create(c.Request.Context(), customerID, req)
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.JSON(http.StatusCreated, createResponse{ID: id.String()})
}

// UpdateShipment godoc
// @Summary Update shipment
// @Tags Shipment
// @Security BearerAuth
// @Accept json
// @Param id path string true "Shipment ID"
// @Param request body dto.ShipmentRequest true "Shipment payload"
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /shipments/{id} [put]
func (h *Handler) UpdateShipment(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.ShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	if err := h.shipment.Update(c.Request.Context(), id, req); err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteShipment godoc
// @Summary Delete shipment
// @Tags Shipment
// @Security BearerAuth
// @Param id path string true "Shipment ID"
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /shipments/{id} [delete]
func (h *Handler) DeleteShipment(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.shipment.Delete(c.Request.Context(), id); err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.Status(http.StatusNoContent)
}
