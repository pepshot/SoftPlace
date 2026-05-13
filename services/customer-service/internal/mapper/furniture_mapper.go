package mapper

import (
	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
)

func ToFurnitureResponse(item model.Furniture) dto.FurnitureResponse {
	return dto.FurnitureResponse{
		ID:         item.ID.String(),
		Name:       item.Name,
		Code:       item.Code,
		Price:      item.Price,
		StockCount: item.StockCount,
	}
}

func ToFurnitureResponseList(items []model.Furniture) []dto.FurnitureResponse {
	result := make([]dto.FurnitureResponse, 0, len(items))

	for _, item := range items {
		result = append(result, ToFurnitureResponse(item))
	}

	return result
}
