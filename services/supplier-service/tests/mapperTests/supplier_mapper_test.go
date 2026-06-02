package mappertests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
)

func TestToSupplierResponse(t *testing.T) {
	id := uuid.New()
	item := models.Supplier{ID: id, Login: "supplier", Email: "supplier@example.com"}

	result := mapper.ToSupplierResponse(item)

	require.Equal(t, id.String(), result.ID)
	require.Equal(t, item.Login, result.Login)
	require.Equal(t, item.Email, result.Email)
}
