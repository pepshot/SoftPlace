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

func TestShipmentGarnitureRepositoryCreateMany(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentGarnitureRepository(t, mockPool)

	items := []relations.ShipmentGarniture{{
		ShipmentID:  uuid.New(),
		GarnitureID: uuid.New(),
		Count:       1,
	}}

	mockPool.ExpectExec(sqlInsertShipmentGarniture).
		WithArgs(items[0].ShipmentID, items[0].GarnitureID, items[0].Count).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.CreateMany(context.Background(), items)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentGarnitureRepositoryCreateManyError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentGarnitureRepository(t, mockPool)

	items := []relations.ShipmentGarniture{{
		ShipmentID:  uuid.New(),
		GarnitureID: uuid.New(),
		Count:       1,
	}}

	mockPool.ExpectExec(sqlInsertShipmentGarniture).
		WithArgs(items[0].ShipmentID, items[0].GarnitureID, items[0].Count).
		WillReturnError(errors.New("exec error"))

	err := repo.CreateMany(context.Background(), items)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentGarnitureRepositoryGetByShipmentID(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentGarnitureRepository(t, mockPool)

	shipmentID := uuid.New()
	items := []relations.ShipmentGarniture{{
		ShipmentID:  shipmentID,
		GarnitureID: uuid.New(),
		Count:       1,
	}}

	mockPool.ExpectQuery(sqlSelectShipmentGarniture).
		WithArgs(shipmentID).
		WillReturnRows(shipmentGarnitureRows(items...))

	result, err := repo.GetByShipmentID(context.Background(), shipmentID)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentGarnitureRepositoryGetByShipmentIDQueryError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentGarnitureRepository(t, mockPool)

	shipmentID := uuid.New()
	mockPool.ExpectQuery(sqlSelectShipmentGarniture).
		WithArgs(shipmentID).
		WillReturnError(errors.New("query error"))

	_, err := repo.GetByShipmentID(context.Background(), shipmentID)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentGarnitureRepositoryDeleteByShipmentID(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentGarnitureRepository(t, mockPool)

	shipmentID := uuid.New()
	mockPool.ExpectExec(sqlDeleteShipmentGarniture).
		WithArgs(shipmentID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.DeleteByShipmentID(context.Background(), shipmentID)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentGarnitureRepositoryDeleteByShipmentIDError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentGarnitureRepository(t, mockPool)

	shipmentID := uuid.New()
	mockPool.ExpectExec(sqlDeleteShipmentGarniture).
		WithArgs(shipmentID).
		WillReturnError(errors.New("exec error"))

	err := repo.DeleteByShipmentID(context.Background(), shipmentID)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}
