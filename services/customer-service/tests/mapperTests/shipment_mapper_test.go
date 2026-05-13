package mappertests

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
)

func TestToShipmentResponse(t *testing.T) {
	id := uuid.New()
	item := model.Shipment{
		ID:    id,
		Code:  "SH-1",
		Date:  time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC),
		Price: 200,
	}

	result := mapper.ToShipmentResponse(item)

	require.Equal(t, item.ID.String(), result.ID)
	require.Equal(t, "2025-05-10", result.Date)
	require.Equal(t, item.Price, result.Price)
}

func TestToShipmentResponseList(t *testing.T) {
	items := []model.Shipment{
		{ID: uuid.New(), Code: "SH-1", Date: time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC), Price: 200},
		{ID: uuid.New(), Code: "SH-2", Date: time.Date(2025, 5, 11, 0, 0, 0, 0, time.UTC), Price: 300},
	}

	result := mapper.ToShipmentResponseList(items)

	require.Len(t, result, 2)
	require.Equal(t, items[0].ID.String(), result[0].ID)
	require.Equal(t, items[1].ID.String(), result[1].ID)
	require.Equal(t, items[0].Code, result[0].Code)
	require.Equal(t, items[1].Code, result[1].Code)
	require.Equal(t, "2025-05-10", result[0].Date)
	require.Equal(t, "2025-05-11", result[1].Date)
}
