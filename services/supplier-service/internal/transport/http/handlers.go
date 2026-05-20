package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/service"
	"github.com/pepshot/SoftPlace/shared/logger"
)

const userIDKey = "userID"

type Handler struct {
	module   service.ModuleUseCase
	material service.MaterialUseCase
	supply   service.SupplyUseCase
	supplier service.SupplierUseCase
	logger   *logger.Logger
}

func NewHandler(
	module service.ModuleUseCase,
	material service.MaterialUseCase,
	supply service.SupplyUseCase,
	supplier service.SupplierUseCase,
	logger *logger.Logger,
) *Handler {
	return &Handler{
		module:   module,
		material: material,
		supply:   supply,
		supplier: supplier,
		logger:   logger,
	}
}

func parseUUIDParam(c *gin.Context, name string) (uuid.UUID, bool) {
	idParam := c.Param(name)
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return uuid.Nil, false
	}

	return id, true
}

func getSupplierID(c *gin.Context) (uuid.UUID, bool) {
	value, ok := c.Get(userIDKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return uuid.Nil, false
	}

	userID, ok := value.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return uuid.Nil, false
	}

	parsed, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return uuid.Nil, false
	}

	return parsed, true
}
