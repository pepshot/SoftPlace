package mappertests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
)

func TestToModuleResponse(t *testing.T) {
	id := uuid.New()
	item := models.Module{ID: id, Name: "Module", Code: "MOD-1", Price: 10, StockCount: 2}

	result := mapper.ToModuleResponse(item)
	require.Equal(t, id.String(), result.ID)
	require.Equal(t, item.Name, result.Name)
}

func TestToModuleResponseList(t *testing.T) {
	items := []models.Module{
		{ID: uuid.New(), Name: "Module", Code: "MOD-1", Price: 10, StockCount: 2},
		{ID: uuid.New(), Name: "Module 2", Code: "MOD-2", Price: 20, StockCount: 5},
	}
	result := mapper.ToModuleResponseList(items)
	require.Len(t, result, 2)
	require.Equal(t, items[0].ID.String(), result[0].ID)
	require.Equal(t, items[1].ID.String(), result[1].ID)
	require.Equal(t, items[0].Code, result[0].Code)
	require.Equal(t, items[1].Code, result[1].Code)
}

func TestToModuleResponseListEmpty(t *testing.T) {
	result := mapper.ToModuleResponseList([]models.Module{})
	require.Len(t, result, 0)
}
