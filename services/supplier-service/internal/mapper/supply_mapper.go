package mapper

import (
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
)

func ToSupplyResponse(item models.Supply) dto.SupplyResponse {
	return dto.SupplyResponse{
		ID:         item.ID.String(),
		SupplierID: item.SupplierID.String(),
		Code:       item.Code,
		Date:       item.Date.Format("2006-01-02"),
		Price:      item.Price,
	}
}

func ToSupplyResponseList(items []models.Supply) []dto.SupplyResponse {
	result := make([]dto.SupplyResponse, 0, len(items))
	for _, item := range items {
		result = append(result, ToSupplyResponse(item))
	}
	return result
}
