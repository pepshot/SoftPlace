package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type MaterialRepository struct {
	db     DBExecutor
	logger *logger.Logger
}

func NewMaterialRepository(db DBExecutor, logger *logger.Logger) *MaterialRepository {
	return &MaterialRepository{db: db, logger: logger}
}

func (r *MaterialRepository) GetList(ctx context.Context) ([]models.Material, error) {
	query := `SELECT id, name, code, price, stock_count FROM materials ORDER BY name;`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Material, 0)
	for rows.Next() {
		var item models.Material
		if err := rows.Scan(&item.ID, &item.Name, &item.Code, &item.Price, &item.StockCount); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MaterialRepository) GetByID(ctx context.Context, id uuid.UUID) (models.Material, error) {
	return r.GetByIDTx(ctx, r.db, id)
}

func (r *MaterialRepository) GetByIDTx(ctx context.Context, db DBExecutor, id uuid.UUID) (models.Material, error) {
	var item models.Material
	query := `SELECT id, name, code, price, stock_count FROM materials WHERE id = $1;`
	err := db.QueryRow(ctx, query, id).Scan(&item.ID, &item.Name, &item.Code, &item.Price, &item.StockCount)
	if err != nil {
		return models.Material{}, err
	}
	return item, nil
}

func (r *MaterialRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT id FROM materials WHERE code = $1);`
	if err := r.db.QueryRow(ctx, query, code).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *MaterialRepository) Create(ctx context.Context, item models.Material) error {
	return r.CreateTx(ctx, r.db, item)
}

func (r *MaterialRepository) CreateTx(ctx context.Context, db DBExecutor, item models.Material) error {
	query := `INSERT INTO materials (id, name, code, price, stock_count) VALUES ($1, $2, $3, $4, $5);`
	commandTag, err := db.Exec(ctx, query, item.ID, item.Name, item.Code, item.Price, item.StockCount)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *MaterialRepository) Update(ctx context.Context, item models.Material) error {
	return r.UpdateTx(ctx, r.db, item)
}

func (r *MaterialRepository) UpdateTx(ctx context.Context, db DBExecutor, item models.Material) error {
	query := `UPDATE materials SET name = $2, code = $3, price = $4, stock_count = $5 WHERE id = $1;`
	commandTag, err := db.Exec(ctx, query, item.ID, item.Name, item.Code, item.Price, item.StockCount)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *MaterialRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DeleteTx(ctx, r.db, id)
}

func (r *MaterialRepository) DeleteTx(ctx context.Context, db DBExecutor, id uuid.UUID) error {
	query := `DELETE FROM materials WHERE id = $1;`
	commandTag, err := db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *MaterialRepository) IncreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	return r.IncreaseStockTx(ctx, r.db, id, count)
}

func (r *MaterialRepository) IncreaseStockTx(ctx context.Context, db DBExecutor, id uuid.UUID, count int) error {
	query := `UPDATE materials SET stock_count = stock_count + $2 WHERE id = $1;`
	commandTag, err := db.Exec(ctx, query, id, count)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *MaterialRepository) DecreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	return r.DecreaseStockTx(ctx, r.db, id, count)
}

func (r *MaterialRepository) DecreaseStockTx(ctx context.Context, db DBExecutor, id uuid.UUID, count int) error {
	query := `UPDATE materials SET stock_count = stock_count - $2 WHERE id = $1 AND stock_count >= $2;`
	commandTag, err := db.Exec(ctx, query, id, count)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
