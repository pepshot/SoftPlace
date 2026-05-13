package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model/relations"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/repository"
)

type TransactionManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context, tx repository.DBExecutor) error) error
}

type CustomerRepository interface {
	Create(ctx context.Context, customer model.Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (model.Customer, error)
	GetByLogin(ctx context.Context, login string) (model.Customer, error)
	ExistsByLogin(ctx context.Context, login string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type FurnitureRepository interface {
	GetList(ctx context.Context) ([]model.Furniture, error)

	GetByID(ctx context.Context, id uuid.UUID) (model.Furniture, error)
	GetByIDTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) (model.Furniture, error)

	Create(ctx context.Context, furniture model.Furniture) error
	CreateTx(ctx context.Context, db repository.DBExecutor, furniture model.Furniture) error

	Update(ctx context.Context, furniture model.Furniture) error
	UpdateTx(ctx context.Context, db repository.DBExecutor, furniture model.Furniture) error

	Delete(ctx context.Context, id uuid.UUID) error
	DeleteTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) error

	IncreaseStock(ctx context.Context, id uuid.UUID, count int) error
	IncreaseStockTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID, count int) error

	DecreaseStock(ctx context.Context, id uuid.UUID, count int) error
	DecreaseStockTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID, count int) error
}

type FurnitureModuleRepository interface {
	CreateMany(ctx context.Context, items []relations.FurnitureModule) error
	CreateManyTx(ctx context.Context, db repository.DBExecutor, items []relations.FurnitureModule) error

	GetByFurnitureID(ctx context.Context, furnitureID uuid.UUID) ([]relations.FurnitureModule, error)
	GetByFurnitureIDTx(ctx context.Context, db repository.DBExecutor, furnitureID uuid.UUID) ([]relations.FurnitureModule, error)

	DeleteByFurnitureID(ctx context.Context, furnitureID uuid.UUID) error
	DeleteByFurnitureIDTx(ctx context.Context, db repository.DBExecutor, furnitureID uuid.UUID) error
}

type GarnitureRepository interface {
	GetList(ctx context.Context) ([]model.Garniture, error)

	GetByID(ctx context.Context, id uuid.UUID) (model.Garniture, error)
	GetByIDTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) (model.Garniture, error)

	Create(ctx context.Context, garniture model.Garniture) error
	CreateTx(ctx context.Context, db repository.DBExecutor, garniture model.Garniture) error

	Update(ctx context.Context, garniture model.Garniture) error
	UpdateTx(ctx context.Context, db repository.DBExecutor, garniture model.Garniture) error

	Delete(ctx context.Context, id uuid.UUID) error
	DeleteTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) error

	IncreaseStock(ctx context.Context, id uuid.UUID, count int) error
	IncreaseStockTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID, count int) error

	DecreaseStock(ctx context.Context, id uuid.UUID, count int) error
	DecreaseStockTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID, count int) error
}

type GarnitureFurnitureRepository interface {
	CreateMany(ctx context.Context, items []relations.GarnitureFurniture) error
	CreateManyTx(ctx context.Context, db repository.DBExecutor, items []relations.GarnitureFurniture) error

	GetByGarnitureID(ctx context.Context, garnitureID uuid.UUID) ([]relations.GarnitureFurniture, error)
	GetByGarnitureIDTx(ctx context.Context, db repository.DBExecutor, garnitureID uuid.UUID) ([]relations.GarnitureFurniture, error)

	DeleteByGarnitureID(ctx context.Context, garnitureID uuid.UUID) error
	DeleteByGarnitureIDTx(ctx context.Context, db repository.DBExecutor, garnitureID uuid.UUID) error
}

type ShipmentRepository interface {
	GetList(ctx context.Context) ([]model.Shipment, error)

	GetByID(ctx context.Context, id uuid.UUID) (model.Shipment, error)
	GetByIDTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) (model.Shipment, error)

	Create(ctx context.Context, shipment model.Shipment) error
	CreateTx(ctx context.Context, db repository.DBExecutor, shipment model.Shipment) error

	Update(ctx context.Context, shipment model.Shipment) error
	UpdateTx(ctx context.Context, db repository.DBExecutor, shipment model.Shipment) error

	Delete(ctx context.Context, id uuid.UUID) error
	DeleteTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) error
}

type ShipmentGarnitureRepository interface {
	CreateMany(ctx context.Context, items []relations.ShipmentGarniture) error
	CreateManyTx(ctx context.Context, db repository.DBExecutor, items []relations.ShipmentGarniture) error

	GetByShipmentID(ctx context.Context, shipmentID uuid.UUID) ([]relations.ShipmentGarniture, error)
	GetByShipmentIDTx(ctx context.Context, db repository.DBExecutor, shipmentID uuid.UUID) ([]relations.ShipmentGarniture, error)

	DeleteByShipmentID(ctx context.Context, shipmentID uuid.UUID) error
	DeleteByShipmentIDTx(ctx context.Context, db repository.DBExecutor, shipmentID uuid.UUID) error
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
