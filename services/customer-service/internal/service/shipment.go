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
	"github.com/pepshot/SoftPlace/shared/logger"
)

type ShipmentService struct {
	shipmentRepo          ShipmentRepository
	shipmentGarnitureRepo ShipmentGarnitureRepository
	garnitureRepo         GarnitureRepository
	logger                *logger.Logger
}

func NewShipmentService(
	shipmentRepo ShipmentRepository,
	shipmentGarnitureRepo ShipmentGarnitureRepository,
	garnitureRepo GarnitureRepository,
	logger *logger.Logger,
) *ShipmentService {
	return &ShipmentService{
		shipmentRepo:          shipmentRepo,
		shipmentGarnitureRepo: shipmentGarnitureRepo,
		garnitureRepo:         garnitureRepo,
		logger:                logger,
	}
}

func (s *ShipmentService) GetList(ctx context.Context) ([]dto.ShipmentResponse, error) {
	s.logger.Debug("getting shipment list")

	shipmentList, err := s.shipmentRepo.GetList(ctx)
	if err != nil {
		s.logger.Error("failed to get shipment list", "error", err)
		return nil, err
	}

	return mapper.ToShipmentResponseList(shipmentList), nil
}

func (s *ShipmentService) GetByID(ctx context.Context, id uuid.UUID) (dto.ShipmentResponse, error) {
	s.logger.Debug("getting shipment by id", "shipmentID", id)

	item, err := s.shipmentRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get shipment by id", "shipmentID", id, "error", err)
		return dto.ShipmentResponse{}, ErrNotFound
	}

	return mapper.ToShipmentResponse(item), nil
}

func (s *ShipmentService) Create(ctx context.Context, customerID uuid.UUID, req dto.ShipmentRequest) (uuid.UUID, error) {
	s.logger.Info("starting shipment creation", "code", req.Code, "customerID", customerID)

	if len(req.Garnitures) == 0 {
		return uuid.Nil, ErrEmptyComposition
	}

	shipmentDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return uuid.Nil, ErrInvalidDate
	}

	shipmentID := uuid.New()

	stockItems := make([]helpers.StockItem, 0, len(req.Garnitures))
	relationItems := make([]relations.ShipmentGarniture, 0, len(req.Garnitures))
	priceItems := make([]helpers.PriceItem, 0, len(req.Garnitures))

	for _, item := range req.Garnitures {
		if item.Count <= 0 {
			return uuid.Nil, ErrInvalidCount
		}

		garnitureID, err := uuid.Parse(item.GarnitureID)
		if err != nil {
			return uuid.Nil, ErrInvalidID
		}

		garniture, err := s.garnitureRepo.GetByID(ctx, garnitureID)
		if err != nil {
			s.logger.Error("failed to get garniture", "garnitureID", garnitureID, "error", err)
			return uuid.Nil, ErrNotFound
		}

		if garniture.StockCount < item.Count {
			s.logger.Warn(
				"not enough garniture stock",
				"garnitureID", garnitureID,
				"need", item.Count,
				"available", garniture.StockCount,
			)

			return uuid.Nil, ErrNotEnoughStock
		}

		stockItems = append(stockItems, helpers.StockItem{
			ID:    garnitureID,
			Count: item.Count,
		})

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

	reservedItems := make([]helpers.StockItem, 0, len(stockItems))

	for _, item := range stockItems {
		if err := s.garnitureRepo.DecreaseStock(ctx, item.ID, item.Count); err != nil {
			s.logger.Error(
				"failed to decrease garniture stock",
				"garnitureID", item.ID,
				"count", item.Count,
				"error", err,
			)

			s.releaseGarnitureStock(ctx, reservedItems)

			return uuid.Nil, err
		}

		reservedItems = append(reservedItems, item)
	}

	shipment := model.Shipment{
		ID:         shipmentID,
		CustomerID: customerID,
		Code:       req.Code,
		Date:       shipmentDate,
		Price:      totalPrice,
	}

	if err := s.shipmentRepo.Create(ctx, shipment); err != nil {
		s.logger.Error("failed to create shipment", "error", err)

		s.releaseGarnitureStock(ctx, reservedItems)

		return uuid.Nil, err
	}

	if err := s.shipmentGarnitureRepo.CreateMany(ctx, relationItems); err != nil {
		s.logger.Error("failed to create shipment-garniture relations", "error", err)

		s.releaseGarnitureStock(ctx, reservedItems)
		_ = s.shipmentRepo.Delete(ctx, shipmentID)

		return uuid.Nil, err
	}

	s.logger.Info(
		"shipment created successfully",
		"shipmentID", shipmentID,
		"price", totalPrice,
	)

	return shipmentID, nil
}

func (s *ShipmentService) Update(ctx context.Context, id uuid.UUID, req dto.ShipmentRequest) error {
	s.logger.Info("starting shipment update", "shipmentID", id)

	if len(req.Garnitures) == 0 {
		return ErrEmptyComposition
	}

	shipmentDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return ErrInvalidDate
	}

	oldShipment, err := s.shipmentRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get old shipment", "shipmentID", id, "error", err)
		return ErrNotFound
	}

	oldRelations, err := s.shipmentGarnitureRepo.GetByShipmentID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get old shipment-garniture relations", "shipmentID", id, "error", err)
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

		garniture, err := s.garnitureRepo.GetByID(ctx, garnitureID)
		if err != nil {
			s.logger.Error("failed to get garniture", "garnitureID", garnitureID, "error", err)
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
		garniture, err := s.garnitureRepo.GetByID(ctx, item.ID)
		if err != nil {
			s.logger.Error("failed to get garniture for stock checking", "garnitureID", item.ID, "error", err)
			return ErrNotFound
		}

		if garniture.StockCount < item.Count {
			s.logger.Warn(
				"not enough garniture stock for update",
				"garnitureID", item.ID,
				"need", item.Count,
				"available", garniture.StockCount,
			)

			return ErrNotEnoughStock
		}
	}

	reservedItems := make([]helpers.StockItem, 0, len(toReserve))

	for _, item := range toReserve {
		if err := s.garnitureRepo.DecreaseStock(ctx, item.ID, item.Count); err != nil {
			s.logger.Error(
				"failed to decrease garniture stock",
				"garnitureID", item.ID,
				"count", item.Count,
				"error", err,
			)

			s.releaseGarnitureStock(ctx, reservedItems)

			return err
		}

		reservedItems = append(reservedItems, item)
	}

	totalPrice := helpers.CalculateTotalPrice(priceItems)

	updatedShipment := model.Shipment{
		ID:         id,
		CustomerID: oldShipment.CustomerID,
		Code:       req.Code,
		Date:       shipmentDate,
		Price:      totalPrice,
	}

	if err := s.shipmentRepo.Update(ctx, updatedShipment); err != nil {
		s.logger.Error("failed to update shipment", "shipmentID", id, "error", err)

		s.releaseGarnitureStock(ctx, reservedItems)

		return err
	}

	if err := s.shipmentGarnitureRepo.DeleteByShipmentID(ctx, id); err != nil {
		s.logger.Error("failed to delete old shipment-garniture relations", "shipmentID", id, "error", err)

		s.releaseGarnitureStock(ctx, reservedItems)

		return err
	}

	if err := s.shipmentGarnitureRepo.CreateMany(ctx, newRelations); err != nil {
		s.logger.Error("failed to create new shipment-garniture relations", "shipmentID", id, "error", err)

		s.releaseGarnitureStock(ctx, reservedItems)

		return err
	}

	s.releaseGarnitureStock(ctx, toRelease)

	s.logger.Info(
		"shipment updated successfully",
		"shipmentID", id,
		"price", totalPrice,
	)

	return nil
}

func (s *ShipmentService) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting shipment", "shipmentID", id)

	relationsItems, err := s.shipmentGarnitureRepo.GetByShipmentID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get shipment-garniture relations", "shipmentID", id, "error", err)
		return err
	}

	stockItems := make([]helpers.StockItem, 0, len(relationsItems))

	for _, item := range relationsItems {
		stockItems = append(stockItems, helpers.StockItem{
			ID:    item.GarnitureID,
			Count: item.Count,
		})
	}

	if err := s.shipmentGarnitureRepo.DeleteByShipmentID(ctx, id); err != nil {
		s.logger.Error("failed to delete shipment-garniture relations", "shipmentID", id, "error", err)
		return err
	}

	if err := s.shipmentRepo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete shipment", "shipmentID", id, "error", err)
		return err
	}

	s.releaseGarnitureStock(ctx, stockItems)

	s.logger.Info("shipment deleted successfully", "shipmentID", id)

	return nil
}

func (s *ShipmentService) releaseGarnitureStock(ctx context.Context, items []helpers.StockItem) {
	for _, item := range items {
		if err := s.garnitureRepo.IncreaseStock(ctx, item.ID, item.Count); err != nil {
			s.logger.Error(
				"failed to release garniture stock",
				"garnitureID", item.ID,
				"count", item.Count,
				"error", err,
			)
		}
	}
}
