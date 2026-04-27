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

type FurnitureService struct {
	furnitureRepo       FurnitureRepository
	furnitureModuleRepo FurnitureModuleRepository
	supplierClient      SupplierClient
	logger              *logger.Logger
}

func NewFurnitureService(furnitureRepo FurnitureRepository, furnitureModuleRepo FurnitureModuleRepository,
	supplierClient SupplierClient, logger *logger.Logger) *FurnitureService {
	return &FurnitureService{
		furnitureRepo:       furnitureRepo,
		furnitureModuleRepo: furnitureModuleRepo,
		supplierClient:      supplierClient,
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

	totalPrice := helpers.CalculateTotalPrice(priceItems)

	if err := s.supplierClient.ReserveModules(ctx, moduleItems); err != nil {
		s.logger.Error("failed to reserve modules", "error", err)
		return uuid.Nil, err
	}

	furniture := model.Furniture{
		ID:         furnitureID,
		Name:       req.Name,
		Code:       req.Code,
		Price:      totalPrice,
		StockCount: 1,
	}

	if err := s.furnitureRepo.Create(ctx, furniture); err != nil {
		s.logger.Error("failed to create furniture", "error", err)

		_ = s.supplierClient.ReleaseModules(ctx, moduleItems)

		return uuid.Nil, err
	}

	if err := s.furnitureModuleRepo.CreateMany(ctx, relationItems); err != nil {
		s.logger.Error("failed to create furniture-module relations", "error", err)

		_ = s.supplierClient.ReleaseModules(ctx, moduleItems)
		_ = s.furnitureRepo.Delete(ctx, furnitureID)

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

	oldFurniture, err := s.furnitureRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(
			"failed to get old furniture",
			"furnitureID", id,
			"error", err,
		)

		return err
	}

	oldRelations, err := s.furnitureModuleRepo.GetByFurnitureID(ctx, id)
	if err != nil {
		s.logger.Error(
			"failed to get old furniture-module relations",
			"furnitureID", id,
			"error", err,
		)

		return err
	}

	oldModuleItems := make([]ModuleItem, 0, len(oldRelations))

	for _, item := range oldRelations {
		oldModuleItems = append(oldModuleItems, ModuleItem{
			ModuleID: item.ModuleID.String(),
			Count:    item.Count,
		})
	}

	newModuleItems := make([]ModuleItem, 0, len(req.Modules))
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

		module, err := s.supplierClient.GetModule(ctx, item.ModuleID)
		if err != nil {
			s.logger.Error(
				"failed to get module from supplier-service",
				"moduleID", item.ModuleID,
				"error", err,
			)

			return err
		}

		if module.StockCount < item.Count {
			s.logger.Warn(
				"not enough module stock",
				"moduleID", item.ModuleID,
				"need", item.Count,
				"available", module.StockCount,
			)

			return ErrNotEnoughStock
		}

		newModuleItems = append(newModuleItems, ModuleItem{
			ModuleID: item.ModuleID,
			Count:    item.Count,
		})

		newRelations = append(newRelations, relations.FurnitureModule{
			FurnitureID: id,
			ModuleID:    moduleID,
			Count:       item.Count,
		})

		priceItems = append(priceItems, helpers.PriceItem{
			Price: module.Price,
			Count: item.Count,
		})
	}

	totalPrice := helpers.CalculateTotalPrice(priceItems)

	if err := s.supplierClient.ReserveModules(ctx, newModuleItems); err != nil {
		s.logger.Error("failed to reserve new modules", "furnitureID", id, "error", err)
		return err
	}

	updatedFurniture := model.Furniture{
		ID:         id,
		Name:       req.Name,
		Code:       req.Code,
		Price:      totalPrice,
		StockCount: oldFurniture.StockCount,
	}

	if err := s.furnitureRepo.Update(ctx, updatedFurniture); err != nil {
		s.logger.Error("failed to update furniture", "furnitureID", id, "error", err)

		_ = s.supplierClient.ReleaseModules(ctx, newModuleItems)

		return err
	}

	if err := s.furnitureModuleRepo.DeleteByFurnitureID(ctx, id); err != nil {
		s.logger.Error(
			"failed to delete old furniture-module relations",
			"furnitureID", id,
			"error", err,
		)

		_ = s.supplierClient.ReleaseModules(ctx, newModuleItems)

		return err
	}

	if err := s.furnitureModuleRepo.CreateMany(ctx, newRelations); err != nil {
		s.logger.Error(
			"failed to create new furniture-module relations",
			"furnitureID", id,
			"error", err,
		)

		_ = s.supplierClient.ReleaseModules(ctx, newModuleItems)

		return err
	}

	_ = s.supplierClient.ReleaseModules(ctx, oldModuleItems)

	s.logger.Info(
		"furniture updated successfully",
		"furnitureID", id,
		"price", totalPrice,
	)

	return nil
}

func (s *FurnitureService) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting furniture", "furnitureID", id)

	relationsItems, err := s.furnitureModuleRepo.GetByFurnitureID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get furniture-module relations", "furnitureID", id, "error", err)
		return err
	}

	moduleItems := make([]ModuleItem, 0, len(relationsItems))

	for _, item := range relationsItems {
		moduleItems = append(moduleItems, ModuleItem{
			ModuleID: item.ModuleID.String(),
			Count:    item.Count,
		})
	}

	if err := s.furnitureModuleRepo.DeleteByFurnitureID(ctx, id); err != nil {
		return err
	}

	if err := s.furnitureRepo.Delete(ctx, id); err != nil {
		return err
	}

	_ = s.supplierClient.ReleaseModules(ctx, moduleItems)

	s.logger.Info("furniture deleted successfully", "furnitureID", id)

	return nil
}
