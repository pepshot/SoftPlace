package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
)

// Register godoc
// @Summary Register customer
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterCustomerRequest true "Register payload"
// @Success 201 {object} dto.AuthCustomerResponse
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	response, err := h.customer.Register(c.Request.Context(), req)
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary Login customer
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginCustomerRequest true "Login payload"
// @Success 200 {object} dto.AuthCustomerResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	response, err := h.customer.Login(c.Request.Context(), req)
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Profile godoc
// @Summary Get customer profile
// @Tags Customer
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.CustomerResponse
// @Failure 401 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /profile [get]
func (h *Handler) Profile(c *gin.Context) {
	customerID, ok := getUserID(c)
	if !ok {
		return
	}

	response, err := h.customer.GetProfile(c.Request.Context(), customerID)
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, errorResponse{Error: message})
		return
	}

	c.JSON(http.StatusOK, response)
}
