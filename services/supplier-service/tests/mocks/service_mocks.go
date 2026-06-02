package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models/relations"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/repository"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/service"
)

type JWTUseCaseMock struct{ mock.Mock }

func (m *JWTUseCaseMock) GenerateToken(userID string, role string) (string, error) {
	args := m.Called(userID, role)
	return args.String(0), args.Error(1)
}

type SupplierRepositoryMock struct{ mock.Mock }

func (m *SupplierRepositoryMock) Create(ctx context.Context, supplier models.Supplier) error {
	args := m.Called(ctx, supplier)
	return args.Error(0)
}

func (m *SupplierRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (models.Supplier, error) {
	args := m.Called(ctx, id)
	var result models.Supplier
	if value := args.Get(0); value != nil {
		result = value.(models.Supplier)
	}
	return result, args.Error(1)
}

func (m *SupplierRepositoryMock) GetByLogin(ctx context.Context, login string) (models.Supplier, error) {
	args := m.Called(ctx, login)
	var result models.Supplier
	if value := args.Get(0); value != nil {
		result = value.(models.Supplier)
	}
	return result, args.Error(1)
}

func (m *SupplierRepositoryMock) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	args := m.Called(ctx, login)
	return args.Bool(0), args.Error(1)
}

func (m *SupplierRepositoryMock) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

type SupplierUseCaseMock struct{ mock.Mock }

func (m *SupplierUseCaseMock) Register(ctx context.Context, req dto.RegisterSupplierRequest) (dto.AuthSupplierResponse, error) {
	args := m.Called(ctx, req)
	var result dto.AuthSupplierResponse
	if value := args.Get(0); value != nil {
		result = value.(dto.AuthSupplierResponse)
	}
	return result, args.Error(1)
}

func (m *SupplierUseCaseMock) Login(ctx context.Context, req dto.LoginSupplierRequest) (dto.AuthSupplierResponse, error) {
	args := m.Called(ctx, req)
	var result dto.AuthSupplierResponse
	if value := args.Get(0); value != nil {
		result = value.(dto.AuthSupplierResponse)
	}
	return result, args.Error(1)
}

func (m *SupplierUseCaseMock) GetProfile(ctx context.Context, supplierID uuid.UUID) (dto.SupplierResponse, error) {
	args := m.Called(ctx, supplierID)
	var result dto.SupplierResponse
	if value := args.Get(0); value != nil {
		result = value.(dto.SupplierResponse)
	}
	return result, args.Error(1)
}

type ModuleUseCaseMock struct{ mock.Mock }

func (m *ModuleUseCaseMock) GetList(ctx context.Context) ([]dto.ModuleResponse, error) {
	args := m.Called(ctx)
	var result []dto.ModuleResponse
	if value := args.Get(0); value != nil {
		result = value.([]dto.ModuleResponse)
	}
	return result, args.Error(1)
}

func (m *ModuleUseCaseMock) GetByID(ctx context.Context, id uuid.UUID) (dto.ModuleResponse, error) {
	args := m.Called(ctx, id)
	var result dto.ModuleResponse
	if value := args.Get(0); value != nil {
		result = value.(dto.ModuleResponse)
	}
	return result, args.Error(1)
}

func (m *ModuleUseCaseMock) Create(ctx context.Context, req dto.ModuleRequest) (uuid.UUID, error) {
	args := m.Called(ctx, req)
	var result uuid.UUID
	if value := args.Get(0); value != nil {
		result = value.(uuid.UUID)
	}
	return result, args.Error(1)
}

func (m *ModuleUseCaseMock) Update(ctx context.Context, id uuid.UUID, req dto.ModuleRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *ModuleUseCaseMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ModuleUseCaseMock) ReserveModules(ctx context.Context, items []service.ModuleItem) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *ModuleUseCaseMock) ReleaseModules(ctx context.Context, items []service.ModuleItem) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

type MaterialUseCaseMock struct{ mock.Mock }

func (m *MaterialUseCaseMock) GetList(ctx context.Context) ([]dto.MaterialResponse, error) {
	args := m.Called(ctx)
	var result []dto.MaterialResponse
	if value := args.Get(0); value != nil {
		result = value.([]dto.MaterialResponse)
	}
	return result, args.Error(1)
}

func (m *MaterialUseCaseMock) GetByID(ctx context.Context, id uuid.UUID) (dto.MaterialResponse, error) {
	args := m.Called(ctx, id)
	var result dto.MaterialResponse
	if value := args.Get(0); value != nil {
		result = value.(dto.MaterialResponse)
	}
	return result, args.Error(1)
}

func (m *MaterialUseCaseMock) Create(ctx context.Context, req dto.MaterialRequest) (uuid.UUID, error) {
	args := m.Called(ctx, req)
	var result uuid.UUID
	if value := args.Get(0); value != nil {
		result = value.(uuid.UUID)
	}
	return result, args.Error(1)
}

func (m *MaterialUseCaseMock) Update(ctx context.Context, id uuid.UUID, req dto.MaterialRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MaterialUseCaseMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type SupplyUseCaseMock struct{ mock.Mock }

func (m *SupplyUseCaseMock) GetList(ctx context.Context) ([]dto.SupplyResponse, error) {
	args := m.Called(ctx)
	var result []dto.SupplyResponse
	if value := args.Get(0); value != nil {
		result = value.([]dto.SupplyResponse)
	}
	return result, args.Error(1)
}

func (m *SupplyUseCaseMock) GetByID(ctx context.Context, id uuid.UUID) (dto.SupplyResponse, error) {
	args := m.Called(ctx, id)
	var result dto.SupplyResponse
	if value := args.Get(0); value != nil {
		result = value.(dto.SupplyResponse)
	}
	return result, args.Error(1)
}

func (m *SupplyUseCaseMock) Create(ctx context.Context, supplierID uuid.UUID, req dto.SupplyRequest) (uuid.UUID, error) {
	args := m.Called(ctx, supplierID, req)
	var result uuid.UUID
	if value := args.Get(0); value != nil {
		result = value.(uuid.UUID)
	}
	return result, args.Error(1)
}

func (m *SupplyUseCaseMock) Update(ctx context.Context, id uuid.UUID, req dto.SupplyRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *SupplyUseCaseMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type ModuleRepositoryMock struct{ mock.Mock }

func (m *ModuleRepositoryMock) GetList(ctx context.Context) ([]models.Module, error) {
	args := m.Called(ctx)
	var result []models.Module
	if value := args.Get(0); value != nil {
		result = value.([]models.Module)
	}
	return result, args.Error(1)
}

func (m *ModuleRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (models.Module, error) {
	args := m.Called(ctx, id)
	var result models.Module
	if value := args.Get(0); value != nil {
		result = value.(models.Module)
	}
	return result, args.Error(1)
}

func (m *ModuleRepositoryMock) GetByIDTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) (models.Module, error) {
	args := m.Called(ctx, db, id)
	var result models.Module
	if value := args.Get(0); value != nil {
		result = value.(models.Module)
	}
	return result, args.Error(1)
}

func (m *ModuleRepositoryMock) GetByCode(ctx context.Context, code string) (models.Module, error) {
	args := m.Called(ctx, code)
	var result models.Module
	if value := args.Get(0); value != nil {
		result = value.(models.Module)
	}
	return result, args.Error(1)
}

func (m *ModuleRepositoryMock) GetByCodeTx(ctx context.Context, db repository.DBExecutor, code string) (models.Module, error) {
	args := m.Called(ctx, db, code)
	var result models.Module
	if value := args.Get(0); value != nil {
		result = value.(models.Module)
	}
	return result, args.Error(1)
}

func (m *ModuleRepositoryMock) ExistsByCode(ctx context.Context, code string) (bool, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Error(1)
}

func (m *ModuleRepositoryMock) Create(ctx context.Context, item models.Module) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *ModuleRepositoryMock) CreateTx(ctx context.Context, db repository.DBExecutor, item models.Module) error {
	args := m.Called(ctx, db, item)
	return args.Error(0)
}

func (m *ModuleRepositoryMock) Update(ctx context.Context, item models.Module) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *ModuleRepositoryMock) UpdateTx(ctx context.Context, db repository.DBExecutor, item models.Module) error {
	args := m.Called(ctx, db, item)
	return args.Error(0)
}

func (m *ModuleRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ModuleRepositoryMock) DeleteTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) error {
	args := m.Called(ctx, db, id)
	return args.Error(0)
}

func (m *ModuleRepositoryMock) IncreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	args := m.Called(ctx, id, count)
	return args.Error(0)
}

func (m *ModuleRepositoryMock) IncreaseStockTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID, count int) error {
	args := m.Called(ctx, db, id, count)
	return args.Error(0)
}

func (m *ModuleRepositoryMock) DecreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	args := m.Called(ctx, id, count)
	return args.Error(0)
}

func (m *ModuleRepositoryMock) DecreaseStockTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID, count int) error {
	args := m.Called(ctx, db, id, count)
	return args.Error(0)
}

type ModuleMaterialRepositoryMock struct{ mock.Mock }

func (m *ModuleMaterialRepositoryMock) CreateMany(ctx context.Context, items []relations.ModuleMaterial) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *ModuleMaterialRepositoryMock) CreateManyTx(ctx context.Context, db repository.DBExecutor, items []relations.ModuleMaterial) error {
	args := m.Called(ctx, db, items)
	return args.Error(0)
}

func (m *ModuleMaterialRepositoryMock) GetByModuleID(ctx context.Context, id uuid.UUID) ([]relations.ModuleMaterial, error) {
	args := m.Called(ctx, id)
	var result []relations.ModuleMaterial
	if value := args.Get(0); value != nil {
		result = value.([]relations.ModuleMaterial)
	}
	return result, args.Error(1)
}

func (m *ModuleMaterialRepositoryMock) GetByModuleIDTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) ([]relations.ModuleMaterial, error) {
	args := m.Called(ctx, db, id)
	var result []relations.ModuleMaterial
	if value := args.Get(0); value != nil {
		result = value.([]relations.ModuleMaterial)
	}
	return result, args.Error(1)
}

func (m *ModuleMaterialRepositoryMock) DeleteByModuleID(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ModuleMaterialRepositoryMock) DeleteByModuleIDTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) error {
	args := m.Called(ctx, db, id)
	return args.Error(0)
}

type MaterialRepositoryMock struct{ mock.Mock }

func (m *MaterialRepositoryMock) GetList(ctx context.Context) ([]models.Material, error) {
	args := m.Called(ctx)
	var result []models.Material
	if value := args.Get(0); value != nil {
		result = value.([]models.Material)
	}
	return result, args.Error(1)
}

func (m *MaterialRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (models.Material, error) {
	args := m.Called(ctx, id)
	var result models.Material
	if value := args.Get(0); value != nil {
		result = value.(models.Material)
	}
	return result, args.Error(1)
}

func (m *MaterialRepositoryMock) GetByIDTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) (models.Material, error) {
	args := m.Called(ctx, db, id)
	var result models.Material
	if value := args.Get(0); value != nil {
		result = value.(models.Material)
	}
	return result, args.Error(1)
}

func (m *MaterialRepositoryMock) ExistsByCode(ctx context.Context, code string) (bool, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Error(1)
}

func (m *MaterialRepositoryMock) Create(ctx context.Context, item models.Material) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MaterialRepositoryMock) CreateTx(ctx context.Context, db repository.DBExecutor, item models.Material) error {
	args := m.Called(ctx, db, item)
	return args.Error(0)
}

func (m *MaterialRepositoryMock) Update(ctx context.Context, item models.Material) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MaterialRepositoryMock) UpdateTx(ctx context.Context, db repository.DBExecutor, item models.Material) error {
	args := m.Called(ctx, db, item)
	return args.Error(0)
}

func (m *MaterialRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MaterialRepositoryMock) DeleteTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) error {
	args := m.Called(ctx, db, id)
	return args.Error(0)
}

func (m *MaterialRepositoryMock) IncreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	args := m.Called(ctx, id, count)
	return args.Error(0)
}

func (m *MaterialRepositoryMock) IncreaseStockTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID, count int) error {
	args := m.Called(ctx, db, id, count)
	return args.Error(0)
}

func (m *MaterialRepositoryMock) DecreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	args := m.Called(ctx, id, count)
	return args.Error(0)
}

func (m *MaterialRepositoryMock) DecreaseStockTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID, count int) error {
	args := m.Called(ctx, db, id, count)
	return args.Error(0)
}

type SupplyRepositoryMock struct{ mock.Mock }

func (m *SupplyRepositoryMock) GetList(ctx context.Context) ([]models.Supply, error) {
	args := m.Called(ctx)
	var result []models.Supply
	if value := args.Get(0); value != nil {
		result = value.([]models.Supply)
	}
	return result, args.Error(1)
}

func (m *SupplyRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (models.Supply, error) {
	args := m.Called(ctx, id)
	var result models.Supply
	if value := args.Get(0); value != nil {
		result = value.(models.Supply)
	}
	return result, args.Error(1)
}

func (m *SupplyRepositoryMock) GetByIDTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) (models.Supply, error) {
	args := m.Called(ctx, db, id)
	var result models.Supply
	if value := args.Get(0); value != nil {
		result = value.(models.Supply)
	}
	return result, args.Error(1)
}

func (m *SupplyRepositoryMock) ExistsByCode(ctx context.Context, code string) (bool, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Error(1)
}

func (m *SupplyRepositoryMock) Create(ctx context.Context, item models.Supply) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *SupplyRepositoryMock) CreateTx(ctx context.Context, db repository.DBExecutor, item models.Supply) error {
	args := m.Called(ctx, db, item)
	return args.Error(0)
}

func (m *SupplyRepositoryMock) Update(ctx context.Context, item models.Supply) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *SupplyRepositoryMock) UpdateTx(ctx context.Context, db repository.DBExecutor, item models.Supply) error {
	args := m.Called(ctx, db, item)
	return args.Error(0)
}

func (m *SupplyRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *SupplyRepositoryMock) DeleteTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) error {
	args := m.Called(ctx, db, id)
	return args.Error(0)
}

type SupplyMaterialRepositoryMock struct{ mock.Mock }

func (m *SupplyMaterialRepositoryMock) CreateMany(ctx context.Context, items []relations.SupplyMaterial) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *SupplyMaterialRepositoryMock) CreateManyTx(ctx context.Context, db repository.DBExecutor, items []relations.SupplyMaterial) error {
	args := m.Called(ctx, db, items)
	return args.Error(0)
}

func (m *SupplyMaterialRepositoryMock) GetBySupplyID(ctx context.Context, id uuid.UUID) ([]relations.SupplyMaterial, error) {
	args := m.Called(ctx, id)
	var result []relations.SupplyMaterial
	if value := args.Get(0); value != nil {
		result = value.([]relations.SupplyMaterial)
	}
	return result, args.Error(1)
}

func (m *SupplyMaterialRepositoryMock) GetBySupplyIDTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) ([]relations.SupplyMaterial, error) {
	args := m.Called(ctx, db, id)
	var result []relations.SupplyMaterial
	if value := args.Get(0); value != nil {
		result = value.([]relations.SupplyMaterial)
	}
	return result, args.Error(1)
}

func (m *SupplyMaterialRepositoryMock) DeleteBySupplyID(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *SupplyMaterialRepositoryMock) DeleteBySupplyIDTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) error {
	args := m.Called(ctx, db, id)
	return args.Error(0)
}

type TransactionManagerMock struct {
	mock.Mock
	Tx repository.DBExecutor
}

func (m *TransactionManagerMock) RunInTx(ctx context.Context, fn func(ctx context.Context, tx repository.DBExecutor) error) error {
	args := m.Called(ctx)
	if args.Error(0) != nil {
		return args.Error(0)
	}

	tx := m.Tx
	if tx == nil {
		tx = &noopDBExecutor{}
	}

	return fn(ctx, tx)
}

type noopDBExecutor struct{}

func (n *noopDBExecutor) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (n *noopDBExecutor) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return pgx.Rows(nil), nil
}

func (n *noopDBExecutor) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return nil
}
