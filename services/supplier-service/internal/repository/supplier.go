package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
)

type SupplierRepository struct {
	db DBExecutor
}

func NewSupplierRepository(db DBExecutor) *SupplierRepository {
	return &SupplierRepository{db: db}
}

func (r *SupplierRepository) Create(ctx context.Context, supplier models.Supplier) error {
	query := `INSERT INTO suppliers (id, login, password, email) VALUES ($1, $2, $3, $4);`
	_, err := r.db.Exec(ctx, query, supplier.ID, supplier.Login, supplier.Password, supplier.Email)
	return err
}

func (r *SupplierRepository) GetByID(ctx context.Context, id uuid.UUID) (models.Supplier, error) {
	var item models.Supplier
	query := `SELECT id, login, password, email FROM suppliers WHERE id = $1;`
	err := r.db.QueryRow(ctx, query, id).Scan(&item.ID, &item.Login, &item.Password, &item.Email)
	if err != nil {
		return models.Supplier{}, err
	}
	return item, nil
}

func (r *SupplierRepository) GetByLogin(ctx context.Context, login string) (models.Supplier, error) {
	var item models.Supplier
	query := `SELECT id, login, password, email FROM suppliers WHERE login = $1;`
	err := r.db.QueryRow(ctx, query, login).Scan(&item.ID, &item.Login, &item.Password, &item.Email)
	if err != nil {
		return models.Supplier{}, err
	}
	return item, nil
}

func (r *SupplierRepository) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT id FROM suppliers WHERE login = $1);`
	err := r.db.QueryRow(ctx, query, login).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *SupplierRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT id FROM suppliers WHERE email = $1);`
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
