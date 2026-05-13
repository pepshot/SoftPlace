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

func TestCustomerRepositoryCreate(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newRepository(t, mockPool)

	customer := newCustomer(uuid.New())

	mockPool.ExpectExec(sqlInsertCustomer).
		WithArgs(customer.ID, customer.Login, customer.Password, customer.Email).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.Create(context.Background(), customer)
	require.NoError(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestCustomerRepositoryCreateExecError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newRepository(t, mockPool)

	customer := newCustomer(uuid.New())

	mockPool.ExpectExec(sqlInsertCustomer).
		WithArgs(customer.ID, customer.Login, customer.Password, customer.Email).
		WillReturnError(errors.New("db error"))

	err := repo.Create(context.Background(), customer)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestCustomerRepositoryCreateNoRowsAffected(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newRepository(t, mockPool)

	customer := newCustomer(uuid.New())

	mockPool.ExpectExec(sqlInsertCustomer).
		WithArgs(customer.ID, customer.Login, customer.Password, customer.Email).
		WillReturnResult(pgxmock.NewResult("INSERT", 0))

	err := repo.Create(context.Background(), customer)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestCustomerRepositoryGetByID(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newRepository(t, mockPool)

	customer := newCustomer(uuid.New())

	mockPool.ExpectQuery(sqlSelectCustomer).
		WithArgs(customer.ID).
		WillReturnRows(customerRows(customer))

	result, err := repo.GetByID(context.Background(), customer.ID)
	require.NoError(t, err)
	require.Equal(t, customer.ID, result.ID)
	require.Equal(t, customer.Email, result.Email)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestCustomerRepositoryGetByIDNotFound(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newRepository(t, mockPool)

	customerID := uuid.New()

	mockPool.ExpectQuery(sqlSelectCustomer).
		WithArgs(customerID).
		WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetByID(context.Background(), customerID)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestCustomerRepositoryGetByIDQueryError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newRepository(t, mockPool)

	customerID := uuid.New()

	mockPool.ExpectQuery(sqlSelectCustomer).
		WithArgs(customerID).
		WillReturnError(errors.New("query error"))

	_, err := repo.GetByID(context.Background(), customerID)
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestCustomerRepositoryGetByLogin(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newRepository(t, mockPool)

	customer := newCustomer(uuid.New())

	mockPool.ExpectQuery(sqlSelectFromLogin).
		WithArgs(customer.Login).
		WillReturnRows(customerRows(customer))

	result, err := repo.GetByLogin(context.Background(), customer.Login)
	require.NoError(t, err)
	require.Equal(t, customer.ID, result.ID)
	require.Equal(t, customer.Email, result.Email)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestCustomerRepositoryGetByLoginNotFound(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newRepository(t, mockPool)

	mockPool.ExpectQuery(sqlSelectFromLogin).
		WithArgs("login").
		WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetByLogin(context.Background(), "login")
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestCustomerRepositoryGetByLoginQueryError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newRepository(t, mockPool)

	mockPool.ExpectQuery(sqlSelectFromLogin).
		WithArgs("login").
		WillReturnError(errors.New("query error"))

	_, err := repo.GetByLogin(context.Background(), "login")
	require.Error(t, err)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestCustomerRepositoryExistsByLoginTrue(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newRepository(t, mockPool)

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(true)
	mockPool.ExpectQuery(sqlExistFromLogin).
		WithArgs("login").
		WillReturnRows(rows)

	exists, err := repo.ExistsByLogin(context.Background(), "login")
	require.NoError(t, err)
	require.True(t, exists)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestCustomerRepositoryExistsByLoginFalse(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newRepository(t, mockPool)

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(false)
	mockPool.ExpectQuery(sqlExistFromLogin).
		WithArgs("login").
		WillReturnRows(rows)

	exists, err := repo.ExistsByLogin(context.Background(), "login")
	require.NoError(t, err)
	require.False(t, exists)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestCustomerRepositoryExistsByLoginError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newRepository(t, mockPool)

	mockPool.ExpectQuery(sqlExistFromLogin).
		WithArgs("login").
		WillReturnError(errors.New("query error"))

	exists, err := repo.ExistsByLogin(context.Background(), "login")
	require.Error(t, err)
	require.False(t, exists)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestCustomerRepositoryExistsByEmailTrue(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newRepository(t, mockPool)

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(true)
	mockPool.ExpectQuery(sqlExistFromEmail).
		WithArgs("user@example.com").
		WillReturnRows(rows)

	exists, err := repo.ExistsByEmail(context.Background(), "user@example.com")
	require.NoError(t, err)
	require.True(t, exists)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestCustomerRepositoryExistsByEmailFalse(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newRepository(t, mockPool)

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(false)
	mockPool.ExpectQuery(sqlExistFromEmail).
		WithArgs("user@example.com").
		WillReturnRows(rows)

	exists, err := repo.ExistsByEmail(context.Background(), "user@example.com")
	require.NoError(t, err)
	require.False(t, exists)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestCustomerRepositoryExistsByEmailError(t *testing.T) {
	mockPool := newMockPool(t)
	defer mockPool.Close()

	repo := newRepository(t, mockPool)

	mockPool.ExpectQuery(sqlExistFromEmail).
		WithArgs("user@example.com").
		WillReturnError(errors.New("query error"))

	exists, err := repo.ExistsByEmail(context.Background(), "user@example.com")
	require.Error(t, err)
	require.False(t, exists)
	require.NoError(t, mockPool.ExpectationsWereMet())
}
