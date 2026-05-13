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

func TestGarnitureRepositoryGetList(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureRepository(t, mockPool)

	item := sampleGarniture(uuid.New())
	mockPool.ExpectQuery(sqlSelectGarniture).
		WillReturnRows(garnitureRows(item))

	result, err := repo.GetList(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, item.ID, result[0].ID)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureRepositoryGetListQueryError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureRepository(t, mockPool)

	mockPool.ExpectQuery(sqlSelectGarniture).
		WillReturnError(errors.New("query error"))

	_, err := repo.GetList(context.Background())
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureRepositoryGetListScanError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureRepository(t, mockPool)

	rows := garnitureRows(sampleGarniture(uuid.New()))
	rows.RowError(0, errors.New("scan error"))
	mockPool.ExpectQuery(sqlSelectGarniture).
		WillReturnRows(rows)

	_, err := repo.GetList(context.Background())
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureRepositoryGetByID(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureRepository(t, mockPool)

	item := sampleGarniture(uuid.New())
	mockPool.ExpectQuery(sqlSelectGarnitureFromID).
		WithArgs(item.ID).
		WillReturnRows(garnitureRows(item))

	result, err := repo.GetByID(context.Background(), item.ID)
	require.NoError(t, err)
	require.Equal(t, item.ID, result.ID)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureRepositoryGetByIDNotFound(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectQuery(sqlSelectGarnitureFromID).
		WithArgs(itemID).
		WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetByID(context.Background(), itemID)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureRepositoryGetByIDQueryError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectQuery(sqlSelectGarnitureFromID).
		WithArgs(itemID).
		WillReturnError(errors.New("query error"))

	_, err := repo.GetByID(context.Background(), itemID)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureRepositoryCreate(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureRepository(t, mockPool)

	item := sampleGarniture(uuid.New())
	mockPool.ExpectExec(sqlInsertGarniture).
		WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.Create(context.Background(), item)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureRepositoryCreateExecError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureRepository(t, mockPool)

	item := sampleGarniture(uuid.New())
	mockPool.ExpectExec(sqlInsertGarniture).
		WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).
		WillReturnError(errors.New("exec error"))

	err := repo.Create(context.Background(), item)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureRepositoryUpdate(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureRepository(t, mockPool)

	item := sampleGarniture(uuid.New())
	mockPool.ExpectExec(sqlUpdateGarniture).
		WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.Update(context.Background(), item)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureRepositoryUpdateNoRows(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureRepository(t, mockPool)

	item := sampleGarniture(uuid.New())
	mockPool.ExpectExec(sqlUpdateGarniture).
		WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := repo.Update(context.Background(), item)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureRepositoryDelete(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectExec(sqlDeleteGarniture).
		WithArgs(itemID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.Delete(context.Background(), itemID)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureRepositoryDeleteNoRows(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectExec(sqlDeleteGarniture).
		WithArgs(itemID).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err := repo.Delete(context.Background(), itemID)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureRepositoryIncreaseStock(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectExec(sqlIncreaseGarniture).
		WithArgs(itemID, 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.IncreaseStock(context.Background(), itemID, 1)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureRepositoryIncreaseStockNoRows(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectExec(sqlIncreaseGarniture).
		WithArgs(itemID, 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := repo.IncreaseStock(context.Background(), itemID, 1)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureRepositoryDecreaseStock(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectExec(sqlDecreaseGarniture).
		WithArgs(itemID, 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.DecreaseStock(context.Background(), itemID, 1)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestGarnitureRepositoryDecreaseStockNoRows(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newGarnitureRepository(t, mockPool)

	itemID := uuid.New()
	mockPool.ExpectExec(sqlDecreaseGarniture).
		WithArgs(itemID, 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := repo.DecreaseStock(context.Background(), itemID, 1)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}
