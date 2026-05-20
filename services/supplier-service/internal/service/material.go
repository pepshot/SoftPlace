package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/mapper"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type MaterialService struct {
	materialRepo MaterialRepository
	logger       *logger.Logger
}

func NewMaterialService(materialRepo MaterialRepository, logger *logger.Logger) *MaterialService {
	return &MaterialService{materialRepo: materialRepo, logger: logger}
}

func (s *MaterialService) GetList(ctx context.Context) ([]dto.MaterialResponse, error) {
	items, err := s.materialRepo.GetList(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.ToMaterialResponseList(items), nil
}

func (s *MaterialService) GetByID(ctx context.Context, id uuid.UUID) (dto.MaterialResponse, error) {
	item, err := s.materialRepo.GetByID(ctx, id)
	if err != nil {
		return dto.MaterialResponse{}, ErrNotFound
	}
	return mapper.ToMaterialResponse(item), nil
}

func (s *MaterialService) Create(ctx context.Context, req dto.MaterialRequest) (uuid.UUID, error) {
	if req.Price < 0 {
		return uuid.Nil, ErrInvalidPrice
	}
	if req.StockCount < 0 {
		return uuid.Nil, ErrInvalidCount
	}

	exists, err := s.materialRepo.ExistsByCode(ctx, req.Code)
	if err != nil {
		return uuid.Nil, err
	}
	if exists {
		return uuid.Nil, ErrCodeAlreadyUsed
	}

	item := models.Material{
		ID:         uuid.New(),
		Name:       req.Name,
		Code:       req.Code,
		Price:      req.Price,
		StockCount: req.StockCount,
	}

	if err := s.materialRepo.Create(ctx, item); err != nil {
		return uuid.Nil, err
	}
	return item.ID, nil
}

func (s *MaterialService) Update(ctx context.Context, id uuid.UUID, req dto.MaterialRequest) error {
	if req.Price < 0 {
		return ErrInvalidPrice
	}
	if req.StockCount < 0 {
		return ErrInvalidCount
	}

	oldItem, err := s.materialRepo.GetByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}

	if oldItem.Code != req.Code {
		exists, err := s.materialRepo.ExistsByCode(ctx, req.Code)
		if err != nil {
			return err
		}
		if exists {
			return ErrCodeAlreadyUsed
		}
	}

	return s.materialRepo.Update(ctx, models.Material{
		ID:         id,
		Name:       req.Name,
		Code:       req.Code,
		Price:      req.Price,
		StockCount: req.StockCount,
	})
}

func (s *MaterialService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.materialRepo.Delete(ctx, id); err != nil {
		return ErrNotFound
	}
	return nil
}
