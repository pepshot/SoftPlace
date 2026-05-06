package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/helpers"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model/relations"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/repository"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type ShipmentService struct {
	shipmentRepo          ShipmentRepository
	shipmentGarnitureRepo ShipmentGarnitureRepository
	garnitureRepo         GarnitureRepository
	txManager             TransactionManager
	logger                *logger.Logger
}

func NewShipmentService(shipmentRepo ShipmentRepository, shipmentGarnitureRepo ShipmentGarnitureRepository,
	garnitureRepo GarnitureRepository, txManager TransactionManager, logger *logger.Logger) *ShipmentService {
	return &ShipmentService{
		shipmentRepo:          shipmentRepo,
		shipmentGarnitureRepo: shipmentGarnitureRepo,
		garnitureRepo:         garnitureRepo,
		txManager:             txManager,
		logger:                logger,
	}
}

func (s *ShipmentService) GetList(ctx context.Context) ([]dto.ShipmentResponse, error) {
	shipmentList, err := s.shipmentRepo.GetList(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.ToShipmentResponseList(shipmentList), nil
}

func (s *ShipmentService) GetByID(ctx context.Context, id uuid.UUID) (dto.ShipmentResponse, error) {
	item, err := s.shipmentRepo.GetByID(ctx, id)
	if err != nil {
		return dto.ShipmentResponse{}, ErrNotFound
	}

	return mapper.ToShipmentResponse(item), nil
}

func (s *ShipmentService) Create(ctx context.Context, customerID uuid.UUID, req dto.ShipmentRequest) (uuid.UUID, error) {
	if len(req.Garnitures) == 0 {
		return uuid.Nil, ErrEmptyComposition
	}

	shipmentDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return uuid.Nil, ErrInvalidDate
	}

	shipmentID := uuid.New()

	err = s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		relationItems := make([]relations.ShipmentGarniture, 0, len(req.Garnitures))
		priceItems := make([]helpers.PriceItem, 0, len(req.Garnitures))

		for _, item := range req.Garnitures {
			if item.Count <= 0 {
				return ErrInvalidCount
			}

			garnitureID, err := uuid.Parse(item.GarnitureID)
			if err != nil {
				return ErrInvalidID
			}

			garniture, err := s.garnitureRepo.GetByIDTx(ctx, tx, garnitureID)
			if err != nil {
				return ErrNotFound
			}

			if garniture.StockCount < item.Count {
				return ErrNotEnoughStock
			}

			if err := s.garnitureRepo.DecreaseStockTx(ctx, tx, garnitureID, item.Count); err != nil {
				return ErrNotEnoughStock
			}

			relationItems = append(relationItems, relations.ShipmentGarniture{
				ShipmentID:  shipmentID,
				GarnitureID: garnitureID,
				Count:       item.Count,
			})

			priceItems = append(priceItems, helpers.PriceItem{
				Price: garniture.Price,
				Count: item.Count,
			})
		}

		totalPrice := helpers.CalculateTotalPrice(priceItems)

		shipment := model.Shipment{
			ID:         shipmentID,
			CustomerID: customerID,
			Code:       req.Code,
			Date:       shipmentDate,
			Price:      totalPrice,
		}

		if err := s.shipmentRepo.CreateTx(ctx, tx, shipment); err != nil {
			return err
		}

		if err := s.shipmentGarnitureRepo.CreateManyTx(ctx, tx, relationItems); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("failed to create shipment", "error", err)
		return uuid.Nil, err
	}

	s.logger.Info("shipment created successfully", "shipmentID", shipmentID)

	return shipmentID, nil
}

func (s *ShipmentService) Update(ctx context.Context, id uuid.UUID, req dto.ShipmentRequest) error {
	if len(req.Garnitures) == 0 {
		return ErrEmptyComposition
	}

	shipmentDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return ErrInvalidDate
	}

	err = s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		oldShipment, err := s.shipmentRepo.GetByIDTx(ctx, tx, id)
		if err != nil {
			return ErrNotFound
		}

		oldRelations, err := s.shipmentGarnitureRepo.GetByShipmentIDTx(ctx, tx, id)
		if err != nil {
			return err
		}

		oldStockItems := make([]helpers.StockItem, 0, len(oldRelations))

		for _, item := range oldRelations {
			oldStockItems = append(oldStockItems, helpers.StockItem{
				ID:    item.GarnitureID,
				Count: item.Count,
			})
		}

		newStockItems := make([]helpers.StockItem, 0, len(req.Garnitures))
		newRelations := make([]relations.ShipmentGarniture, 0, len(req.Garnitures))
		priceItems := make([]helpers.PriceItem, 0, len(req.Garnitures))

		for _, item := range req.Garnitures {
			if item.Count <= 0 {
				return ErrInvalidCount
			}

			garnitureID, err := uuid.Parse(item.GarnitureID)
			if err != nil {
				return ErrInvalidID
			}

			garniture, err := s.garnitureRepo.GetByIDTx(ctx, tx, garnitureID)
			if err != nil {
				return ErrNotFound
			}

			newStockItems = append(newStockItems, helpers.StockItem{
				ID:    garnitureID,
				Count: item.Count,
			})

			newRelations = append(newRelations, relations.ShipmentGarniture{
				ShipmentID:  id,
				GarnitureID: garnitureID,
				Count:       item.Count,
			})

			priceItems = append(priceItems, helpers.PriceItem{
				Price: garniture.Price,
				Count: item.Count,
			})
		}

		toReserve, toRelease := helpers.CalculateStockDiff(oldStockItems, newStockItems)

		for _, item := range toReserve {
			garniture, err := s.garnitureRepo.GetByIDTx(ctx, tx, item.ID)
			if err != nil {
				return ErrNotFound
			}

			if garniture.StockCount < item.Count {
				return ErrNotEnoughStock
			}

			if err := s.garnitureRepo.DecreaseStockTx(ctx, tx, item.ID, item.Count); err != nil {
				return ErrNotEnoughStock
			}
		}

		for _, item := range toRelease {
			if err := s.garnitureRepo.IncreaseStockTx(ctx, tx, item.ID, item.Count); err != nil {
				return err
			}
		}

		totalPrice := helpers.CalculateTotalPrice(priceItems)

		updatedShipment := model.Shipment{
			ID:         id,
			CustomerID: oldShipment.CustomerID,
			Code:       req.Code,
			Date:       shipmentDate,
			Price:      totalPrice,
		}

		if err := s.shipmentRepo.UpdateTx(ctx, tx, updatedShipment); err != nil {
			return err
		}

		if err := s.shipmentGarnitureRepo.DeleteByShipmentIDTx(ctx, tx, id); err != nil {
			return err
		}

		if err := s.shipmentGarnitureRepo.CreateManyTx(ctx, tx, newRelations); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("failed to update shipment", "shipmentID", id, "error", err)
		return err
	}

	s.logger.Info("shipment updated successfully", "shipmentID", id)

	return nil
}

func (s *ShipmentService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		relationsItems, err := s.shipmentGarnitureRepo.GetByShipmentIDTx(ctx, tx, id)
		if err != nil {
			return err
		}

		for _, item := range relationsItems {
			if err := s.garnitureRepo.IncreaseStockTx(ctx, tx, item.GarnitureID, item.Count); err != nil {
				return err
			}
		}

		if err := s.shipmentGarnitureRepo.DeleteByShipmentIDTx(ctx, tx, id); err != nil {
			return err
		}

		if err := s.shipmentRepo.DeleteTx(ctx, tx, id); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("failed to delete shipment", "shipmentID", id, "error", err)
		return err
	}

	s.logger.Info("shipment deleted successfully", "shipmentID", id)

	return nil
}
