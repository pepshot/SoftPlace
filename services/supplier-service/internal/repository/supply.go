package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type SupplyRepository struct {
	db     DBExecutor
	logger *logger.Logger
}

func NewSupplyRepository(db DBExecutor, logger *logger.Logger) *SupplyRepository {
	return &SupplyRepository{db: db, logger: logger}
}

func (r *SupplyRepository) GetList(ctx context.Context) ([]models.Supply, error) {
	query := `SELECT id, supplier_id, code, date, price FROM supplies ORDER BY date DESC;`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Supply, 0)
	for rows.Next() {
		var item models.Supply
		if err := rows.Scan(&item.ID, &item.SupplierID, &item.Code, &item.Date, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *SupplyRepository) GetByID(ctx context.Context, id uuid.UUID) (models.Supply, error) {
	return r.GetByIDTx(ctx, r.db, id)
}

func (r *SupplyRepository) GetByIDTx(ctx context.Context, db DBExecutor, id uuid.UUID) (models.Supply, error) {
	var item models.Supply
	query := `SELECT id, supplier_id, code, date, price FROM supplies WHERE id = $1;`
	err := db.QueryRow(ctx, query, id).Scan(&item.ID, &item.SupplierID, &item.Code, &item.Date, &item.Price)
	if err != nil {
		return models.Supply{}, err
	}
	return item, nil
}

func (r *SupplyRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT id FROM supplies WHERE code = $1);`
	if err := r.db.QueryRow(ctx, query, code).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *SupplyRepository) Create(ctx context.Context, item models.Supply) error {
	return r.CreateTx(ctx, r.db, item)
}

func (r *SupplyRepository) CreateTx(ctx context.Context, db DBExecutor, item models.Supply) error {
	query := `INSERT INTO supplies (id, supplier_id, code, date, price) VALUES ($1, $2, $3, $4, $5);`
	_, err := db.Exec(ctx, query, item.ID, item.SupplierID, item.Code, item.Date, item.Price)
	return err
}

func (r *SupplyRepository) Update(ctx context.Context, item models.Supply) error {
	return r.UpdateTx(ctx, r.db, item)
}

func (r *SupplyRepository) UpdateTx(ctx context.Context, db DBExecutor, item models.Supply) error {
	query := `UPDATE supplies SET supplier_id = $2, code = $3, date = $4, price = $5 WHERE id = $1;`
	commandTag, err := db.Exec(ctx, query, item.ID, item.SupplierID, item.Code, item.Date, item.Price)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *SupplyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DeleteTx(ctx, r.db, id)
}

func (r *SupplyRepository) DeleteTx(ctx context.Context, db DBExecutor, id uuid.UUID) error {
	query := `DELETE FROM supplies WHERE id = $1;`
	commandTag, err := db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
