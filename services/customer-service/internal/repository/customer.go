package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type CustomerRepository struct {
	db     DBExecutor
	logger *logger.Logger
}

func NewCustomerRepository(db DBExecutor, logger *logger.Logger) *CustomerRepository {
	return &CustomerRepository{
		db:     db,
		logger: logger,
	}
}

func (r *CustomerRepository) Create(ctx context.Context, customer model.Customer) error {
	r.logger.Debug("creating customer", "login", customer.Login)

	query := `INSERT INTO customers (id, login, password, email) VALUES ($1, $2, $3, $4);`

	commandTag, err := r.db.Exec(ctx, query,
		customer.ID, customer.Login, customer.Password, customer.Email)

	if err != nil {
		r.logger.Error("failed to create customer", "login", customer.Login, "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *CustomerRepository) GetByID(ctx context.Context, id uuid.UUID) (model.Customer, error) {
	r.logger.Debug("getting customer by id", "customerID", id)

	var customer model.Customer

	query := `SELECT id, login, password, email FROM customers WHERE id = $1;`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&customer.ID, &customer.Login, &customer.Password, &customer.Email)

	if err != nil {
		r.logger.Error("failed to get customer by id", "customerID", id, "error", err)
		return model.Customer{}, err
	}

	return customer, nil
}

func (r *CustomerRepository) GetByLogin(ctx context.Context, login string) (model.Customer, error) {
	r.logger.Debug("getting customer by login", "login", login)

	var customer model.Customer

	query := `SELECT id, login, password, email FROM customers WHERE login = $1;`

	err := r.db.QueryRow(ctx, query, login).
		Scan(&customer.ID, &customer.Login, &customer.Password, &customer.Email)

	if err != nil {
		r.logger.Error("failed to get customer by login", "login", login, "error", err)
		return model.Customer{}, err
	}

	return customer, nil
}

func (r *CustomerRepository) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT id FROM customers WHERE login = $1);`

	err := r.db.QueryRow(ctx, query, login).Scan(&exists)

	if err != nil {
		r.logger.Error("failed to check customer login", "login", login, "error", err)
		return false, err
	}

	return exists, nil
}

func (r *CustomerRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT id FROM customers WHERE email = $1);`

	err := r.db.QueryRow(ctx, query, email).Scan(&exists)

	if err != nil {
		r.logger.Error("failed to check customer email", "email", email, "error", err)
		return false, err
	}

	return exists, nil
}
