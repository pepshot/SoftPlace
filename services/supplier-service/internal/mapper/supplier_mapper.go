package mapper

import (
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
)

func ToSupplierResponse(item models.Supplier) dto.SupplierResponse {
	return dto.SupplierResponse{
		ID:    item.ID.String(),
		Login: item.Login,
		Email: item.Email,
	}
}
