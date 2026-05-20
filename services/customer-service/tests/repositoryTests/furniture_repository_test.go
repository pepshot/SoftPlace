package repositorytests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestFurnitureRepositoryGetList(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	item := sampleFurniture(uuid.New())
	rows := furnitureRows(item)
	mockPool.ExpectQuery(sqlSelectFurniture).
		WillReturnRows(rows)

	result, err := repo.GetList(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, item.ID, result[0].ID)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureRepositoryGetListQueryError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	mockPool.ExpectQuery(sqlSelectFurniture).
		WillReturnError(errors.New("query error"))

	_, err := repo.GetList(context.Background())
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureRepositoryGetListScanError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	rows := furnitureRows(sampleFurniture(uuid.New()))
	rows.RowError(0, errors.New("scan error"))
	mockPool.ExpectQuery(sqlSelectFurniture).
		WillReturnRows(rows)

	_, err := repo.GetList(context.Background())
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureRepositoryGetByID(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	item := sampleFurniture(uuid.New())
	mockPool.ExpectQuery(sqlSelectFurnitureByID).
		WithArgs(item.ID).
		WillReturnRows(furnitureRows(item))

	result, err := repo.GetByID(context.Background(), item.ID)
	require.NoError(t, err)
	require.Equal(t, item.ID, result.ID)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureRepositoryGetByIDNotFound(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectQuery(sqlSelectFurnitureByID).
		WithArgs(itemID).
		WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetByID(context.Background(), itemID)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureRepositoryGetByIDQueryError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectQuery(sqlSelectFurnitureByID).
		WithArgs(itemID).
		WillReturnError(errors.New("query error"))

	_, err := repo.GetByID(context.Background(), itemID)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureRepositoryCreate(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	item := sampleFurniture(uuid.New())
	mockPool.ExpectExec("INSERT INTO furniture").
		WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.Create(context.Background(), item)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureRepositoryCreateExecError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	item := sampleFurniture(uuid.New())
	mockPool.ExpectExec("INSERT INTO furniture").
		WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).
		WillReturnError(errors.New("exec error"))

	err := repo.Create(context.Background(), item)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureRepositoryCreateNoRows(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	item := sampleFurniture(uuid.New())
	mockPool.ExpectExec("INSERT INTO furniture").
		WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).
		WillReturnResult(pgxmock.NewResult("INSERT", 0))

	err := repo.Create(context.Background(), item)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureRepositoryUpdate(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	item := sampleFurniture(uuid.New())
	mockPool.ExpectExec("UPDATE furniture").
		WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.Update(context.Background(), item)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureRepositoryUpdateNoRows(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	item := sampleFurniture(uuid.New())
	mockPool.ExpectExec("UPDATE furniture").
		WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := repo.Update(context.Background(), item)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureRepositoryDelete(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectExec(sqlDeleteFurniture).
		WithArgs(itemID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.Delete(context.Background(), itemID)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureRepositoryDeleteNoRows(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectExec(sqlDeleteFurniture).
		WithArgs(itemID).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err := repo.Delete(context.Background(), itemID)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureRepositoryIncreaseStock(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectExec(sqlIncreaseFurniture).
		WithArgs(itemID, 2).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.IncreaseStock(context.Background(), itemID, 2)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureRepositoryIncreaseStockNoRows(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectExec(sqlIncreaseFurniture).
		WithArgs(itemID, 2).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := repo.IncreaseStock(context.Background(), itemID, 2)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureRepositoryDecreaseStock(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectExec(sqlDecreaseFurniture).
		WithArgs(itemID, 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.DecreaseStock(context.Background(), itemID, 1)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestFurnitureRepositoryDecreaseStockNoRows(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newFurnitureRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectExec(sqlDecreaseFurniture).
		WithArgs(itemID, 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := repo.DecreaseStock(context.Background(), itemID, 1)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}
