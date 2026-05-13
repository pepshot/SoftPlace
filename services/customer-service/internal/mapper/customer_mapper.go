package mapper

import (
	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
)

func ToCustomerResponse(customer model.Customer) dto.CustomerResponse {
	return dto.CustomerResponse{
		ID:    customer.ID.String(),
		Login: customer.Login,
		Email: customer.Email,
	}
}
