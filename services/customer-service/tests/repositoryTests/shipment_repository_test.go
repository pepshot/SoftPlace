package repositorytests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestShipmentRepositoryGetList(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentRepository(t, mockPool)

	item := sampleShipment(uuid.New(), uuid.New())
	mockPool.ExpectQuery(sqlSelectShipment).
		WillReturnRows(shipmentRows(item))

	result, err := repo.GetList(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, item.ID, result[0].ID)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentRepositoryGetListQueryError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentRepository(t, mockPool)

	mockPool.ExpectQuery(sqlSelectShipment).
		WillReturnError(errors.New("query error"))

	_, err := repo.GetList(context.Background())
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentRepositoryGetListScanError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentRepository(t, mockPool)

	rows := shipmentRows(sampleShipment(uuid.New(), uuid.New()))
	rows.RowError(0, errors.New("scan error"))
	mockPool.ExpectQuery(sqlSelectShipment).
		WillReturnRows(rows)

	_, err := repo.GetList(context.Background())
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentRepositoryGetByID(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentRepository(t, mockPool)

	item := sampleShipment(uuid.New(), uuid.New())
	mockPool.ExpectQuery(sqlSelectShipmentFromID).
		WithArgs(item.ID).
		WillReturnRows(shipmentRows(item))

	result, err := repo.GetByID(context.Background(), item.ID)
	require.NoError(t, err)
	require.Equal(t, item.ID, result.ID)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentRepositoryGetByIDNotFound(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectQuery(sqlSelectShipmentFromID).
		WithArgs(itemID).
		WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetByID(context.Background(), itemID)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentRepositoryGetByIDQueryError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectQuery(sqlSelectShipmentFromID).
		WithArgs(itemID).
		WillReturnError(errors.New("query error"))

	_, err := repo.GetByID(context.Background(), itemID)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentRepositoryCreate(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentRepository(t, mockPool)

	item := sampleShipment(uuid.New(), uuid.New())
	mockPool.ExpectExec(sqlInsertShipment).
		WithArgs(item.ID, item.CustomerID, item.Code, item.Date, item.Price).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.Create(context.Background(), item)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentRepositoryCreateExecError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentRepository(t, mockPool)

	item := sampleShipment(uuid.New(), uuid.New())
	mockPool.ExpectExec(sqlInsertShipment).
		WithArgs(item.ID, item.CustomerID, item.Code, item.Date, item.Price).
		WillReturnError(errors.New("exec error"))

	err := repo.Create(context.Background(), item)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentRepositoryCreateBadDate(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentRepository(t, mockPool)

	item := sampleShipment(uuid.New(), uuid.New())
	item.Date = time.Time{}

	mockPool.ExpectExec(sqlInsertShipment).
		WithArgs(item.ID, item.CustomerID, item.Code, item.Date, item.Price).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.Create(context.Background(), item)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentRepositoryUpdate(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentRepository(t, mockPool)

	item := sampleShipment(uuid.New(), uuid.New())
	mockPool.ExpectExec(sqlUpdateShipment).
		WithArgs(item.ID, item.CustomerID, item.Code, item.Date, item.Price).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.Update(context.Background(), item)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentRepositoryUpdateNoRows(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentRepository(t, mockPool)

	item := sampleShipment(uuid.New(), uuid.New())
	mockPool.ExpectExec(sqlUpdateShipment).
		WithArgs(item.ID, item.CustomerID, item.Code, item.Date, item.Price).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := repo.Update(context.Background(), item)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentRepositoryDelete(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectExec(sqlDeleteShipment).
		WithArgs(itemID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.Delete(context.Background(), itemID)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestShipmentRepositoryDeleteNoRows(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newShipmentRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectExec(sqlDeleteShipment).
		WithArgs(itemID).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err := repo.Delete(context.Background(), itemID)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}
