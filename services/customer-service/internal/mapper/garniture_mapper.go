package mapper

import (
	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
)

func ToGarnitureResponse(item model.Garniture) dto.GarnitureResponse {
	return dto.GarnitureResponse{
		ID:         item.ID.String(),
		Name:       item.Name,
		Code:       item.Code,
		Price:      item.Price,
		StockCount: item.StockCount,
	}
}

func ToGarnitureResponseList(items []model.Garniture) []dto.GarnitureResponse {
	result := make([]dto.GarnitureResponse, 0, len(items))

	for _, item := range items {
		result = append(result, ToGarnitureResponse(item))
	}

	return result
}
