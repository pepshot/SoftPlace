package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model/relations"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/repository"
	"github.com/pepshot/SoftPlace/shared/helpers"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type GarnitureService struct {
	garnitureRepo          GarnitureRepository
	garnitureFurnitureRepo GarnitureFurnitureRepository
	furnitureRepo          FurnitureRepository
	txManager              TransactionManager
	logger                 *logger.Logger
}

func NewGarnitureService(garnitureRepo GarnitureRepository, garnitureFurnitureRepo GarnitureFurnitureRepository,
	furnitureRepo FurnitureRepository, txManager TransactionManager, logger *logger.Logger) *GarnitureService {
	return &GarnitureService{
		garnitureRepo:          garnitureRepo,
		garnitureFurnitureRepo: garnitureFurnitureRepo,
		furnitureRepo:          furnitureRepo,
		txManager:              txManager,
		logger:                 logger,
	}
}

func (s *GarnitureService) GetList(ctx context.Context) ([]dto.GarnitureResponse, error) {
	garnitureList, err := s.garnitureRepo.GetList(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.ToGarnitureResponseList(garnitureList), nil
}

func (s *GarnitureService) GetByID(ctx context.Context, id uuid.UUID) (dto.GarnitureResponse, error) {
	item, err := s.garnitureRepo.GetByID(ctx, id)
	if err != nil {
		return dto.GarnitureResponse{}, ErrNotFound
	}

	return mapper.ToGarnitureResponse(item), nil
}

func (s *GarnitureService) Create(ctx context.Context, req dto.GarnitureRequest) (uuid.UUID, error) {
	if len(req.Furniture) == 0 {
		return uuid.Nil, ErrEmptyComposition
	}

	garnitureID := uuid.New()

	err := s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		relationItems := make([]relations.GarnitureFurniture, 0, len(req.Furniture))
		priceItems := make([]helpers.PriceItem, 0, len(req.Furniture))

		for _, item := range req.Furniture {
			if item.Count <= 0 {
				return ErrInvalidCount
			}

			furnitureID, err := uuid.Parse(item.FurnitureID)
			if err != nil {
				return ErrInvalidID
			}

			furniture, err := s.furnitureRepo.GetByIDTx(ctx, tx, furnitureID)
			if err != nil {
				return ErrNotFound
			}

			if furniture.StockCount < item.Count {
				return ErrNotEnoughStock
			}

			if err := s.furnitureRepo.DecreaseStockTx(ctx, tx, furnitureID, item.Count); err != nil {
				return ErrNotEnoughStock
			}

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

		garniture := model.Garniture{
			ID:         garnitureID,
			Name:       req.Name,
			Code:       req.Code,
			Price:      totalPrice,
			StockCount: 1,
		}

		if err := s.garnitureRepo.CreateTx(ctx, tx, garniture); err != nil {
			return err
		}

		if err := s.garnitureFurnitureRepo.CreateManyTx(ctx, tx, relationItems); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("failed to create garniture", "error", err)
		return uuid.Nil, err
	}

	s.logger.Info("garniture created successfully", "garnitureID", garnitureID)

	return garnitureID, nil
}

func (s *GarnitureService) Update(ctx context.Context, id uuid.UUID, req dto.GarnitureRequest) error {
	if len(req.Furniture) == 0 {
		return ErrEmptyComposition
	}

	err := s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		oldGarniture, err := s.garnitureRepo.GetByIDTx(ctx, tx, id)
		if err != nil {
			return ErrNotFound
		}

		oldRelations, err := s.garnitureFurnitureRepo.GetByGarnitureIDTx(ctx, tx, id)
		if err != nil {
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

			furniture, err := s.furnitureRepo.GetByIDTx(ctx, tx, furnitureID)
			if err != nil {
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
			furniture, err := s.furnitureRepo.GetByIDTx(ctx, tx, item.ID)
			if err != nil {
				return ErrNotFound
			}

			if furniture.StockCount < item.Count {
				return ErrNotEnoughStock
			}

			if err := s.furnitureRepo.DecreaseStockTx(ctx, tx, item.ID, item.Count); err != nil {
				return ErrNotEnoughStock
			}
		}

		for _, item := range toRelease {
			if err := s.furnitureRepo.IncreaseStockTx(ctx, tx, item.ID, item.Count); err != nil {
				return err
			}
		}

		totalPrice := helpers.CalculateTotalPrice(priceItems)

		updatedGarniture := model.Garniture{
			ID:         id,
			Name:       req.Name,
			Code:       req.Code,
			Price:      totalPrice,
			StockCount: oldGarniture.StockCount,
		}

		if err := s.garnitureRepo.UpdateTx(ctx, tx, updatedGarniture); err != nil {
			return err
		}

		if err := s.garnitureFurnitureRepo.DeleteByGarnitureIDTx(ctx, tx, id); err != nil {
			return err
		}

		if err := s.garnitureFurnitureRepo.CreateManyTx(ctx, tx, newRelations); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("failed to update garniture", "garnitureID", id, "error", err)
		return err
	}

	s.logger.Info("garniture updated successfully", "garnitureID", id)

	return nil
}

func (s *GarnitureService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		relationsItems, err := s.garnitureFurnitureRepo.GetByGarnitureIDTx(ctx, tx, id)
		if err != nil {
			return err
		}

		for _, item := range relationsItems {
			if err := s.furnitureRepo.IncreaseStockTx(ctx, tx, item.FurnitureID, item.Count); err != nil {
				return err
			}
		}

		if err := s.garnitureFurnitureRepo.DeleteByGarnitureIDTx(ctx, tx, id); err != nil {
			return err
		}

		if err := s.garnitureRepo.DeleteTx(ctx, tx, id); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("failed to delete garniture", "garnitureID", id, "error", err)
		return err
	}

	s.logger.Info("garniture deleted successfully", "garnitureID", id)

	return nil
}
