package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models/relations"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/repository"
	"github.com/pepshot/SoftPlace/shared/helpers"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type SupplyService struct {
	supplyRepo         SupplyRepository
	supplyMaterialRepo SupplyMaterialRepository
	materialRepo       MaterialRepository
	txManager          TransactionManager
	logger             *logger.Logger
}

func NewSupplyService(supplyRepo SupplyRepository, supplyMaterialRepo SupplyMaterialRepository, materialRepo MaterialRepository, txManager TransactionManager, logger *logger.Logger) *SupplyService {
	return &SupplyService{
		supplyRepo: supplyRepo, supplyMaterialRepo: supplyMaterialRepo, materialRepo: materialRepo, txManager: txManager, logger: logger,
	}
}

func (s *SupplyService) GetList(ctx context.Context) ([]dto.SupplyResponse, error) {
	items, err := s.supplyRepo.GetList(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.ToSupplyResponseList(items), nil
}

func (s *SupplyService) GetByID(ctx context.Context, id uuid.UUID) (dto.SupplyResponse, error) {
	item, err := s.supplyRepo.GetByID(ctx, id)
	if err != nil {
		return dto.SupplyResponse{}, ErrNotFound
	}
	return mapper.ToSupplyResponse(item), nil
}

func (s *SupplyService) Create(ctx context.Context, supplierID uuid.UUID, req dto.SupplyRequest) (uuid.UUID, error) {
	if len(req.Materials) == 0 {
		return uuid.Nil, ErrEmptyComposition
	}

	supplyDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return uuid.Nil, ErrInvalidDate
	}

	exists, err := s.supplyRepo.ExistsByCode(ctx, req.Code)
	if err != nil {
		return uuid.Nil, err
	}
	if exists {
		return uuid.Nil, ErrCodeAlreadyUsed
	}

	supplyID := uuid.New()
	err = s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		relationsItems := make([]relations.SupplyMaterial, 0, len(req.Materials))
		priceItems := make([]helpers.PriceItem, 0, len(req.Materials))

		for _, item := range req.Materials {
			if item.Count <= 0 {
				return ErrInvalidCount
			}

			materialID, err := uuid.Parse(item.MaterialID)
			if err != nil {
				return ErrInvalidID
			}

			material, err := s.materialRepo.GetByIDTx(ctx, tx, materialID)
			if err != nil {
				return ErrNotFound
			}

			if err := s.materialRepo.IncreaseStockTx(ctx, tx, materialID, item.Count); err != nil {
				return err
			}

			relationsItems = append(relationsItems, relations.SupplyMaterial{SupplyID: supplyID, MaterialID: materialID, Count: item.Count})
			priceItems = append(priceItems, helpers.PriceItem{Price: material.Price, Count: item.Count})
		}

		totalPrice := helpers.CalculateTotalPrice(priceItems)
		supply := models.Supply{ID: supplyID, SupplierID: supplierID, Code: req.Code, Date: supplyDate, Price: totalPrice}

		if err := s.supplyRepo.CreateTx(ctx, tx, supply); err != nil {
			return err
		}
		if err := s.supplyMaterialRepo.CreateManyTx(ctx, tx, relationsItems); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return supplyID, nil
}

func (s *SupplyService) Update(ctx context.Context, id uuid.UUID, req dto.SupplyRequest) error {
	if len(req.Materials) == 0 {
		return ErrEmptyComposition
	}

	supplyDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return ErrInvalidDate
	}

	err = s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		oldSupply, err := s.supplyRepo.GetByIDTx(ctx, tx, id)
		if err != nil {
			return ErrNotFound
		}

		oldRelations, err := s.supplyMaterialRepo.GetBySupplyIDTx(ctx, tx, id)
		if err != nil {
			return err
		}

		oldStockItems := make([]helpers.StockItem, 0, len(oldRelations))
		for _, item := range oldRelations {
			oldStockItems = append(oldStockItems, helpers.StockItem{ID: item.MaterialID, Count: item.Count})
		}

		newStockItems := make([]helpers.StockItem, 0, len(req.Materials))
		newRelations := make([]relations.SupplyMaterial, 0, len(req.Materials))
		priceItems := make([]helpers.PriceItem, 0, len(req.Materials))

		for _, item := range req.Materials {
			if item.Count <= 0 {
				return ErrInvalidCount
			}
			materialID, err := uuid.Parse(item.MaterialID)
			if err != nil {
				return ErrInvalidID
			}
			material, err := s.materialRepo.GetByIDTx(ctx, tx, materialID)
			if err != nil {
				return ErrNotFound
			}

			newStockItems = append(newStockItems, helpers.StockItem{ID: materialID, Count: item.Count})
			newRelations = append(newRelations, relations.SupplyMaterial{SupplyID: id, MaterialID: materialID, Count: item.Count})
			priceItems = append(priceItems, helpers.PriceItem{Price: material.Price, Count: item.Count})
		}

		toReserve, toRelease := helpers.CalculateStockDiff(oldStockItems, newStockItems)
		for _, item := range toReserve {
			if err := s.materialRepo.IncreaseStockTx(ctx, tx, item.ID, item.Count); err != nil {
				return err
			}
		}

		for _, item := range toRelease {
			material, err := s.materialRepo.GetByIDTx(ctx, tx, item.ID)
			if err != nil {
				return ErrNotFound
			}
			if material.StockCount < item.Count {
				return ErrNotEnoughStock
			}
			if err := s.materialRepo.DecreaseStockTx(ctx, tx, item.ID, item.Count); err != nil {
				return ErrNotEnoughStock
			}
		}

		totalPrice := helpers.CalculateTotalPrice(priceItems)
		updated := models.Supply{ID: id, SupplierID: oldSupply.SupplierID, Code: req.Code, Date: supplyDate, Price: totalPrice}
		if err := s.supplyRepo.UpdateTx(ctx, tx, updated); err != nil {
			return err
		}

		if err := s.supplyMaterialRepo.DeleteBySupplyIDTx(ctx, tx, id); err != nil {
			return err
		}
		if err := s.supplyMaterialRepo.CreateManyTx(ctx, tx, newRelations); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *SupplyService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		relationsItems, err := s.supplyMaterialRepo.GetBySupplyIDTx(ctx, tx, id)
		if err != nil {
			return err
		}

		for _, item := range relationsItems {
			material, err := s.materialRepo.GetByIDTx(ctx, tx, item.MaterialID)
			if err != nil {
				return ErrNotFound
			}
			if material.StockCount < item.Count {
				return ErrNotEnoughStock
			}
			if err := s.materialRepo.DecreaseStockTx(ctx, tx, item.MaterialID, item.Count); err != nil {
				return ErrNotEnoughStock
			}
		}

		if err := s.supplyMaterialRepo.DeleteBySupplyIDTx(ctx, tx, id); err != nil {
			return err
		}
		if err := s.supplyRepo.DeleteTx(ctx, tx, id); err != nil {
			return err
		}
		return nil
	})
}
