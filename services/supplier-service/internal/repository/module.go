package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type ModuleRepository struct {
	db     DBExecutor
	logger *logger.Logger
}

func NewModuleRepository(db DBExecutor, logger *logger.Logger) *ModuleRepository {
	return &ModuleRepository{db: db, logger: logger}
}

func (r *ModuleRepository) GetList(ctx context.Context) ([]models.Module, error) {
	r.logger.Debug("getting list of modules")

	query := `SELECT id, name, code, price, stock_count FROM modules ORDER BY name;`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.Error("failed to get module list", "error", err)
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Module, 0)
	for rows.Next() {
		var item models.Module
		if err := rows.Scan(&item.ID, &item.Name, &item.Code, &item.Price, &item.StockCount); err != nil {
			r.logger.Error("failed to scan module", "error", err)
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ModuleRepository) GetByID(ctx context.Context, id uuid.UUID) (models.Module, error) {
	return r.GetByIDTx(ctx, r.db, id)
}

func (r *ModuleRepository) GetByIDTx(ctx context.Context, db DBExecutor, id uuid.UUID) (models.Module, error) {
	var item models.Module

	query := `SELECT id, name, code, price, stock_count FROM modules WHERE id = $1;`
	if err := db.QueryRow(ctx, query, id).Scan(&item.ID, &item.Name, &item.Code, &item.Price, &item.StockCount); err != nil {
		r.logger.Error("failed to get module by id", "moduleID", id, "error", err)
		return models.Module{}, err
	}

	return item, nil
}

func (r *ModuleRepository) GetByCode(ctx context.Context, code string) (models.Module, error) {
	return r.GetByCodeTx(ctx, r.db, code)
}

func (r *ModuleRepository) GetByCodeTx(ctx context.Context, db DBExecutor, code string) (models.Module, error) {
	var item models.Module

	query := `SELECT id, name, code, price, stock_count FROM modules WHERE code = $1;`
	if err := db.QueryRow(ctx, query, code).Scan(&item.ID, &item.Name, &item.Code, &item.Price, &item.StockCount); err != nil {
		r.logger.Error("failed to get module by code", "code", code, "error", err)
		return models.Module{}, err
	}

	return item, nil
}

func (r *ModuleRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT id FROM modules WHERE code = $1);`
	if err := r.db.QueryRow(ctx, query, code).Scan(&exists); err != nil {
		r.logger.Error("failed to check module code", "code", code, "error", err)
		return false, err
	}

	return exists, nil
}

func (r *ModuleRepository) Create(ctx context.Context, module models.Module) error {
	return r.CreateTx(ctx, r.db, module)
}

func (r *ModuleRepository) CreateTx(ctx context.Context, db DBExecutor, module models.Module) error {
	r.logger.Debug("creating module", "moduleID", module.ID)

	query := `INSERT INTO modules (id, name, code, price, stock_count) VALUES ($1, $2, $3, $4, $5);`
	commandTag, err := db.Exec(ctx, query, module.ID, module.Name, module.Code, module.Price, module.StockCount)
	if err != nil {
		r.logger.Error("failed to create module", "moduleID", module.ID, "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *ModuleRepository) Update(ctx context.Context, module models.Module) error {
	return r.UpdateTx(ctx, r.db, module)
}

func (r *ModuleRepository) UpdateTx(ctx context.Context, db DBExecutor, module models.Module) error {
	r.logger.Debug("updating module", "moduleID", module.ID)

	query := `UPDATE modules SET name = $2, code = $3, price = $4, stock_count = $5 WHERE id = $1;`
	commandTag, err := db.Exec(ctx, query, module.ID, module.Name, module.Code, module.Price, module.StockCount)
	if err != nil {
		r.logger.Error("failed to update module", "moduleID", module.ID, "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *ModuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DeleteTx(ctx, r.db, id)
}

func (r *ModuleRepository) DeleteTx(ctx context.Context, db DBExecutor, id uuid.UUID) error {
	r.logger.Debug("deleting module", "moduleID", id)

	query := `DELETE FROM modules WHERE id = $1;`
	commandTag, err := db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete module", "moduleID", id, "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *ModuleRepository) IncreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	return r.IncreaseStockTx(ctx, r.db, id, count)
}

func (r *ModuleRepository) IncreaseStockTx(ctx context.Context, db DBExecutor, id uuid.UUID, count int) error {
	r.logger.Debug("increasing module stock", "moduleID", id, "count", count)

	query := `UPDATE modules SET stock_count = stock_count + $2 WHERE id = $1;`
	commandTag, err := db.Exec(ctx, query, id, count)
	if err != nil {
		r.logger.Error("failed to increase module stock", "moduleID", id, "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *ModuleRepository) DecreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	return r.DecreaseStockTx(ctx, r.db, id, count)
}

func (r *ModuleRepository) DecreaseStockTx(ctx context.Context, db DBExecutor, id uuid.UUID, count int) error {
	r.logger.Debug("decreasing module stock", "moduleID", id, "count", count)

	query := `UPDATE modules SET stock_count = stock_count - $2 WHERE id = $1 AND stock_count >= $2;`
	commandTag, err := db.Exec(ctx, query, id, count)
	if err != nil {
		r.logger.Error("failed to decrease module stock", "moduleID", id, "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
