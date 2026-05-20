package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models/relations"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/repository"
)

type TransactionManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context, tx repository.DBExecutor) error) error
}

type ModuleRepository interface {
	GetList(ctx context.Context) ([]models.Module, error)

	GetByID(ctx context.Context, id uuid.UUID) (models.Module, error)
	GetByIDTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) (models.Module, error)

	GetByCode(ctx context.Context, code string) (models.Module, error)
	GetByCodeTx(ctx context.Context, db repository.DBExecutor, code string) (models.Module, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)

	Create(ctx context.Context, module models.Module) error
	CreateTx(ctx context.Context, db repository.DBExecutor, module models.Module) error

	Update(ctx context.Context, module models.Module) error
	UpdateTx(ctx context.Context, db repository.DBExecutor, module models.Module) error

	Delete(ctx context.Context, id uuid.UUID) error
	DeleteTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) error

	IncreaseStock(ctx context.Context, id uuid.UUID, count int) error
	IncreaseStockTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID, count int) error

	DecreaseStock(ctx context.Context, id uuid.UUID, count int) error
	DecreaseStockTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID, count int) error
}

type ModuleMaterialRepository interface {
	CreateMany(ctx context.Context, items []relations.ModuleMaterial) error
	CreateManyTx(ctx context.Context, db repository.DBExecutor, items []relations.ModuleMaterial) error

	GetByModuleID(ctx context.Context, moduleID uuid.UUID) ([]relations.ModuleMaterial, error)
	GetByModuleIDTx(ctx context.Context, db repository.DBExecutor, moduleID uuid.UUID) ([]relations.ModuleMaterial, error)

	DeleteByModuleID(ctx context.Context, moduleID uuid.UUID) error
	DeleteByModuleIDTx(ctx context.Context, db repository.DBExecutor, moduleID uuid.UUID) error
}

type MaterialRepository interface {
	GetList(ctx context.Context) ([]models.Material, error)

	GetByID(ctx context.Context, id uuid.UUID) (models.Material, error)
	GetByIDTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) (models.Material, error)

	ExistsByCode(ctx context.Context, code string) (bool, error)

	Create(ctx context.Context, item models.Material) error
	CreateTx(ctx context.Context, db repository.DBExecutor, item models.Material) error

	Update(ctx context.Context, item models.Material) error
	UpdateTx(ctx context.Context, db repository.DBExecutor, item models.Material) error

	Delete(ctx context.Context, id uuid.UUID) error
	DeleteTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) error

	IncreaseStock(ctx context.Context, id uuid.UUID, count int) error
	IncreaseStockTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID, count int) error

	DecreaseStock(ctx context.Context, id uuid.UUID, count int) error
	DecreaseStockTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID, count int) error
}

type SupplyRepository interface {
	GetList(ctx context.Context) ([]models.Supply, error)

	GetByID(ctx context.Context, id uuid.UUID) (models.Supply, error)
	GetByIDTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) (models.Supply, error)

	ExistsByCode(ctx context.Context, code string) (bool, error)

	Create(ctx context.Context, item models.Supply) error
	CreateTx(ctx context.Context, db repository.DBExecutor, item models.Supply) error

	Update(ctx context.Context, item models.Supply) error
	UpdateTx(ctx context.Context, db repository.DBExecutor, item models.Supply) error

	Delete(ctx context.Context, id uuid.UUID) error
	DeleteTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) error
}

type SupplyMaterialRepository interface {
	CreateMany(ctx context.Context, items []relations.SupplyMaterial) error
	CreateManyTx(ctx context.Context, db repository.DBExecutor, items []relations.SupplyMaterial) error

	GetBySupplyID(ctx context.Context, supplyID uuid.UUID) ([]relations.SupplyMaterial, error)
	GetBySupplyIDTx(ctx context.Context, db repository.DBExecutor, supplyID uuid.UUID) ([]relations.SupplyMaterial, error)

	DeleteBySupplyID(ctx context.Context, supplyID uuid.UUID) error
	DeleteBySupplyIDTx(ctx context.Context, db repository.DBExecutor, supplyID uuid.UUID) error
}

type SupplierRepository interface {
	Create(ctx context.Context, supplier models.Supplier) error
	GetByID(ctx context.Context, id uuid.UUID) (models.Supplier, error)
	GetByLogin(ctx context.Context, login string) (models.Supplier, error)
	ExistsByLogin(ctx context.Context, login string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
