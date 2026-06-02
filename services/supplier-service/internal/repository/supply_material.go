package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models/relations"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type SupplyMaterialRepository struct {
	db     DBExecutor
	logger *logger.Logger
}

func NewSupplyMaterialRepository(db DBExecutor, logger *logger.Logger) *SupplyMaterialRepository {
	return &SupplyMaterialRepository{db: db, logger: logger}
}

func (r *SupplyMaterialRepository) CreateMany(ctx context.Context, items []relations.SupplyMaterial) error {
	return r.CreateManyTx(ctx, r.db, items)
}

func (r *SupplyMaterialRepository) CreateManyTx(ctx context.Context, db DBExecutor, items []relations.SupplyMaterial) error {
	query := `INSERT INTO supply_materials (supply_id, material_id, count) VALUES ($1, $2, $3);`
	for _, item := range items {
		if _, err := db.Exec(ctx, query, item.SupplyID, item.MaterialID, item.Count); err != nil {
			return err
		}
	}
	return nil
}

func (r *SupplyMaterialRepository) GetBySupplyID(ctx context.Context, supplyID uuid.UUID) ([]relations.SupplyMaterial, error) {
	return r.GetBySupplyIDTx(ctx, r.db, supplyID)
}

func (r *SupplyMaterialRepository) GetBySupplyIDTx(ctx context.Context, db DBExecutor, supplyID uuid.UUID) ([]relations.SupplyMaterial, error) {
	query := `SELECT supply_id, material_id, count FROM supply_materials WHERE supply_id = $1;`
	rows, err := db.Query(ctx, query, supplyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]relations.SupplyMaterial, 0)
	for rows.Next() {
		var item relations.SupplyMaterial
		if err := rows.Scan(&item.SupplyID, &item.MaterialID, &item.Count); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *SupplyMaterialRepository) DeleteBySupplyID(ctx context.Context, supplyID uuid.UUID) error {
	return r.DeleteBySupplyIDTx(ctx, r.db, supplyID)
}

func (r *SupplyMaterialRepository) DeleteBySupplyIDTx(ctx context.Context, db DBExecutor, supplyID uuid.UUID) error {
	query := `DELETE FROM supply_materials WHERE supply_id = $1;`
	_, err := db.Exec(ctx, query, supplyID)
	return err
}
