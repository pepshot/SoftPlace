package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/model/relations"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type GarnitureFurnitureRepository struct {
	db     DBExecutor
	logger *logger.Logger
}

func NewGarnitureFurnitureRepository(db DBExecutor, logger *logger.Logger) *GarnitureFurnitureRepository {
	return &GarnitureFurnitureRepository{
		db:     db,
		logger: logger,
	}
}

func (r *GarnitureFurnitureRepository) CreateMany(ctx context.Context, items []relations.GarnitureFurniture) error {
	return r.CreateManyTx(ctx, r.db, items)
}

func (r *GarnitureFurnitureRepository) CreateManyTx(ctx context.Context, db DBExecutor, items []relations.GarnitureFurniture) error {
	r.logger.Debug("creating garniture-furniture relations", "count", len(items))

	query := `INSERT INTO garniture_furniture (garniture_id, furniture_id, count) VALUES ($1, $2, $3);`

	for _, item := range items {
		_, err := db.Exec(ctx, query,
			item.GarnitureID, item.FurnitureID, item.Count)

		if err != nil {
			r.logger.Error(
				"failed to create garniture-furniture relation",
				"garnitureID", item.GarnitureID,
				"furnitureID", item.FurnitureID,
				"error", err,
			)
			return err
		}
	}

	return nil
}

func (r *GarnitureFurnitureRepository) GetByGarnitureID(ctx context.Context,
	garnitureID uuid.UUID) ([]relations.GarnitureFurniture, error) {
	return r.GetByGarnitureIDTx(ctx, r.db, garnitureID)
}

func (r *GarnitureFurnitureRepository) GetByGarnitureIDTx(ctx context.Context, db DBExecutor, garnitureID uuid.UUID) ([]relations.GarnitureFurniture, error) {
	r.logger.Debug("getting garniture-furniture relations", "garnitureID", garnitureID)

	query := `SELECT garniture_id, furniture_id, count FROM garniture_furniture WHERE garniture_id = $1;`

	rows, err := db.Query(ctx, query, garnitureID)
	if err != nil {
		r.logger.Error("failed to get garniture-furniture relations", "garnitureID", garnitureID, "error", err)
		return nil, err
	}
	defer rows.Close()

	items := make([]relations.GarnitureFurniture, 0)

	for rows.Next() {
		var item relations.GarnitureFurniture

		err := rows.Scan(&item.GarnitureID, &item.FurnitureID, &item.Count)
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

func (r *GarnitureFurnitureRepository) DeleteByGarnitureID(ctx context.Context, garnitureID uuid.UUID) error {
	return r.DeleteByGarnitureIDTx(ctx, r.db, garnitureID)
}

func (r *GarnitureFurnitureRepository) DeleteByGarnitureIDTx(ctx context.Context, db DBExecutor, garnitureID uuid.UUID) error {
	r.logger.Debug("deleting garniture-furniture relations", "garnitureID", garnitureID)

	query := `DELETE FROM garniture_furniture WHERE garniture_id = $1;`

	_, err := db.Exec(ctx, query, garnitureID)

	if err != nil {
		r.logger.Error("failed to delete garniture-furniture relations", "garnitureID", garnitureID, "error", err)
		return err
	}

	return nil
}
