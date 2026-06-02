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

func TestSupplyRepositoryGetList(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplyRepository(t, pool)
	item := sampleSupply(uuid.New(), uuid.New())

	pool.ExpectQuery(sqlSelectSupplies).WillReturnRows(supplyRows(item))

	result, err := repo.GetList(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, item.ID, result[0].ID)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplyRepositoryGetListQueryError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplyRepository(t, pool)

	pool.ExpectQuery(sqlSelectSupplies).WillReturnError(errors.New("query error"))

	_, err := repo.GetList(context.Background())
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplyRepositoryGetListScanError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplyRepository(t, pool)

	rows := pgxmock.NewRows([]string{"id", "supplier_id", "code", "date", "price"}).AddRow(uuid.New(), uuid.New(), "SUP-1", "bad", 20)
	rows.RowError(0, errors.New("scan error"))
	pool.ExpectQuery(sqlSelectSupplies).WillReturnRows(rows)

	_, err := repo.GetList(context.Background())
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplyRepositoryGetByID(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplyRepository(t, pool)
	item := sampleSupply(uuid.New(), uuid.New())

	pool.ExpectQuery(sqlSelectSupplyByID).WithArgs(item.ID).WillReturnRows(supplyRows(item))

	result, err := repo.GetByID(context.Background(), item.ID)
	require.NoError(t, err)
	require.Equal(t, item.ID, result.ID)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplyRepositoryGetByIDNotFound(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplyRepository(t, pool)
	id := uuid.New()

	pool.ExpectQuery(sqlSelectSupplyByID).WithArgs(id).WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetByID(context.Background(), id)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplyRepositoryGetByIDQueryError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplyRepository(t, pool)
	id := uuid.New()

	pool.ExpectQuery(sqlSelectSupplyByID).WithArgs(id).WillReturnError(errors.New("query error"))

	_, err := repo.GetByID(context.Background(), id)
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplyRepositoryExistsByCodeTrueFalseError(t *testing.T) {
	tests := []struct {
		name   string
		value  any
		err    error
		exists bool
	}{
		{name: "true", value: true, exists: true},
		{name: "false", value: false, exists: false},
		{name: "error", err: errors.New("query error"), exists: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := newMockPool(t)
			defer pool.Close()
			repo := newSupplyRepository(t, pool)

			if tt.err != nil {
				pool.ExpectQuery(sqlExistsSupplyCode).WithArgs("SUP-1").WillReturnError(tt.err)
			} else {
				rows := pgxmock.NewRows([]string{"exists"}).AddRow(tt.value)
				pool.ExpectQuery(sqlExistsSupplyCode).WithArgs("SUP-1").WillReturnRows(rows)
			}

			exists, err := repo.ExistsByCode(context.Background(), "SUP-1")
			if tt.err != nil {
				require.Error(t, err)
				require.False(t, exists)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.exists, exists)
			}
			require.NoError(t, pool.ExpectationsWereMet())
		})
	}
}

func TestSupplyRepositoryCreate(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplyRepository(t, pool)
	item := sampleSupply(uuid.New(), uuid.New())

	pool.ExpectExec(sqlInsertSupply).WithArgs(item.ID, item.SupplierID, item.Code, item.Date, item.Price).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.Create(context.Background(), item)
	require.NoError(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplyRepositoryCreateError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplyRepository(t, pool)
	item := sampleSupply(uuid.New(), uuid.New())

	pool.ExpectExec(sqlInsertSupply).WithArgs(item.ID, item.SupplierID, item.Code, item.Date, item.Price).WillReturnError(errors.New("exec error"))

	err := repo.Create(context.Background(), item)
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplyRepositoryUpdate(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplyRepository(t, pool)
	item := sampleSupply(uuid.New(), uuid.New())

	pool.ExpectExec(sqlUpdateSupply).WithArgs(item.ID, item.SupplierID, item.Code, item.Date, item.Price).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.Update(context.Background(), item)
	require.NoError(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplyRepositoryUpdateNoRows(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplyRepository(t, pool)
	item := sampleSupply(uuid.New(), uuid.New())

	pool.ExpectExec(sqlUpdateSupply).WithArgs(item.ID, item.SupplierID, item.Code, item.Date, item.Price).WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := repo.Update(context.Background(), item)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplyRepositoryDelete(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplyRepository(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlDeleteSupply).WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplyRepositoryDeleteNoRows(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplyRepository(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlDeleteSupply).WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err := repo.Delete(context.Background(), id)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, pool.ExpectationsWereMet())
}
