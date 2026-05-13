package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
)

// ListGarnitures godoc
// @Summary List garnitures
// @Tags Garniture
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.GarnitureResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /garnitures [get]
func (h *Handler) ListGarnitures(c *gin.Context) {
	items, err := h.garniture.GetList(c.Request.Context())
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetGarniture godoc
// @Summary Get garniture by ID
// @Tags Garniture
// @Security BearerAuth
// @Produce json
// @Param id path string true "Garniture ID"
// @Success 200 {object} dto.GarnitureResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /garnitures/{id} [get]
func (h *Handler) GetGarniture(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	item, err := h.garniture.GetByID(c.Request.Context(), id)
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.JSON(http.StatusOK, item)
}

// CreateGarniture godoc
// @Summary Create garniture
// @Tags Garniture
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.GarnitureRequest true "Garniture payload"
// @Success 201 {object} createResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /garnitures [post]
func (h *Handler) CreateGarniture(c *gin.Context) {
	var req dto.GarnitureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	id, err := h.garniture.Create(c.Request.Context(), req)
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.JSON(http.StatusCreated, createResponse{ID: id.String()})
}

// UpdateGarniture godoc
// @Summary Update garniture
// @Tags Garniture
// @Security BearerAuth
// @Accept json
// @Param id path string true "Garniture ID"
// @Param request body dto.GarnitureRequest true "Garniture payload"
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /garnitures/{id} [put]
func (h *Handler) UpdateGarniture(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.GarnitureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	if err := h.garniture.Update(c.Request.Context(), id, req); err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteGarniture godoc
// @Summary Delete garniture
// @Tags Garniture
// @Security BearerAuth
// @Param id path string true "Garniture ID"
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /garnitures/{id} [delete]
func (h *Handler) DeleteGarniture(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.garniture.Delete(c.Request.Context(), id); err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.Status(http.StatusNoContent)
}
