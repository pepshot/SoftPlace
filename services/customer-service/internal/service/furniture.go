package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/repository"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/helpers"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model/relations"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type FurnitureService struct {
	furnitureRepo       FurnitureRepository
	furnitureModuleRepo FurnitureModuleRepository
	supplierClient      SupplierClient
	txManager           TransactionManager
	logger              *logger.Logger
}

func NewFurnitureService(furnitureRepo FurnitureRepository, furnitureModuleRepo FurnitureModuleRepository,
	supplierClient SupplierClient, txManager TransactionManager, logger *logger.Logger) *FurnitureService {
	return &FurnitureService{
		furnitureRepo:       furnitureRepo,
		furnitureModuleRepo: furnitureModuleRepo,
		supplierClient:      supplierClient,
		txManager:           txManager,
		logger:              logger,
	}
}

func (s *FurnitureService) GetList(ctx context.Context) ([]dto.FurnitureResponse, error) {
	s.logger.Debug("getting furniture list")

	furnitureList, err := s.furnitureRepo.GetList(ctx)
	if err != nil {
		s.logger.Error("failed to get furniture list", "error", err)
		return nil, err
	}

	return mapper.ToFurnitureResponseList(furnitureList), nil
}

func (s *FurnitureService) GetByID(ctx context.Context, id uuid.UUID) (dto.FurnitureResponse, error) {
	s.logger.Debug("getting furniture by id", "furnitureID", id)

	item, err := s.furnitureRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get furniture by id", "furnitureID", id, "error", err)
		return dto.FurnitureResponse{}, err
	}

	return mapper.ToFurnitureResponse(item), nil
}

func (s *FurnitureService) Create(ctx context.Context, req dto.FurnitureRequest) (uuid.UUID, error) {
	s.logger.Info("starting furniture creation", "name", req.Name, "code", req.Code)

	if len(req.Modules) == 0 {
		return uuid.Nil, ErrEmptyComposition
	}

	furnitureID := uuid.New()

	moduleItems := make([]ModuleItem, 0, len(req.Modules))
	relationItems := make([]relations.FurnitureModule, 0, len(req.Modules))
	priceItems := make([]helpers.PriceItem, 0, len(req.Modules))

	for _, item := range req.Modules {
		if item.Count <= 0 {
			return uuid.Nil, ErrInvalidCount
		}

		moduleID, err := uuid.Parse(item.ModuleID)
		if err != nil {
			return uuid.Nil, ErrInvalidID
		}

		module, err := s.supplierClient.GetModule(ctx, item.ModuleID)
		if err != nil {
			s.logger.Error(
				"failed to get module from supplier-service",
				"moduleID", item.ModuleID,
				"error", err,
			)

			return uuid.Nil, err
		}

		if module.StockCount < item.Count {
			s.logger.Warn(
				"not enough module stock",
				"moduleID", item.ModuleID,
				"need", item.Count,
				"available", module.StockCount,
			)

			return uuid.Nil, ErrNotEnoughStock
		}

		moduleItems = append(moduleItems, ModuleItem{
			ModuleID: item.ModuleID,
			Count:    item.Count,
		})

		relationItems = append(relationItems, relations.FurnitureModule{
			FurnitureID: furnitureID,
			ModuleID:    moduleID,
			Count:       item.Count,
		})

		priceItems = append(priceItems, helpers.PriceItem{
			Price: module.Price,
			Count: item.Count,
		})
	}

	if err := s.supplierClient.ReserveModules(ctx, moduleItems); err != nil {
		s.logger.Error("failed to reserve modules", "error", err)
		return uuid.Nil, err
	}

	totalPrice := helpers.CalculateTotalPrice(priceItems)

	err := s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		furniture := model.Furniture{
			ID:         furnitureID,
			Name:       req.Name,
			Code:       req.Code,
			Price:      totalPrice,
			StockCount: 1,
		}

		if err := s.furnitureRepo.CreateTx(ctx, tx, furniture); err != nil {
			return err
		}

		if err := s.furnitureModuleRepo.CreateManyTx(ctx, tx, relationItems); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		_ = s.supplierClient.ReleaseModules(ctx, moduleItems)
		s.logger.Error("failed to create furniture transaction", "error", err)
		return uuid.Nil, err
	}

	s.logger.Info(
		"furniture created successfully",
		"furnitureID", furnitureID,
		"price", totalPrice,
	)

	return furnitureID, nil
}

func (s *FurnitureService) Update(ctx context.Context, id uuid.UUID, req dto.FurnitureRequest) error {
	s.logger.Info("starting furniture update", "furnitureID", id)

	if len(req.Modules) == 0 {
		return ErrEmptyComposition
	}

	oldRelations, err := s.furnitureModuleRepo.GetByFurnitureID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get old furniture-module relations", "furnitureID", id, "error", err)
		return err
	}

	oldStockItems := make([]helpers.StockItem, 0, len(oldRelations))

	for _, item := range oldRelations {
		oldStockItems = append(oldStockItems, helpers.StockItem{
			ID:    item.ModuleID,
			Count: item.Count,
		})
	}

	newStockItems := make([]helpers.StockItem, 0, len(req.Modules))
	newRelations := make([]relations.FurnitureModule, 0, len(req.Modules))
	priceItems := make([]helpers.PriceItem, 0, len(req.Modules))

	for _, item := range req.Modules {
		if item.Count <= 0 {
			return ErrInvalidCount
		}

		moduleID, err := uuid.Parse(item.ModuleID)
		if err != nil {
			return ErrInvalidID
		}

		moduleInfo, err := s.supplierClient.GetModule(ctx, item.ModuleID)
		if err != nil {
			return ErrExternalService
		}

		newStockItems = append(newStockItems, helpers.StockItem{
			ID:    moduleID,
			Count: item.Count,
		})

		newRelations = append(newRelations, relations.FurnitureModule{
			FurnitureID: id,
			ModuleID:    moduleID,
			Count:       item.Count,
		})

		priceItems = append(priceItems, helpers.PriceItem{
			Price: moduleInfo.Price,
			Count: item.Count,
		})
	}

	toReserveStock, toReleaseStock := helpers.CalculateStockDiff(oldStockItems, newStockItems)

	toReserve := make([]ModuleItem, 0, len(toReserveStock))
	for _, item := range toReserveStock {
		toReserve = append(toReserve, ModuleItem{
			ModuleID: item.ID.String(),
			Count:    item.Count,
		})
	}

	toRelease := make([]ModuleItem, 0, len(toReleaseStock))
	for _, item := range toReleaseStock {
		toRelease = append(toRelease, ModuleItem{
			ModuleID: item.ID.String(),
			Count:    item.Count,
		})
	}

	if err := s.supplierClient.ReserveModules(ctx, toReserve); err != nil {
		return ErrExternalService
	}

	totalPrice := helpers.CalculateTotalPrice(priceItems)

	err = s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		oldFurniture, err := s.furnitureRepo.GetByIDTx(ctx, tx, id)
		if err != nil {
			return ErrNotFound
		}

		updatedFurniture := model.Furniture{
			ID:         id,
			Name:       req.Name,
			Code:       req.Code,
			Price:      totalPrice,
			StockCount: oldFurniture.StockCount,
		}

		if err := s.furnitureRepo.UpdateTx(ctx, tx, updatedFurniture); err != nil {
			return err
		}

		if err := s.furnitureModuleRepo.DeleteByFurnitureIDTx(ctx, tx, id); err != nil {
			return err
		}

		if err := s.furnitureModuleRepo.CreateManyTx(ctx, tx, newRelations); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		_ = s.supplierClient.ReleaseModules(ctx, toReserve)
		return err
	}

	_ = s.supplierClient.ReleaseModules(ctx, toRelease)

	s.logger.Info("furniture updated successfully", "furnitureID", id, "price", totalPrice)

	return nil
}

func (s *FurnitureService) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting furniture", "furnitureID", id)

	oldRelations, err := s.furnitureModuleRepo.GetByFurnitureID(ctx, id)
	if err != nil {
		return err
	}

	toRelease := make([]ModuleItem, 0, len(oldRelations))

	for _, item := range oldRelations {
		toRelease = append(toRelease, ModuleItem{
			ModuleID: item.ModuleID.String(),
			Count:    item.Count,
		})
	}

	err = s.txManager.RunInTx(ctx, func(ctx context.Context, tx repository.DBExecutor) error {
		if err := s.furnitureModuleRepo.DeleteByFurnitureIDTx(ctx, tx, id); err != nil {
			return err
		}

		if err := s.furnitureRepo.DeleteTx(ctx, tx, id); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	_ = s.supplierClient.ReleaseModules(ctx, toRelease)

	s.logger.Info("furniture deleted successfully", "furnitureID", id)

	return nil
}
