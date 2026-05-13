package mappertests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
)

func TestToGarnitureResponse(t *testing.T) {
	item := model.Garniture{ID: uuid.New(), Name: "Set", Code: "GS-1", Price: 100, StockCount: 1}

	result := mapper.ToGarnitureResponse(item)

	require.Equal(t, item.ID.String(), result.ID)
	require.Equal(t, item.Name, result.Name)
	require.Equal(t, item.Code, result.Code)
	require.Equal(t, item.Price, result.Price)
	require.Equal(t, item.StockCount, result.StockCount)
}

func TestToGarnitureResponseList(t *testing.T) {
	items := []model.Garniture{
		{ID: uuid.New(), Name: "Set", Code: "GS-1", Price: 100, StockCount: 1},
		{ID: uuid.New(), Name: "Set2", Code: "GS-2", Price: 200, StockCount: 5},
	}

	result := mapper.ToGarnitureResponseList(items)

	require.Len(t, result, 2)
	require.Equal(t, items[0].ID.String(), result[0].ID)
	require.Equal(t, items[1].ID.String(), result[1].ID)
	require.Equal(t, items[0].Code, result[0].Code)
	require.Equal(t, items[1].Code, result[1].Code)
}

func TestToGarnitureResponseListEmpty(t *testing.T) {
	items := []model.Garniture{}

	result := mapper.ToGarnitureResponseList(items)

	require.Len(t, result, 0)
}
