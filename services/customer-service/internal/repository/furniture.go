package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type FurnitureRepository struct {
	db     *pgxpool.Pool
	logger *logger.Logger
}

func NewFurnitureRepository(db *pgxpool.Pool, logger *logger.Logger) FurnitureRepository {
	return FurnitureRepository{
		db:     db,
		logger: logger,
	}
}

func (r FurnitureRepository) GetList(ctx context.Context) ([]model.Furniture, error) {
	r.logger.Debug("getting list of furniture's")

	query := `SELECT id, name, code, price, stock_count FROM furniture ORDER BY name;`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		r.logger.Error("failed to get furniture list", "error", err)
		return nil, err
	}

	defer rows.Close()

	items := make([]model.Furniture, 0)

	for rows.Next() {
		var item model.Furniture

		err := rows.Scan(&item.ID, &item.Name, &item.Code, &item.Price, &item.StockCount)
		if err != nil {
			r.logger.Error("failed to scan furniture", "error", err)
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *FurnitureRepository) GetByID(ctx context.Context, id uuid.UUID) (model.Furniture, error) {
	return r.GetByIDTx(ctx, r.db, id)
}

func (r *FurnitureRepository) GetByIDTx(ctx context.Context, db DBExecutor, id uuid.UUID) (model.Furniture, error) {
	var item model.Furniture

	query := `SELECT id, name, code, price, stock_count	FROM furniture	WHERE id = $1;`

	err := db.QueryRow(ctx, query, id).
		Scan(&item.ID, &item.Name, &item.Code, &item.Price, &item.StockCount)

	if err != nil {
		r.logger.Error("failed to get furniture by id", "furnitureID", id, "error", err)
		return model.Furniture{}, err
	}

	return item, nil
}

func (r *FurnitureRepository) Create(ctx context.Context, furniture model.Furniture) error {
	return r.CreateTx(ctx, r.db, furniture)
}

func (r FurnitureRepository) CreateTx(ctx context.Context, db DBExecutor, furniture model.Furniture) error {
	r.logger.Debug("creating furniture", "furniture", furniture)

	query := `INSERT INTO furniture (id, name, code, price, stock_count) VALUES ($1, $2, $3, $4, $5);`

	commandTag, err := db.Exec(ctx, query,
		furniture.ID, furniture.Name, furniture.Code, furniture.Price, furniture.StockCount)

	if err != nil {
		r.logger.Error("failed to create furniture", "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *FurnitureRepository) Update(ctx context.Context, furniture model.Furniture) error {
	return r.UpdateTx(ctx, r.db, furniture)
}

func (r FurnitureRepository) UpdateTx(ctx context.Context, db DBExecutor, furniture model.Furniture) error {
	r.logger.Debug("updating furniture", "furnitureID", furniture.ID)

	query := `
UPDATE furniture
SET
	name = $2,
	code = $3,
	price = $4,
	stock_count = $5
WHERE id = $1;`

	commandTag, err := db.Exec(ctx, query,
		furniture.ID, furniture.Name, furniture.Code, furniture.Price, furniture.StockCount)

	if err != nil {
		r.logger.Error("failed to update furniture", "furnitureID", furniture.ID, "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *FurnitureRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DeleteTx(ctx, r.db, id)
}

func (r FurnitureRepository) DeleteTx(ctx context.Context, db DBExecutor, id uuid.UUID) error {
	r.logger.Debug("deleting furniture", "furnitureID", id)

	query := `DELETE FROM furniture WHERE id = $1;`

	commandTag, err := db.Exec(ctx, query, id)

	if err != nil {
		r.logger.Error("failed to delete furniture", "furnitureID", id, "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *FurnitureRepository) IncreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	return r.IncreaseStockTx(ctx, r.db, id, count)
}

func (r FurnitureRepository) IncreaseStockTx(ctx context.Context, db DBExecutor, id uuid.UUID, count int) error {
	r.logger.Debug("increasing furniture stock", "furnitureID", id, "count", count)

	query := `UPDATE furniture SET stock_count = stock_count + $2 WHERE id = $1;`

	commandTag, err := db.Exec(ctx, query, id, count)

	if err != nil {
		r.logger.Error("failed to increase furniture stock", "furnitureID", id, "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *FurnitureRepository) DecreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	return r.DecreaseStockTx(ctx, r.db, id, count)
}

func (r FurnitureRepository) DecreaseStockTx(ctx context.Context, db DBExecutor, id uuid.UUID, count int) error {
	r.logger.Debug("decreasing furniture stock", "furnitureID", id, "count", count)

	query := `UPDATE furniture SET stock_count = stock_count - $2 WHERE id = $1 AND stock_count >= $2;`

	commandTag, err := db.Exec(ctx, query, id, count)

	if err != nil {
		r.logger.Error("failed to decrease furniture stock", "furnitureID", id, "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
