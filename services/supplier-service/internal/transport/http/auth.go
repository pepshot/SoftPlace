package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
)

func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	resp, err := h.supplier.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	resp, err := h.supplier.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Profile(c *gin.Context) {
	supplierID, ok := getSupplierID(c)
	if !ok {
		return
	}

	resp, err := h.supplier.GetProfile(c.Request.Context(), supplierID)
	if err != nil {
		c.JSON(mapServiceError(err), errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
