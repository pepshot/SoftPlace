package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
)

type createResponse struct {
	ID string `json:"id"`
}

// ListFurniture godoc
// @Summary List furniture
// @Tags Furniture
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.FurnitureResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /furniture [get]
func (h *Handler) ListFurniture(c *gin.Context) {
	items, err := h.furniture.GetList(c.Request.Context())
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetFurniture godoc
// @Summary Get furniture by ID
// @Tags Furniture
// @Security BearerAuth
// @Produce json
// @Param id path string true "Furniture ID"
// @Success 200 {object} dto.FurnitureResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /furniture/{id} [get]
func (h *Handler) GetFurniture(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	item, err := h.furniture.GetByID(c.Request.Context(), id)
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.JSON(http.StatusOK, item)
}

// CreateFurniture godoc
// @Summary Create furniture
// @Tags Furniture
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.FurnitureRequest true "Furniture payload"
// @Success 201 {object} createResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /furniture [post]
func (h *Handler) CreateFurniture(c *gin.Context) {
	var req dto.FurnitureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	id, err := h.furniture.Create(c.Request.Context(), req)
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.JSON(http.StatusCreated, createResponse{ID: id.String()})
}

// UpdateFurniture godoc
// @Summary Update furniture
// @Tags Furniture
// @Security BearerAuth
// @Accept json
// @Param id path string true "Furniture ID"
// @Param request body dto.FurnitureRequest true "Furniture payload"
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /furniture/{id} [put]
func (h *Handler) UpdateFurniture(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.FurnitureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	if err := h.furniture.Update(c.Request.Context(), id, req); err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteFurniture godoc
// @Summary Delete furniture
// @Tags Furniture
// @Security BearerAuth
// @Param id path string true "Furniture ID"
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /furniture/{id} [delete]
func (h *Handler) DeleteFurniture(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.furniture.Delete(c.Request.Context(), id); err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.Status(http.StatusNoContent)
}
