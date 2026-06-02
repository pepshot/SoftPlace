package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models/relations"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type ModuleMaterialRepository struct {
	db     DBExecutor
	logger *logger.Logger
}

func NewModuleMaterialRepository(db DBExecutor, logger *logger.Logger) *ModuleMaterialRepository {
	return &ModuleMaterialRepository{db: db, logger: logger}
}

func (r *ModuleMaterialRepository) CreateMany(ctx context.Context, items []relations.ModuleMaterial) error {
	return r.CreateManyTx(ctx, r.db, items)
}

func (r *ModuleMaterialRepository) CreateManyTx(ctx context.Context, db DBExecutor, items []relations.ModuleMaterial) error {
	query := `INSERT INTO module_materials (module_id, material_id, count) VALUES ($1, $2, $3);`
	for _, item := range items {
		if _, err := db.Exec(ctx, query, item.ModuleID, item.MaterialID, item.Count); err != nil {
			return err
		}
	}
	return nil
}

func (r *ModuleMaterialRepository) GetByModuleID(ctx context.Context, moduleID uuid.UUID) ([]relations.ModuleMaterial, error) {
	return r.GetByModuleIDTx(ctx, r.db, moduleID)
}

func (r *ModuleMaterialRepository) GetByModuleIDTx(ctx context.Context, db DBExecutor, moduleID uuid.UUID) ([]relations.ModuleMaterial, error) {
	query := `SELECT module_id, material_id, count FROM module_materials WHERE module_id = $1;`
	rows, err := db.Query(ctx, query, moduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]relations.ModuleMaterial, 0)
	for rows.Next() {
		var item relations.ModuleMaterial
		if err := rows.Scan(&item.ModuleID, &item.MaterialID, &item.Count); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ModuleMaterialRepository) DeleteByModuleID(ctx context.Context, moduleID uuid.UUID) error {
	return r.DeleteByModuleIDTx(ctx, r.db, moduleID)
}

func (r *ModuleMaterialRepository) DeleteByModuleIDTx(ctx context.Context, db DBExecutor, moduleID uuid.UUID) error {
	query := `DELETE FROM module_materials WHERE module_id = $1;`
	_, err := db.Exec(ctx, query, moduleID)
	return err
}
