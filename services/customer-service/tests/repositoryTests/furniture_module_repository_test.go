package repositorytests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/model/relations"
)

func TestFurnitureModuleRepositoryCreateMany(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureModuleRepository(t, mockPool)

	items := []relations.FurnitureModule{{
		FurnitureID: uuid.New(),
		ModuleID:    uuid.New(),
		Count:       2,
	}}

	mockPool.ExpectExec(sqlInsertFurnitureModule).
		WithArgs(items[0].FurnitureID, items[0].ModuleID, items[0].Count).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.CreateMany(context.Background(), items)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureModuleRepositoryCreateManyError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureModuleRepository(t, mockPool)

	items := []relations.FurnitureModule{{
		FurnitureID: uuid.New(),
		ModuleID:    uuid.New(),
		Count:       2,
	}}

	mockPool.ExpectExec(sqlInsertFurnitureModule).
		WithArgs(items[0].FurnitureID, items[0].ModuleID, items[0].Count).
		WillReturnError(errors.New("exec error"))

	err := repo.CreateMany(context.Background(), items)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureModuleRepositoryGetByFurnitureID(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureModuleRepository(t, mockPool)

	furnitureID := uuid.New()
	items := []relations.FurnitureModule{{
		FurnitureID: furnitureID,
		ModuleID:    uuid.New(),
		Count:       2,
	}}

	mockPool.ExpectQuery(sqlSelectFurnitureModule).
		WithArgs(furnitureID).
		WillReturnRows(furnitureModuleRows(items...))

	result, err := repo.GetByFurnitureID(context.Background(), furnitureID)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureModuleRepositoryGetByFurnitureIDQueryError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureModuleRepository(t, mockPool)

	furnitureID := uuid.New()
	mockPool.ExpectQuery(sqlSelectFurnitureModule).
		WithArgs(furnitureID).
		WillReturnError(errors.New("query error"))

	_, err := repo.GetByFurnitureID(context.Background(), furnitureID)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureModuleRepositoryDeleteByFurnitureID(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureModuleRepository(t, mockPool)

	furnitureID := uuid.New()
	mockPool.ExpectExec(sqlDeleteFurnitureModule).
		WithArgs(furnitureID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.DeleteByFurnitureID(context.Background(), furnitureID)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureModuleRepositoryDeleteByFurnitureIDError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureModuleRepository(t, mockPool)

	furnitureID := uuid.New()
	mockPool.ExpectExec(sqlDeleteFurnitureModule).
		WithArgs(furnitureID).
		WillReturnError(errors.New("exec error"))

	err := repo.DeleteByFurnitureID(context.Background(), furnitureID)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}
