package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/model/relations"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type FurnitureModuleRepository struct {
	db     DBExecutor
	logger *logger.Logger
}

func NewFurnitureModuleRepository(db DBExecutor, logger *logger.Logger) *FurnitureModuleRepository {
	return &FurnitureModuleRepository{
		db:     db,
		logger: logger,
	}
}

func (r *FurnitureModuleRepository) CreateMany(ctx context.Context, items []relations.FurnitureModule) error {
	return r.CreateManyTx(ctx, r.db, items)
}

func (r *FurnitureModuleRepository) CreateManyTx(ctx context.Context, db DBExecutor, items []relations.FurnitureModule) error {
	r.logger.Debug("creating furniture-module relations", "count", len(items))

	query := `INSERT INTO furniture_modules (furniture_id, module_id, count) VALUES ($1, $2, $3);`

	for _, item := range items {
		_, err := db.Exec(ctx, query, item.FurnitureID, item.ModuleID, item.Count)

		if err != nil {
			r.logger.Error(
				"failed to create furniture-module relation",
				"furnitureID", item.FurnitureID,
				"moduleID", item.ModuleID,
				"error", err,
			)
			return err
		}
	}

	return nil
}

func (r *FurnitureModuleRepository) GetByFurnitureID(ctx context.Context,
	furnitureID uuid.UUID) ([]relations.FurnitureModule, error) {
	return r.GetByFurnitureIDTx(ctx, r.db, furnitureID)
}

func (r *FurnitureModuleRepository) GetByFurnitureIDTx(ctx context.Context, db DBExecutor,
	furnitureID uuid.UUID) ([]relations.FurnitureModule, error) {
	r.logger.Debug("getting furniture-module relations", "furnitureID", furnitureID)

	query := `SELECT furniture_id, module_id, count FROM furniture_modules WHERE furniture_id = $1;`

	rows, err := db.Query(ctx, query, furnitureID)

	if err != nil {
		r.logger.Error("failed to get furniture-module relations", "furnitureID", furnitureID, "error", err)
		return nil, err
	}

	defer rows.Close()

	items := make([]relations.FurnitureModule, 0)

	for rows.Next() {
		var item relations.FurnitureModule

		err := rows.Scan(&item.FurnitureID, &item.ModuleID, &item.Count)

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

func (r *FurnitureModuleRepository) DeleteByFurnitureID(ctx context.Context, furnitureID uuid.UUID) error {
	return r.DeleteByFurnitureIDTx(ctx, r.db, furnitureID)
}

func (r *FurnitureModuleRepository) DeleteByFurnitureIDTx(ctx context.Context, db DBExecutor, furnitureID uuid.UUID) error {
	r.logger.Debug("deleting furniture-module relations", "furnitureID", furnitureID)

	query := `DELETE FROM furniture_modules WHERE furniture_id = $1;`

	_, err := db.Exec(ctx, query, furnitureID)

	if err != nil {
		r.logger.Error("failed to delete furniture-module relations", "furnitureID", furnitureID, "error", err)
		return err
	}

	return nil
}
