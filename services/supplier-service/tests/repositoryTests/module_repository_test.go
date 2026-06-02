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

func TestModuleRepositoryGetList(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)
	item := sampleModule(uuid.New())

	pool.ExpectQuery(sqlSelectModules).WillReturnRows(moduleRows(item))

	result, err := repo.GetList(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, item.ID, result[0].ID)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryGetListQueryError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)

	pool.ExpectQuery(sqlSelectModules).WillReturnError(errors.New("query error"))

	_, err := repo.GetList(context.Background())
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryGetByID(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)
	item := sampleModule(uuid.New())

	pool.ExpectQuery(sqlSelectModuleByID).WithArgs(item.ID).WillReturnRows(moduleRows(item))

	result, err := repo.GetByID(context.Background(), item.ID)
	require.NoError(t, err)
	require.Equal(t, item.ID, result.ID)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryGetByIDNotFound(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)
	id := uuid.New()

	pool.ExpectQuery(sqlSelectModuleByID).WithArgs(id).WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetByID(context.Background(), id)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryGetByIDQueryError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)
	id := uuid.New()

	pool.ExpectQuery(sqlSelectModuleByID).WithArgs(id).WillReturnError(errors.New("query error"))

	_, err := repo.GetByID(context.Background(), id)
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryGetByCode(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)
	item := sampleModule(uuid.New())

	pool.ExpectQuery(sqlSelectModuleByCode).WithArgs(item.Code).WillReturnRows(moduleRows(item))

	result, err := repo.GetByCode(context.Background(), item.Code)
	require.NoError(t, err)
	require.Equal(t, item.ID, result.ID)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryGetByCodeNotFound(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)

	pool.ExpectQuery(sqlSelectModuleByCode).WithArgs("MOD-1").WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetByCode(context.Background(), "MOD-1")
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryGetByCodeQueryError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)

	pool.ExpectQuery(sqlSelectModuleByCode).WithArgs("MOD-1").WillReturnError(errors.New("query error"))

	_, err := repo.GetByCode(context.Background(), "MOD-1")
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryExistsByCodeTrueFalseError(t *testing.T) {
	tests := []struct {
		name   string
		result any
		err    error
		exists bool
	}{
		{name: "true", result: true, exists: true},
		{name: "false", result: false, exists: false},
		{name: "error", err: errors.New("query error"), exists: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := newMockPool(t)
			defer pool.Close()
			repo := newModuleRepository(t, pool)

			if tt.err != nil {
				pool.ExpectQuery(sqlExistsModuleCode).WithArgs("MOD-1").WillReturnError(tt.err)
			} else {
				rows := pgxmock.NewRows([]string{"exists"}).AddRow(tt.result)
				pool.ExpectQuery(sqlExistsModuleCode).WithArgs("MOD-1").WillReturnRows(rows)
			}

			exists, err := repo.ExistsByCode(context.Background(), "MOD-1")
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

func TestModuleRepositoryCreate(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)
	item := sampleModule(uuid.New())

	pool.ExpectExec(sqlInsertModule).WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.Create(context.Background(), item)
	require.NoError(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryCreateExecError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)
	item := sampleModule(uuid.New())

	pool.ExpectExec(sqlInsertModule).WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).WillReturnError(errors.New("exec error"))

	err := repo.Create(context.Background(), item)
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryCreateNoRows(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)
	item := sampleModule(uuid.New())

	pool.ExpectExec(sqlInsertModule).WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).WillReturnResult(pgxmock.NewResult("INSERT", 0))

	err := repo.Create(context.Background(), item)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryUpdate(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)
	item := sampleModule(uuid.New())

	pool.ExpectExec(sqlUpdateModule).WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.Update(context.Background(), item)
	require.NoError(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryUpdateNoRows(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)
	item := sampleModule(uuid.New())

	pool.ExpectExec(sqlUpdateModule).WithArgs(item.ID, item.Name, item.Code, item.Price, item.StockCount).WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := repo.Update(context.Background(), item)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryDelete(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlDeleteModule).WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryDeleteNoRows(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlDeleteModule).WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err := repo.Delete(context.Background(), id)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryIncreaseStock(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlIncreaseModule).WithArgs(id, 2).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.IncreaseStock(context.Background(), id, 2)
	require.NoError(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryIncreaseStockNoRows(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlIncreaseModule).WithArgs(id, 2).WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := repo.IncreaseStock(context.Background(), id, 2)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryDecreaseStock(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlDecreaseModule).WithArgs(id, 1).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.DecreaseStock(context.Background(), id, 1)
	require.NoError(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestModuleRepositoryDecreaseStockNoRows(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newModuleRepository(t, pool)
	id := uuid.New()

	pool.ExpectExec(sqlDecreaseModule).WithArgs(id, 1).WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := repo.DecreaseStock(context.Background(), id, 1)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, pool.ExpectationsWereMet())
}
