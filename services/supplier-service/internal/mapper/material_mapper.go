package mapper

import (
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
)

func ToMaterialResponse(item models.Material) dto.MaterialResponse {
	return dto.MaterialResponse{
		ID:         item.ID.String(),
		Name:       item.Name,
		Code:       item.Code,
		Price:      item.Price,
		StockCount: item.StockCount,
	}
}

func ToMaterialResponseList(items []models.Material) []dto.MaterialResponse {
	result := make([]dto.MaterialResponse, 0, len(items))
	for _, item := range items {
		result = append(result, ToMaterialResponse(item))
	}
	return result
}
