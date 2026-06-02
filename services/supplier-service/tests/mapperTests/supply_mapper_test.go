package mappertests

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
)

func TestToSupplyResponse(t *testing.T) {
	id := uuid.New()
	supplierID := uuid.New()
	item := models.Supply{ID: id, SupplierID: supplierID, Code: "SUP-1", Date: time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC), Price: 20}

	result := mapper.ToSupplyResponse(item)
	require.Equal(t, id.String(), result.ID)
	require.Equal(t, supplierID.String(), result.SupplierID)
	require.Equal(t, "2026-05-20", result.Date)
}

func TestToSupplyResponseList(t *testing.T) {
	items := []models.Supply{
		{ID: uuid.New(), SupplierID: uuid.New(), Code: "SUP-1", Date: time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC), Price: 20},
		{ID: uuid.New(), SupplierID: uuid.New(), Code: "SUP-2", Date: time.Date(2026, 5, 21, 0, 0, 0, 0, time.UTC), Price: 30},
	}
	result := mapper.ToSupplyResponseList(items)
	require.Len(t, result, 2)
	require.Equal(t, items[0].ID.String(), result[0].ID)
	require.Equal(t, items[1].ID.String(), result[1].ID)
	require.Equal(t, "2026-05-20", result[0].Date)
	require.Equal(t, "2026-05-21", result[1].Date)
}

func TestToSupplyResponseListEmpty(t *testing.T) {
	result := mapper.ToSupplyResponseList([]models.Supply{})
	require.Len(t, result, 0)
}
