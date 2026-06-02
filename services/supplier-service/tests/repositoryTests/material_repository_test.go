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

func TestMaterialRepositoryCreate(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	item := sampleMaterial(uuid.New())

	pool.ExpectExec(sqlInsertMaterial).WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.Create(context.Background(), item)
	require.NoError(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryGetList(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	item := sampleMaterial(uuid.New())

	pool.ExpectQuery(sqlSelectMaterials).WillReturnRows(pgxmock.NewRows([]string{"id", "name", "code", "price", "stock_count"}).AddRow(item.ID, item.Name, item.Code, item.Price, item.StockCount))

	result, err := repo.GetList(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, item.ID, result[0].ID)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryGetListQueryError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)

	pool.ExpectQuery(sqlSelectMaterials).WillReturnError(errors.New("query error"))

	_, err := repo.GetList(context.Background())
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryGetListScanError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)

	rows := pgxmock.NewRows([]string{"id", "name", "code", "price", "stock_count"}).AddRow(uuid.New(), "Material", "MAT-1", 10, 5)
	rows.RowError(0, errors.New("scan error"))
	pool.ExpectQuery(sqlSelectMaterials).WillReturnRows(rows)

	_, err := repo.GetList(context.Background())
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryGetByID(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	item := sampleMaterial(uuid.New())

	pool.ExpectQuery(sqlSelectMaterialByID).WithArgs(item.ID).WillReturnRows(materialRows(item))

	result, err := repo.GetByID(context.Background(), item.ID)
	require.NoError(t, err)
	require.Equal(t, item.ID, result.ID)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryGetByIDNotFound(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	id := uuid.New()

	pool.ExpectQuery(sqlSelectMaterialByID).WithArgs(id).WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetByID(context.Background(), id)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryGetByIDQueryError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	id := uuid.New()

	pool.ExpectQuery(sqlSelectMaterialByID).WithArgs(id).WillReturnError(errors.New("query error"))

	_, err := repo.GetByID(context.Background(), id)
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryExistsByCodeTrue(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(true)
	pool.ExpectQuery(sqlExistsMaterialCode).WithArgs("MAT-1").WillReturnRows(rows)

	exists, err := repo.ExistsByCode(context.Background(), "MAT-1")
	require.NoError(t, err)
	require.True(t, exists)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryExistsByCodeFalse(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(false)
	pool.ExpectQuery(sqlExistsMaterialCode).WithArgs("MAT-1").WillReturnRows(rows)

	exists, err := repo.ExistsByCode(context.Background(), "MAT-1")
	require.NoError(t, err)
	require.False(t, exists)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryExistsByCodeError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)

	pool.ExpectQuery(sqlExistsMaterialCode).WithArgs("MAT-1").WillReturnError(errors.New("query error"))

	exists, err := repo.ExistsByCode(context.Background(), "MAT-1")
	require.Error(t, err)
	require.False(t, exists)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryCreateNoRows(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	item := sampleMaterial(uuid.New())

	pool.ExpectExec(sqlInsertMaterial).WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).WillReturnResult(pgxmock.NewResult("INSERT", 0))

	err := repo.Create(context.Background(), item)
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestMaterialRepositoryCreateError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	item := sampleMaterial(uuid.New())

	pool.ExpectExec(sqlInsertMaterial).WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).WillReturnError(errors.New("exec error"))

	err := repo.Create(context.Background(), item)
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryUpdateNoRows(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	item := sampleMaterial(uuid.New())

	pool.ExpectExec(sqlUpdateMaterial).WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).WillReturnResult(pgxmock.NewResult("UPDATE", 0))
	err := repo.Update(context.Background(), item)
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestMaterialRepositoryUpdate(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	item := sampleMaterial(uuid.New())

	pool.ExpectExec(sqlUpdateMaterial).WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.Update(context.Background(), item)
	require.NoError(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryUpdateError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	item := sampleMaterial(uuid.New())

	pool.ExpectExec(sqlUpdateMaterial).WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).WillReturnError(errors.New("update error"))

	err := repo.Update(context.Background(), item)
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryDeleteError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlDeleteMaterial).WithArgs(id).WillReturnError(errors.New("delete error"))
	err := repo.Delete(context.Background(), id)
	require.Error(t, err)
}

func TestMaterialRepositoryDelete(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlDeleteMaterial).WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryDeleteNoRows(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlDeleteMaterial).WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err := repo.Delete(context.Background(), id)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryIncreaseStock(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlIncreaseMaterial).WithArgs(id, 2).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.IncreaseStock(context.Background(), id, 2)
	require.NoError(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryIncreaseStockNoRows(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlIncreaseMaterial).WithArgs(id, 2).WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := repo.IncreaseStock(context.Background(), id, 2)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryIncreaseStockError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlIncreaseMaterial).WithArgs(id, 2).WillReturnError(errors.New("exec error"))

	err := repo.IncreaseStock(context.Background(), id, 2)
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryDecreaseStock(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlDecreaseMaterial).WithArgs(id, 1).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.DecreaseStock(context.Background(), id, 1)
	require.NoError(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryDecreaseStockNoRows(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlDecreaseMaterial).WithArgs(id, 1).WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := repo.DecreaseStock(context.Background(), id, 1)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestMaterialRepositoryDecreaseStockError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newMaterialRepo(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlDecreaseMaterial).WithArgs(id, 1).WillReturnError(errors.New("exec error"))

	err := repo.DecreaseStock(context.Background(), id, 1)
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}
