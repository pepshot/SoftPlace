package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/helpers"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model/relations"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type GarnitureService struct {
	garnitureRepo          GarnitureRepository
	garnitureFurnitureRepo GarnitureFurnitureRepository
	furnitureRepo          FurnitureRepository
	logger                 *logger.Logger
}

func NewGarnitureService(
	garnitureRepo GarnitureRepository,
	garnitureFurnitureRepo GarnitureFurnitureRepository,
	furnitureRepo FurnitureRepository,
	logger *logger.Logger,
) *GarnitureService {
	return &GarnitureService{
		garnitureRepo:          garnitureRepo,
		garnitureFurnitureRepo: garnitureFurnitureRepo,
		furnitureRepo:          furnitureRepo,
		logger:                 logger,
	}
}

func (s *GarnitureService) GetList(ctx context.Context) ([]dto.GarnitureResponse, error) {
	s.logger.Debug("getting garniture list")

	garnitureList, err := s.garnitureRepo.GetList(ctx)
	if err != nil {
		s.logger.Error("failed to get garniture list", "error", err)
		return nil, err
	}

	return mapper.ToGarnitureResponseList(garnitureList), nil
}

func (s *GarnitureService) GetByID(ctx context.Context, id uuid.UUID) (dto.GarnitureResponse, error) {
	s.logger.Debug("getting garniture by id", "garnitureID", id)

	item, err := s.garnitureRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get garniture by id", "garnitureID", id, "error", err)
		return dto.GarnitureResponse{}, ErrNotFound
	}

	return mapper.ToGarnitureResponse(item), nil
}

func (s *GarnitureService) Create(ctx context.Context, req dto.GarnitureRequest) (uuid.UUID, error) {
	s.logger.Info("starting garniture creation", "name", req.Name, "code", req.Code)

	if len(req.Furniture) == 0 {
		return uuid.Nil, ErrEmptyComposition
	}

	garnitureID := uuid.New()

	stockItems := make([]helpers.StockItem, 0, len(req.Furniture))
	relationItems := make([]relations.GarnitureFurniture, 0, len(req.Furniture))
	priceItems := make([]helpers.PriceItem, 0, len(req.Furniture))

	for _, item := range req.Furniture {
		if item.Count <= 0 {
			return uuid.Nil, ErrInvalidCount
		}

		furnitureID, err := uuid.Parse(item.FurnitureID)
		if err != nil {
			return uuid.Nil, ErrInvalidID
		}

		furniture, err := s.furnitureRepo.GetByID(ctx, furnitureID)
		if err != nil {
			s.logger.Error("failed to get furniture", "furnitureID", furnitureID, "error", err)
			return uuid.Nil, ErrNotFound
		}

		if furniture.StockCount < item.Count {
			s.logger.Warn(
				"not enough furniture stock",
				"furnitureID", furnitureID,
				"need", item.Count,
				"available", furniture.StockCount,
			)

			return uuid.Nil, ErrNotEnoughStock
		}

		stockItems = append(stockItems, helpers.StockItem{
			ID:    furnitureID,
			Count: item.Count,
		})

		relationItems = append(relationItems, relations.GarnitureFurniture{
			GarnitureID: garnitureID,
			FurnitureID: furnitureID,
			Count:       item.Count,
		})

		priceItems = append(priceItems, helpers.PriceItem{
			Price: furniture.Price,
			Count: item.Count,
		})
	}

	totalPrice := helpers.CalculateTotalPrice(priceItems)

	reservedItems := make([]helpers.StockItem, 0, len(stockItems))

	for _, item := range stockItems {
		if err := s.furnitureRepo.DecreaseStock(ctx, item.ID, item.Count); err != nil {
			s.logger.Error(
				"failed to decrease furniture stock",
				"furnitureID", item.ID,
				"count", item.Count,
				"error", err,
			)

			s.releaseFurnitureStock(ctx, reservedItems)

			return uuid.Nil, err
		}

		reservedItems = append(reservedItems, item)
	}

	garniture := model.Garniture{
		ID:         garnitureID,
		Name:       req.Name,
		Code:       req.Code,
		Price:      totalPrice,
		StockCount: 1,
	}

	if err := s.garnitureRepo.Create(ctx, garniture); err != nil {
		s.logger.Error("failed to create garniture", "error", err)

		s.releaseFurnitureStock(ctx, reservedItems)

		return uuid.Nil, err
	}

	if err := s.garnitureFurnitureRepo.CreateMany(ctx, relationItems); err != nil {
		s.logger.Error("failed to create garniture-furniture relations", "error", err)

		s.releaseFurnitureStock(ctx, reservedItems)
		_ = s.garnitureRepo.Delete(ctx, garnitureID)

		return uuid.Nil, err
	}

	s.logger.Info(
		"garniture created successfully",
		"garnitureID", garnitureID,
		"price", totalPrice,
	)

	return garnitureID, nil
}

func (s *GarnitureService) Update(ctx context.Context, id uuid.UUID, req dto.GarnitureRequest) error {
	s.logger.Info("starting garniture update", "garnitureID", id)

	if len(req.Furniture) == 0 {
		return ErrEmptyComposition
	}

	oldGarniture, err := s.garnitureRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get old garniture", "garnitureID", id, "error", err)
		return ErrNotFound
	}

	oldRelations, err := s.garnitureFurnitureRepo.GetByGarnitureID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get old garniture-furniture relations", "garnitureID", id, "error", err)
		return err
	}

	oldStockItems := make([]helpers.StockItem, 0, len(oldRelations))

	for _, item := range oldRelations {
		oldStockItems = append(oldStockItems, helpers.StockItem{
			ID:    item.FurnitureID,
			Count: item.Count,
		})
	}

	newStockItems := make([]helpers.StockItem, 0, len(req.Furniture))
	newRelations := make([]relations.GarnitureFurniture, 0, len(req.Furniture))
	priceItems := make([]helpers.PriceItem, 0, len(req.Furniture))

	for _, item := range req.Furniture {
		if item.Count <= 0 {
			return ErrInvalidCount
		}

		furnitureID, err := uuid.Parse(item.FurnitureID)
		if err != nil {
			return ErrInvalidID
		}

		furniture, err := s.furnitureRepo.GetByID(ctx, furnitureID)
		if err != nil {
			s.logger.Error("failed to get furniture", "furnitureID", furnitureID, "error", err)
			return ErrNotFound
		}

		newStockItems = append(newStockItems, helpers.StockItem{
			ID:    furnitureID,
			Count: item.Count,
		})

		newRelations = append(newRelations, relations.GarnitureFurniture{
			GarnitureID: id,
			FurnitureID: furnitureID,
			Count:       item.Count,
		})

		priceItems = append(priceItems, helpers.PriceItem{
			Price: furniture.Price,
			Count: item.Count,
		})
	}

	toReserve, toRelease := helpers.CalculateStockDiff(oldStockItems, newStockItems)

	for _, item := range toReserve {
		furniture, err := s.furnitureRepo.GetByID(ctx, item.ID)
		if err != nil {
			s.logger.Error("failed to get furniture for stock checking", "furnitureID", item.ID, "error", err)
			return ErrNotFound
		}

		if furniture.StockCount < item.Count {
			s.logger.Warn(
				"not enough furniture stock for update",
				"furnitureID", item.ID,
				"need", item.Count,
				"available", furniture.StockCount,
			)

			return ErrNotEnoughStock
		}
	}

	reservedItems := make([]helpers.StockItem, 0, len(toReserve))

	for _, item := range toReserve {
		if err := s.furnitureRepo.DecreaseStock(ctx, item.ID, item.Count); err != nil {
			s.logger.Error(
				"failed to decrease furniture stock",
				"furnitureID", item.ID,
				"count", item.Count,
				"error", err,
			)

			s.releaseFurnitureStock(ctx, reservedItems)

			return err
		}

		reservedItems = append(reservedItems, item)
	}

	totalPrice := helpers.CalculateTotalPrice(priceItems)

	updatedGarniture := model.Garniture{
		ID:         id,
		Name:       req.Name,
		Code:       req.Code,
		Price:      totalPrice,
		StockCount: oldGarniture.StockCount,
	}

	if err := s.garnitureRepo.Update(ctx, updatedGarniture); err != nil {
		s.logger.Error("failed to update garniture", "garnitureID", id, "error", err)

		s.releaseFurnitureStock(ctx, reservedItems)

		return err
	}

	if err := s.garnitureFurnitureRepo.DeleteByGarnitureID(ctx, id); err != nil {
		s.logger.Error("failed to delete old garniture-furniture relations", "garnitureID", id, "error", err)

		s.releaseFurnitureStock(ctx, reservedItems)

		return err
	}

	if err := s.garnitureFurnitureRepo.CreateMany(ctx, newRelations); err != nil {
		s.logger.Error("failed to create new garniture-furniture relations", "garnitureID", id, "error", err)

		s.releaseFurnitureStock(ctx, reservedItems)

		return err
	}

	s.releaseFurnitureStock(ctx, toRelease)

	s.logger.Info(
		"garniture updated successfully",
		"garnitureID", id,
		"price", totalPrice,
	)

	return nil
}

func (s *GarnitureService) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting garniture", "garnitureID", id)

	relationsItems, err := s.garnitureFurnitureRepo.GetByGarnitureID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get garniture-furniture relations", "garnitureID", id, "error", err)
		return err
	}

	stockItems := make([]helpers.StockItem, 0, len(relationsItems))

	for _, item := range relationsItems {
		stockItems = append(stockItems, helpers.StockItem{
			ID:    item.FurnitureID,
			Count: item.Count,
		})
	}

	if err := s.garnitureFurnitureRepo.DeleteByGarnitureID(ctx, id); err != nil {
		s.logger.Error("failed to delete garniture-furniture relations", "garnitureID", id, "error", err)
		return err
	}

	if err := s.garnitureRepo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete garniture", "garnitureID", id, "error", err)
		return err
	}

	s.releaseFurnitureStock(ctx, stockItems)

	s.logger.Info("garniture deleted successfully", "garnitureID", id)

	return nil
}

func (s *GarnitureService) releaseFurnitureStock(ctx context.Context, items []helpers.StockItem) {
	for _, item := range items {
		if err := s.furnitureRepo.IncreaseStock(ctx, item.ID, item.Count); err != nil {
			s.logger.Error(
				"failed to release furniture stock",
				"furnitureID", item.ID,
				"count", item.Count,
				"error", err,
			)
		}
	}
}
