package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
)

type ModuleUseCase interface {
	GetList(ctx context.Context) ([]dto.ModuleResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (dto.ModuleResponse, error)
	Create(ctx context.Context, req dto.ModuleRequest) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, req dto.ModuleRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	ReserveModules(ctx context.Context, items []ModuleItem) error
	ReleaseModules(ctx context.Context, items []ModuleItem) error
}

type MaterialUseCase interface {
	GetList(ctx context.Context) ([]dto.MaterialResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (dto.MaterialResponse, error)
	Create(ctx context.Context, req dto.MaterialRequest) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, req dto.MaterialRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SupplyUseCase interface {
	GetList(ctx context.Context) ([]dto.SupplyResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (dto.SupplyResponse, error)
	Create(ctx context.Context, supplierID uuid.UUID, req dto.SupplyRequest) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, req dto.SupplyRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SupplierUseCase interface {
	Register(ctx context.Context, req dto.RegisterSupplierRequest) (dto.AuthSupplierResponse, error)
	Login(ctx context.Context, req dto.LoginSupplierRequest) (dto.AuthSupplierResponse, error)
	GetProfile(ctx context.Context, supplierID uuid.UUID) (dto.SupplierResponse, error)
}

type JWTUseCase interface {
	GenerateToken(userID string, role string) (string, error)
}

type ModuleItem struct {
	ModuleID string
	Count    int
}
