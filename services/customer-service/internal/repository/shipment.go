package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type ShipmentRepository struct {
	db     *pgxpool.Pool
	logger *logger.Logger
}

func NewShipmentRepository(db *pgxpool.Pool, logger *logger.Logger) *ShipmentRepository {
	return &ShipmentRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ShipmentRepository) GetList(ctx context.Context) ([]model.Shipment, error) {
	r.logger.Debug("getting shipment list")

	query := `SELECT id, customer_id, code, date, price FROM shipments ORDER BY date DESC;`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		r.logger.Error("failed to get shipment list", "error", err)
		return nil, err
	}

	defer rows.Close()

	items := make([]model.Shipment, 0)

	for rows.Next() {
		var item model.Shipment

		err := rows.Scan(&item.ID, &item.CustomerID, &item.Code, &item.Date, &item.Price)

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

func (r *ShipmentRepository) GetByID(ctx context.Context, id uuid.UUID) (model.Shipment, error) {
	return r.GetByIDTx(ctx, r.db, id)
}

func (r *ShipmentRepository) GetByIDTx(ctx context.Context, db DBExecutor, id uuid.UUID) (model.Shipment, error) {
	r.logger.Debug("getting shipment by id", "shipmentID", id)

	var item model.Shipment

	query := `SELECT id, customer_id, code, date, price FROM shipments WHERE id = $1;`

	err := db.QueryRow(ctx, query, id).
		Scan(&item.ID, &item.CustomerID, &item.Code, &item.Date, &item.Price)

	if err != nil {
		r.logger.Error("failed to get shipment by id", "shipmentID", id, "error", err)
		return model.Shipment{}, err
	}

	return item, nil
}

func (r *ShipmentRepository) Create(ctx context.Context, shipment model.Shipment) error {
	return r.CreateTx(ctx, r.db, shipment)
}

func (r *ShipmentRepository) CreateTx(ctx context.Context, db DBExecutor, shipment model.Shipment) error {
	r.logger.Debug("creating shipment", "shipmentID", shipment.ID, "code", shipment.Code)

	query := `INSERT INTO shipments (id, customer_id, code,date, price) VALUES ($1, $2, $3, $4, $5);`

	_, err := db.Exec(ctx, query,
		shipment.ID, shipment.CustomerID, shipment.Code, shipment.Date, shipment.Price)

	if err != nil {
		r.logger.Error("failed to create shipment", "shipmentID", shipment.ID, "error", err)
		return err
	}

	return nil
}

func (r *ShipmentRepository) Update(ctx context.Context, shipment model.Shipment) error {
	return r.UpdateTx(ctx, r.db, shipment)
}

func (r *ShipmentRepository) UpdateTx(ctx context.Context, db DBExecutor, shipment model.Shipment) error {
	r.logger.Debug("updating shipment", "shipmentID", shipment.ID)

	query := `
UPDATE shipments
SET
	customer_id = $2,
	code = $3,
	date = $4,
	price = $5
WHERE id = $1;
`

	commandTag, err := db.Exec(ctx, query,
		shipment.ID, shipment.CustomerID, shipment.Code, shipment.Date, shipment.Price)

	if err != nil {
		r.logger.Error("failed to update shipment", "shipmentID", shipment.ID, "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *ShipmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DeleteTx(ctx, r.db, id)
}

func (r *ShipmentRepository) DeleteTx(ctx context.Context, db DBExecutor, id uuid.UUID) error {
	r.logger.Debug("deleting shipment", "shipmentID", id)

	query := `DELETE FROM shipments WHERE id = $1;`

	commandTag, err := db.Exec(ctx, query, id)

	if err != nil {
		r.logger.Error("failed to delete shipment", "shipmentID", id, "error", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
