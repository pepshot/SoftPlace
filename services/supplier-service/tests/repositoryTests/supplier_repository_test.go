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

func TestSupplierRepositoryCreate(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplierRepo(t, pool)
	supplier := sampleSupplier(uuid.New())

	pool.ExpectExec(sqlInsertSupplier).WithArgs(supplier.ID, supplier.Login, supplier.Password, supplier.Email).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.Create(context.Background(), supplier)
	require.NoError(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplierRepositoryCreateError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplierRepo(t, pool)
	supplier := sampleSupplier(uuid.New())

	pool.ExpectExec(sqlInsertSupplier).WithArgs(supplier.ID, supplier.Login, supplier.Password, supplier.Email).WillReturnError(errors.New("exec error"))

	err := repo.Create(context.Background(), supplier)
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplierRepositoryGetByIDNotFound(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplierRepo(t, pool)
	id := uuid.New()

	pool.ExpectQuery(sqlSelectSupplierByID).WithArgs(id).WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetByID(context.Background(), id)
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestSupplierRepositoryGetByID(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplierRepo(t, pool)
	supplier := sampleSupplier(uuid.New())

	pool.ExpectQuery(sqlSelectSupplierByID).WithArgs(supplier.ID).WillReturnRows(supplierRows(supplier))

	result, err := repo.GetByID(context.Background(), supplier.ID)
	require.NoError(t, err)
	require.Equal(t, supplier.ID, result.ID)
	require.Equal(t, supplier.Login, result.Login)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplierRepositoryGetByIDQueryError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplierRepo(t, pool)
	id := uuid.New()

	pool.ExpectQuery(sqlSelectSupplierByID).WithArgs(id).WillReturnError(errors.New("query error"))

	_, err := repo.GetByID(context.Background(), id)
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplierRepositoryGetByLogin(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplierRepo(t, pool)
	supplier := sampleSupplier(uuid.New())

	rows := pgxmock.NewRows([]string{"id", "login", "password", "email"}).AddRow(supplier.ID, supplier.Login, supplier.Password, supplier.Email)
	pool.ExpectQuery(sqlSelectSupplierByLogin).WithArgs(supplier.Login).WillReturnRows(rows)

	result, err := repo.GetByLogin(context.Background(), supplier.Login)
	require.NoError(t, err)
	require.Equal(t, supplier.ID, result.ID)
}

func TestSupplierRepositoryGetByLoginNotFound(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplierRepo(t, pool)

	pool.ExpectQuery(sqlSelectSupplierByLogin).WithArgs("supplier").WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetByLogin(context.Background(), "supplier")
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplierRepositoryGetByLoginQueryError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplierRepo(t, pool)

	pool.ExpectQuery(sqlSelectSupplierByLogin).WithArgs("supplier").WillReturnError(errors.New("query error"))

	_, err := repo.GetByLogin(context.Background(), "supplier")
	require.Error(t, err)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplierRepositoryExistsByEmailError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplierRepo(t, pool)

	pool.ExpectQuery(sqlExistsSupplierEmail).WithArgs("supplier@example.com").WillReturnError(errors.New("query error"))

	exists, err := repo.ExistsByEmail(context.Background(), "supplier@example.com")
	require.Error(t, err)
	require.False(t, exists)
}

func TestSupplierRepositoryExistsByLoginTrue(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplierRepo(t, pool)

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(true)
	pool.ExpectQuery(sqlExistsSupplierLogin).WithArgs("supplier").WillReturnRows(rows)

	exists, err := repo.ExistsByLogin(context.Background(), "supplier")
	require.NoError(t, err)
	require.True(t, exists)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplierRepositoryExistsByLoginFalse(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplierRepo(t, pool)

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(false)
	pool.ExpectQuery(sqlExistsSupplierLogin).WithArgs("supplier").WillReturnRows(rows)

	exists, err := repo.ExistsByLogin(context.Background(), "supplier")
	require.NoError(t, err)
	require.False(t, exists)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplierRepositoryExistsByLoginError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplierRepo(t, pool)

	pool.ExpectQuery(sqlExistsSupplierLogin).WithArgs("supplier").WillReturnError(errors.New("query error"))

	exists, err := repo.ExistsByLogin(context.Background(), "supplier")
	require.Error(t, err)
	require.False(t, exists)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplierRepositoryExistsByEmailTrue(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplierRepo(t, pool)

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(true)
	pool.ExpectQuery(sqlExistsSupplierEmail).WithArgs("supplier@example.com").WillReturnRows(rows)

	exists, err := repo.ExistsByEmail(context.Background(), "supplier@example.com")
	require.NoError(t, err)
	require.True(t, exists)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplierRepositoryExistsByEmailFalse(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplierRepo(t, pool)

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(false)
	pool.ExpectQuery(sqlExistsSupplierEmail).WithArgs("supplier@example.com").WillReturnRows(rows)

	exists, err := repo.ExistsByEmail(context.Background(), "supplier@example.com")
	require.NoError(t, err)
	require.False(t, exists)
	require.NoError(t, pool.ExpectationsWereMet())
}

func TestSupplierRepositoryExistsByEmailQueryError(t *testing.T) {
	pool := newMockPool(t)
	defer pool.Close()
	repo := newSupplierRepo(t, pool)

	pool.ExpectQuery(sqlExistsSupplierEmail).WithArgs("supplier@example.com").WillReturnError(errors.New("query error"))

	exists, err := repo.ExistsByEmail(context.Background(), "supplier@example.com")
	require.Error(t, err)
	require.False(t, exists)
	require.NoError(t, pool.ExpectationsWereMet())
}
