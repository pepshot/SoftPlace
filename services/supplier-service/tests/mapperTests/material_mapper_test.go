package mappertests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
)

func TestToMaterialResponse(t *testing.T) {
	id := uuid.New()
	item := models.Material{ID: id, Name: "Material", Code: "MAT-1", Price: 10, StockCount: 5}

	result := mapper.ToMaterialResponse(item)

	require.Equal(t, id.String(), result.ID)
	require.Equal(t, item.Name, result.Name)
	require.Equal(t, item.Code, result.Code)
	require.Equal(t, item.Price, result.Price)
	require.Equal(t, item.StockCount, result.StockCount)
}

func TestToMaterialResponseList(t *testing.T) {
	items := []models.Material{
		{ID: uuid.New(), Name: "Material", Code: "MAT-1", Price: 10, StockCount: 5},
		{ID: uuid.New(), Name: "Material 2", Code: "MAT-2", Price: 20, StockCount: 7},
	}
	result := mapper.ToMaterialResponseList(items)
	require.Len(t, result, 2)
	require.Equal(t, items[0].ID.String(), result[0].ID)
	require.Equal(t, items[1].ID.String(), result[1].ID)
	require.Equal(t, items[0].Code, result[0].Code)
	require.Equal(t, items[1].Code, result[1].Code)
}

func TestToMaterialResponseListEmpty(t *testing.T) {
	result := mapper.ToMaterialResponseList([]models.Material{})
	require.Len(t, result, 0)
}
