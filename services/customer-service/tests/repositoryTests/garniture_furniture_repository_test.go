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

func TestGarnitureFurnitureRepositoryCreateMany(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureFurnitureRepository(t, mockPool)

	items := []relations.GarnitureFurniture{{
		GarnitureID: uuid.New(),
		FurnitureID: uuid.New(),
		Count:       1,
	}}

	mockPool.ExpectExec(sqlInsertGarnitureFurniture).
		WithArgs(items[0].GarnitureID, items[0].FurnitureID, items[0].Count).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.CreateMany(context.Background(), items)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureFurnitureRepositoryCreateManyError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureFurnitureRepository(t, mockPool)

	items := []relations.GarnitureFurniture{{
		GarnitureID: uuid.New(),
		FurnitureID: uuid.New(),
		Count:       1,
	}}

	mockPool.ExpectExec(sqlInsertGarnitureFurniture).
		WithArgs(items[0].GarnitureID, items[0].FurnitureID, items[0].Count).
		WillReturnError(errors.New("exec error"))

	err := repo.CreateMany(context.Background(), items)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureFurnitureRepositoryGetByGarnitureID(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureFurnitureRepository(t, mockPool)

	garnitureID := uuid.New()
	items := []relations.GarnitureFurniture{{
		GarnitureID: garnitureID,
		FurnitureID: uuid.New(),
		Count:       1,
	}}

	mockPool.ExpectQuery(sqlSelectGarnitureFurniture).
		WithArgs(garnitureID).
		WillReturnRows(garnitureFurnitureRows(items...))

	result, err := repo.GetByGarnitureID(context.Background(), garnitureID)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureFurnitureRepositoryGetByGarnitureIDQueryError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureFurnitureRepository(t, mockPool)

	garnitureID := uuid.New()
	mockPool.ExpectQuery(sqlSelectGarnitureFurniture).
		WithArgs(garnitureID).
		WillReturnError(errors.New("query error"))

	_, err := repo.GetByGarnitureID(context.Background(), garnitureID)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureFurnitureRepositoryDeleteByGarnitureID(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureFurnitureRepository(t, mockPool)

	garnitureID := uuid.New()
	mockPool.ExpectExec(sqlDeleteGarnitureFurniture).
		WithArgs(garnitureID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.DeleteByGarnitureID(context.Background(), garnitureID)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureFurnitureRepositoryDeleteByGarnitureIDError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureFurnitureRepository(t, mockPool)

	garnitureID := uuid.New()
	mockPool.ExpectExec(sqlDeleteGarnitureFurniture).
		WithArgs(garnitureID).
		WillReturnError(errors.New("exec error"))

	err := repo.DeleteByGarnitureID(context.Background(), garnitureID)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}
