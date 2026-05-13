package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type GarnitureRepository struct {
	db     DBExecutor
	logger *logger.Logger
}

func NewGarnitureRepository(db DBExecutor, logger *logger.Logger) *GarnitureRepository {
	return &GarnitureRepository{
		db:     db,
		logger: logger,
	}
}

func (r *GarnitureRepository) GetList(ctx context.Context) ([]model.Garniture, error) {
	r.logger.Debug("getting garniture list")

	query := `SELECT id, name, code, price, stock_count FROM garnitures ORDER BY name;`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.Error("failed to get garniture list", "error", err)
		return nil, err
	}
	defer rows.Close()

	items := make([]model.Garniture, 0)

	for rows.Next() {
		var item model.Garniture

		err := rows.Scan(&item.ID, &item.Name, &item.Code, &item.Price, &item.StockCount)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *GarnitureRepository) GetByID(ctx context.Context, id uuid.UUID) (model.Garniture, error) {
	return r.GetByIDTx(ctx, r.db, id)
}

func (r *GarnitureRepository) GetByIDTx(ctx context.Context, db DBExecutor, id uuid.UUID) (model.Garniture, error) {
	r.logger.Debug("getting garniture by id", "garnitureID", id)

	var item model.Garniture

	query := `SELECT id, name, code, price, stock_count FROM garnitures WHERE id = $1;`

	err := db.QueryRow(ctx, query, id).
		Scan(&item.ID, &item.Name, &item.Code, &item.Price, &item.StockCount)

	if err != nil {
		r.logger.Error("failed to get garniture by id", "garnitureID", id, "error", err)
		return model.Garniture{}, err
	}

	return item, nil
}

func (r *GarnitureRepository) Create(ctx context.Context, garniture model.Garniture) error {
	return r.CreateTx(ctx, r.db, garniture)
}

func (r *GarnitureRepository) CreateTx(ctx context.Context, db DBExecutor, garniture model.Garniture) error {
	r.logger.Debug("creating garniture", "garnitureID", garniture.ID)

	query := `INSERT INTO garnitures (id, name, code, price, stock_count) VALUES ($1, $2, $3, $4, $5);`

	_, err := db.Exec(ctx, query,
		garniture.ID, garniture.Name, garniture.Code, garniture.Price, garniture.StockCount)

	if err != nil {
		r.logger.Error("failed to create garniture", "garnitureID", garniture.ID, "error", err)
		return err
	}

	return nil
}

func (r *GarnitureRepository) Update(ctx context.Context, garniture model.Garniture) error {
	return r.UpdateTx(ctx, r.db, garniture)
}

func (r *GarnitureRepository) UpdateTx(ctx context.Context, db DBExecutor, garniture model.Garniture) error {
	r.logger.Debug("updating garniture", "garnitureID", garniture.ID)

	query := `
UPDATE garnitures
SET
	name = $2,
	code = $3,
	price = $4,
	stock_count = $5
WHERE id = $1;`

	commandTag, err := db.Exec(ctx, query,
		garniture.ID, garniture.Name, garniture.Code, garniture.Price, garniture.StockCount)

	if err != nil {
		r.logger.Error("failed to update garniture", "garnitureID", garniture.ID, "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *GarnitureRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DeleteTx(ctx, r.db, id)
}

func (r *GarnitureRepository) DeleteTx(ctx context.Context, db DBExecutor, id uuid.UUID) error {
	r.logger.Debug("deleting garniture", "garnitureID", id)

	query := `DELETE FROM garnitures WHERE id = $1;`

	commandTag, err := db.Exec(ctx, query, id)

	if err != nil {
		r.logger.Error("failed to delete garniture", "garnitureID", id, "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *GarnitureRepository) IncreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	return r.IncreaseStockTx(ctx, r.db, id, count)
}

func (r *GarnitureRepository) IncreaseStockTx(ctx context.Context, db DBExecutor, id uuid.UUID, count int) error {
	r.logger.Debug("increasing garniture stock", "garnitureID", id, "count", count)

	query := `UPDATE garnitures SET stock_count = stock_count + $2 WHERE id = $1;`

	commandTag, err := db.Exec(ctx, query, id, count)

	if err != nil {
		r.logger.Error("failed to increase garniture stock", "garnitureID", id, "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *GarnitureRepository) DecreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	return r.DecreaseStockTx(ctx, r.db, id, count)
}

func (r *GarnitureRepository) DecreaseStockTx(ctx context.Context, db DBExecutor, id uuid.UUID, count int) error {
	r.logger.Debug("decreasing garniture stock", "garnitureID", id, "count", count)

	query := `UPDATE garnitures SET stock_count = stock_count - $2 WHERE id = $1 AND stock_count >= $2;`

	commandTag, err := db.Exec(ctx, query, id, count)

	if err != nil {
		r.logger.Error("failed to decrease garniture stock", "garnitureID", id, "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
