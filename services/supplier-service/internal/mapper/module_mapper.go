package mapper

import (
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
)

func ToModuleResponse(item models.Module) dto.ModuleResponse {
	return dto.ModuleResponse{
		ID:         item.ID.String(),
		Name:       item.Name,
		Code:       item.Code,
		Price:      item.Price,
		StockCount: item.StockCount,
	}
}

func ToModuleResponseList(items []models.Module) []dto.ModuleResponse {
	result := make([]dto.ModuleResponse, 0, len(items))

	for _, item := range items {
		result = append(result, ToModuleResponse(item))
	}

	return result
}
