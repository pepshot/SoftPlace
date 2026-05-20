package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models/relations"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/repository"
	"github.com/pepshot/SoftPlace/shared/helpers"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type ModuleService struct {
	moduleRepo         ModuleRepository
	moduleMaterialRepo ModuleMaterialRepository
	materialRepo       MaterialRepository
	txManager          TransactionManager
	logger             *logger.Logger
}

func NewModuleService(moduleRepo ModuleRepository, moduleMaterialRepo ModuleMaterialRepository,
	materialRepo MaterialRepository, txManager TransactionManager, logger *logger.Logger) *ModuleService {
	return &ModuleService{
		moduleRepo:         moduleRepo,
		moduleMaterialRepo: moduleMaterialRepo,
		materialRepo:       materialRepo,
		txManager:          txManager,
		logger:             logger,
	}
}

func (s *ModuleService) GetList(ctx context.Context) ([]dto.ModuleResponse, error) {
	modules, err := s.moduleRepo.GetList(ctx)
	if err != nil {
		s.logger.Error("failed to get module list", "error", err)
		return nil, err
	}

	return mapper.ToModuleResponseList(modules), nil
}

func (s *ModuleService) GetByID(ctx context.Context, id uuid.UUID) (dto.ModuleResponse, error) {
	module, err := s.moduleRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get module by id", "moduleID", id, "error", err)
		return dto.ModuleResponse{}, ErrNotFound
	}

	return mapper.ToModuleResponse(module), nil
}

func (s *ModuleService) Create(ctx context.Context, req dto.ModuleRequest) (uuid.UUID, error) {
	s.logger.Info("starting module creation", "name", req.Name, "code", req.Code)

	if len(req.Materials) == 0 {
		return uuid.Nil, ErrEmptyComposition
	}

	exists, err := s.moduleRepo.ExistsByCode(ctx, req.Code)
	if err != nil {
		s.logger.Error("failed to check module code", "code", req.Code, "error", err)
		return uuid.Nil, err
	}

	if exists {
		return uuid.Nil, ErrCodeAlreadyUsed
	}

	moduleID := uuid.New()

	err = s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		relationItems := make([]relations.ModuleMaterial, 0, len(req.Materials))
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

			if material.StockCount < item.Count {
				return ErrNotEnoughStock
			}

			if err := s.materialRepo.DecreaseStockTx(ctx, tx, materialID, item.Count); err != nil {
				return ErrNotEnoughStock
			}

			relationItems = append(relationItems, relations.ModuleMaterial{
				ModuleID:   moduleID,
				MaterialID: materialID,
				Count:      item.Count,
			})

			priceItems = append(priceItems, helpers.PriceItem{Price: material.Price, Count: item.Count})
		}

		module := models.Module{
			ID:         moduleID,
			Name:       req.Name,
			Code:       req.Code,
			Price:      helpers.CalculateTotalPrice(priceItems),
			StockCount: 1,
		}

		if err := s.moduleRepo.CreateTx(ctx, tx, module); err != nil {
			return err
		}

		if err := s.moduleMaterialRepo.CreateManyTx(ctx, tx, relationItems); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		s.logger.Error("failed to create module", "code", req.Code, "error", err)
		return uuid.Nil, err
	}

	s.logger.Info("module created successfully", "moduleID", moduleID)
	return moduleID, nil
}

func (s *ModuleService) Update(ctx context.Context, id uuid.UUID, req dto.ModuleRequest) error {
	s.logger.Info("starting module update", "moduleID", id)

	if len(req.Materials) == 0 {
		return ErrEmptyComposition
	}

	err := s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		oldModule, err := s.moduleRepo.GetByIDTx(ctx, tx, id)
		if err != nil {
			return ErrNotFound
		}

		if req.Code != oldModule.Code {
			exists, err := s.moduleRepo.ExistsByCode(ctx, req.Code)
			if err != nil {
				return err
			}
			if exists {
				return ErrCodeAlreadyUsed
			}
		}

		oldRelations, err := s.moduleMaterialRepo.GetByModuleIDTx(ctx, tx, id)
		if err != nil {
			return err
		}

		oldStockItems := make([]helpers.StockItem, 0, len(oldRelations))
		for _, item := range oldRelations {
			oldStockItems = append(oldStockItems, helpers.StockItem{ID: item.MaterialID, Count: item.Count})
		}

		newStockItems := make([]helpers.StockItem, 0, len(req.Materials))
		newRelations := make([]relations.ModuleMaterial, 0, len(req.Materials))
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
			newRelations = append(newRelations, relations.ModuleMaterial{ModuleID: id, MaterialID: materialID, Count: item.Count})
			priceItems = append(priceItems, helpers.PriceItem{Price: material.Price, Count: item.Count})
		}

		toReserve, toRelease := helpers.CalculateStockDiff(oldStockItems, newStockItems)

		for _, item := range toReserve {
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

		for _, item := range toRelease {
			if err := s.materialRepo.IncreaseStockTx(ctx, tx, item.ID, item.Count); err != nil {
				return err
			}
		}

		updatedModule := models.Module{
			ID:         id,
			Name:       req.Name,
			Code:       req.Code,
			Price:      helpers.CalculateTotalPrice(priceItems),
			StockCount: oldModule.StockCount,
		}

		if err := s.moduleRepo.UpdateTx(ctx, tx, updatedModule); err != nil {
			return err
		}

		if err := s.moduleMaterialRepo.DeleteByModuleIDTx(ctx, tx, id); err != nil {
			return err
		}

		if err := s.moduleMaterialRepo.CreateManyTx(ctx, tx, newRelations); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		s.logger.Error("failed to update module", "moduleID", id, "error", err)
		return err
	}

	s.logger.Info("module updated successfully", "moduleID", id)
	return nil
}

func (s *ModuleService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		relationsItems, err := s.moduleMaterialRepo.GetByModuleIDTx(ctx, tx, id)
		if err != nil {
			return err
		}

		for _, item := range relationsItems {
			if err := s.materialRepo.IncreaseStockTx(ctx, tx, item.MaterialID, item.Count); err != nil {
				return err
			}
		}

		if err := s.moduleMaterialRepo.DeleteByModuleIDTx(ctx, tx, id); err != nil {
			return err
		}

		if err := s.moduleRepo.DeleteTx(ctx, tx, id); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		s.logger.Error("failed to delete module", "moduleID", id, "error", err)
		return err
	}

	s.logger.Info("module deleted successfully", "moduleID", id)

	return nil
}

func (s *ModuleService) ReserveModules(ctx context.Context, items []ModuleItem) error {
	if len(items) == 0 {
		return ErrInvalidComposition
	}

	return s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		for _, item := range items {
			if item.Count <= 0 {
				return ErrInvalidCount
			}

			moduleID, err := uuid.Parse(item.ModuleID)
			if err != nil {
				return ErrInvalidID
			}

			module, err := s.moduleRepo.GetByIDTx(ctx, tx, moduleID)
			if err != nil {
				return ErrNotFound
			}

			if module.StockCount < item.Count {
				return ErrNotEnoughStock
			}

			if err := s.moduleRepo.DecreaseStockTx(ctx, tx, moduleID, item.Count); err != nil {
				return ErrNotEnoughStock
			}
		}

		return nil
	})
}

func (s *ModuleService) ReleaseModules(ctx context.Context, items []ModuleItem) error {
	if len(items) == 0 {
		return ErrInvalidComposition
	}

	return s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		for _, item := range items {
			if item.Count <= 0 {
				return ErrInvalidCount
			}

			moduleID, err := uuid.Parse(item.ModuleID)
			if err != nil {
				return ErrInvalidID
			}

			if _, err := s.moduleRepo.GetByIDTx(ctx, tx, moduleID); err != nil {
				return ErrNotFound
			}

			if err := s.moduleRepo.IncreaseStockTx(ctx, tx, moduleID, item.Count); err != nil {
				return err
			}
		}

		return nil
	})
}
