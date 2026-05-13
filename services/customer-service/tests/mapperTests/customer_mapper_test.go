package mappertests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
)

func TestToCustomerResponse(t *testing.T) {
	id := uuid.New()
	customer := model.Customer{
		ID:    id,
		Login: "login",
		Email: "user@example.com",
	}

	result := mapper.ToCustomerResponse(customer)

	require.Equal(t, id.String(), result.ID)
	require.Equal(t, customer.Login, result.Login)
	require.Equal(t, customer.Email, result.Email)
}
