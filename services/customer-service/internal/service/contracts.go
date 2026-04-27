package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model/relations"
)

type FurnitureRepository interface {
	GetList(ctx context.Context) ([]model.Furniture, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Furniture, error)
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (model.Furniture, error)
	Create(ctx context.Context, furniture model.Furniture) error
	Update(ctx context.Context, furniture model.Furniture) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncreaseStock(ctx context.Context, id uuid.UUID, count int) error
	DecreaseStock(ctx context.Context, id uuid.UUID, count int) error
}

type FurnitureModuleRepository interface {
	GetByFurnitureID(ctx context.Context, furnitureID uuid.UUID) ([]relations.FurnitureModule, error)
	CreateMany(ctx context.Context, items []relations.FurnitureModule) error
	DeleteByFurnitureID(ctx context.Context, furnitureID uuid.UUID) error
}

type GarnitureRepository interface {
	GetList(ctx context.Context) ([]model.Garniture, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Garniture, error)
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (model.Garniture, error)
	Create(ctx context.Context, garniture model.Garniture) error
	Update(ctx context.Context, garniture model.Garniture) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncreaseStock(ctx context.Context, id uuid.UUID, count int) error
	DecreaseStock(ctx context.Context, id uuid.UUID, count int) error
}

type GarnitureFurnitureRepository interface {
	GetByGarnitureID(ctx context.Context, garnitureID uuid.UUID) ([]relations.GarnitureFurniture, error)
	CreateMany(ctx context.Context, items []relations.GarnitureFurniture) error
	DeleteByGarnitureID(ctx context.Context, garnitureID uuid.UUID) error
}

type ShipmentRepository interface {
	GetList(ctx context.Context) ([]model.Shipment, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Shipment, error)
	Create(ctx context.Context, shipment model.Shipment) error
	Update(ctx context.Context, shipment model.Shipment) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ShipmentGarnitureRepository interface {
	GetByShipmentID(ctx context.Context, shipmentID uuid.UUID) ([]relations.ShipmentGarniture, error)
	CreateMany(ctx context.Context, items []relations.ShipmentGarniture) error
	DeleteByShipmentID(ctx context.Context, shipmentID uuid.UUID) error
}

type CustomerRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Customer, error)
	GetByLogin(ctx context.Context, login string) (model.Customer, error)
	Create(ctx context.Context, customer model.Customer) error
	ExistsByLogin(ctx context.Context, login string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type SupplierClient interface {
	GetModule(ctx context.Context, moduleID string) (ModuleInfo, error)
	ReserveModules(ctx context.Context, items []ModuleItem) error
	ReleaseModules(ctx context.Context, items []ModuleItem) error
}

type ModuleInfo struct {
	ID         string
	Name       string
	Code       string
	Price      float64
	StockCount int
}

type ModuleItem struct {
	ModuleID string
	Count    int
}
