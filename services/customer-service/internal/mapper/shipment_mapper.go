package mapper

import (
	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
)

func ToShipmentResponse(item model.Shipment) dto.ShipmentResponse {
	return dto.ShipmentResponse{
		ID:    item.ID.String(),
		Code:  item.Code,
		Date:  item.Date.Format("2006-01-02"),
		Price: item.Price,
	}
}

func ToShipmentResponseList(items []model.Shipment) []dto.ShipmentResponse {
	result := make([]dto.ShipmentResponse, 0, len(items))

	for _, item := range items {
		result = append(result, ToShipmentResponse(item))
	}

	return result
}
