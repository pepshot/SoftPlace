package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
)

type FurnitureUseCase interface {
	GetList(ctx context.Context) ([]dto.FurnitureResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (dto.FurnitureResponse, error)
	Create(ctx context.Context, req dto.FurnitureRequest) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, req dto.FurnitureRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type GarnitureUseCase interface {
	GetList(ctx context.Context) ([]dto.GarnitureResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (dto.GarnitureResponse, error)
	Create(ctx context.Context, req dto.GarnitureRequest) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, req dto.GarnitureRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ShipmentUseCase interface {
	GetList(ctx context.Context) ([]dto.ShipmentResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (dto.ShipmentResponse, error)
	Create(ctx context.Context, customerID uuid.UUID, req dto.ShipmentRequest) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, req dto.ShipmentRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CustomerUseCase interface {
	Register(ctx context.Context, req dto.RegisterCustomerRequest) (dto.AuthCustomerResponse, error)
	Login(ctx context.Context, req dto.LoginCustomerRequest) (dto.AuthCustomerResponse, error)
	GetProfile(ctx context.Context, customerID uuid.UUID) (dto.CustomerResponse, error)
}

type JWTUseCase interface {
	GenerateToken(userID string, role string) (string, error)
}
