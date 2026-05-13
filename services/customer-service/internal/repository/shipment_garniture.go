package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/model/relations"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type ShipmentGarnitureRepository struct {
	db     DBExecutor
	logger *logger.Logger
}

func NewShipmentGarnitureRepository(db DBExecutor, logger *logger.Logger) *ShipmentGarnitureRepository {
	return &ShipmentGarnitureRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ShipmentGarnitureRepository) CreateMany(ctx context.Context, items []relations.ShipmentGarniture) error {
	return r.CreateManyTx(ctx, r.db, items)
}

func (r *ShipmentGarnitureRepository) CreateManyTx(ctx context.Context, db DBExecutor, items []relations.ShipmentGarniture) error {
	r.logger.Debug("creating shipment-garniture relations", "count", len(items))

	query := `INSERT INTO shipment_garnitures (shipment_id, garniture_id, count) VALUES ($1, $2, $3);`

	for _, item := range items {
		_, err := db.Exec(ctx, query, item.ShipmentID, item.GarnitureID, item.Count)

		if err != nil {
			r.logger.Error(
				"failed to create shipment-garniture relation",
				"shipmentID", item.ShipmentID,
				"garnitureID", item.GarnitureID,
				"error", err,
			)
			return err
		}
	}

	return nil
}

func (r *ShipmentGarnitureRepository) GetByShipmentID(ctx context.Context,
	shipmentID uuid.UUID) ([]relations.ShipmentGarniture, error) {
	return r.GetByShipmentIDTx(ctx, r.db, shipmentID)
}

func (r *ShipmentGarnitureRepository) GetByShipmentIDTx(ctx context.Context, db DBExecutor, shipmentID uuid.UUID) ([]relations.ShipmentGarniture, error) {
	r.logger.Debug("getting shipment-garniture relations", "shipmentID", shipmentID)

	query := `SELECT shipment_id, garniture_id, count FROM shipment_garnitures WHERE shipment_id = $1;`

	rows, err := db.Query(ctx, query, shipmentID)

	if err != nil {
		r.logger.Error("failed to get shipment-garniture relations", "shipmentID", shipmentID, "error", err)
		return nil, err
	}
	defer rows.Close()

	items := make([]relations.ShipmentGarniture, 0)

	for rows.Next() {
		var item relations.ShipmentGarniture

		err := rows.Scan(&item.ShipmentID, &item.GarnitureID, &item.Count)
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

func (r *ShipmentGarnitureRepository) DeleteByShipmentID(ctx context.Context, shipmentID uuid.UUID) error {
	return r.DeleteByShipmentIDTx(ctx, r.db, shipmentID)
}

func (r *ShipmentGarnitureRepository) DeleteByShipmentIDTx(ctx context.Context, db DBExecutor, shipmentID uuid.UUID) error {
	r.logger.Debug("deleting shipment-garniture relations", "shipmentID", shipmentID)

	query := `DELETE FROM shipment_garnitures WHERE shipment_id = $1;`

	_, err := db.Exec(ctx, query, shipmentID)

	if err != nil {
		r.logger.Error("failed to delete shipment-garniture relations", "shipmentID", shipmentID, "error", err)
		return err
	}

	return nil
}
