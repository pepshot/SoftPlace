package mappertests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
)

func TestToFurnitureResponse(t *testing.T) {
	id := uuid.New()
	item := model.Furniture{
		ID:         id,
		Name:       "Chair",
		Code:       "CH-1",
		Price:      10,
		StockCount: 3,
	}

	result := mapper.ToFurnitureResponse(item)

	require.Equal(t, id.String(), result.ID)
	require.Equal(t, item.Name, result.Name)
	require.Equal(t, item.Code, result.Code)
	require.Equal(t, item.Price, result.Price)
	require.Equal(t, item.StockCount, result.StockCount)
}

func TestToFurnitureResponseList(t *testing.T) {
	items := []model.Furniture{
		{ID: uuid.New(), Name: "Chair", Code: "CH-1", Price: 10, StockCount: 3},
		{ID: uuid.New(), Name: "Table", Code: "TB-1", Price: 20, StockCount: 5},
	}

	result := mapper.ToFurnitureResponseList(items)

	require.Len(t, result, 2)
	require.Equal(t, items[0].ID.String(), result[0].ID)
	require.Equal(t, items[1].ID.String(), result[1].ID)
	require.Equal(t, items[0].Name, result[0].Name)
	require.Equal(t, items[1].Name, result[1].Name)
}

func TestToFurnitureResponseListEmpty(t *testing.T) {
	items := []model.Furniture{}

	result := mapper.ToFurnitureResponseList(items)

	require.Len(t, result, 0)
}
