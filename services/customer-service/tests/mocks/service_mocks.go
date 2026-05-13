package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model/relations"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/repository"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/service"
)

type CustomerRepositoryMock struct {
	mock.Mock
}

func (m *CustomerRepositoryMock) Create(ctx context.Context, customer model.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *CustomerRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (model.Customer, error) {
	args := m.Called(ctx, id)
	var result model.Customer
	if value := args.Get(0); value != nil {
		result = value.(model.Customer)
	}
	return result, args.Error(1)
}

func (m *CustomerRepositoryMock) GetByLogin(ctx context.Context, login string) (model.Customer, error) {
	args := m.Called(ctx, login)
	var result model.Customer
	if value := args.Get(0); value != nil {
		result = value.(model.Customer)
	}
	return result, args.Error(1)
}

func (m *CustomerRepositoryMock) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	args := m.Called(ctx, login)
	return args.Bool(0), args.Error(1)
}

func (m *CustomerRepositoryMock) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

type JWTUseCaseMock struct {
	mock.Mock
}

func (m *JWTUseCaseMock) GenerateToken(userID string, role string) (string, error) {
	args := m.Called(userID, role)
	return args.String(0), args.Error(1)
}

type CustomerUseCaseMock struct {
	mock.Mock
}

func (m *CustomerUseCaseMock) Register(ctx context.Context, req dto.RegisterCustomerRequest) (dto.AuthCustomerResponse, error) {
	args := m.Called(ctx, req)
	var result dto.AuthCustomerResponse
	if value := args.Get(0); value != nil {
		result = value.(dto.AuthCustomerResponse)
	}
	return result, args.Error(1)
}

func (m *CustomerUseCaseMock) Login(ctx context.Context, req dto.LoginCustomerRequest) (dto.AuthCustomerResponse, error) {
	args := m.Called(ctx, req)
	var result dto.AuthCustomerResponse
	if value := args.Get(0); value != nil {
		result = value.(dto.AuthCustomerResponse)
	}
	return result, args.Error(1)
}

func (m *CustomerUseCaseMock) GetProfile(ctx context.Context, customerID uuid.UUID) (dto.CustomerResponse, error) {
	args := m.Called(ctx, customerID)
	var result dto.CustomerResponse
	if value := args.Get(0); value != nil {
		result = value.(dto.CustomerResponse)
	}
	return result, args.Error(1)
}

type FurnitureUseCaseMock struct {
	mock.Mock
}

func (m *FurnitureUseCaseMock) GetList(ctx context.Context) ([]dto.FurnitureResponse, error) {
	args := m.Called(ctx)
	var result []dto.FurnitureResponse
	if value := args.Get(0); value != nil {
		result = value.([]dto.FurnitureResponse)
	}
	return result, args.Error(1)
}

func (m *FurnitureUseCaseMock) GetByID(ctx context.Context, id uuid.UUID) (dto.FurnitureResponse, error) {
	args := m.Called(ctx, id)
	var result dto.FurnitureResponse
	if value := args.Get(0); value != nil {
		result = value.(dto.FurnitureResponse)
	}
	return result, args.Error(1)
}

func (m *FurnitureUseCaseMock) Create(ctx context.Context, req dto.FurnitureRequest) (uuid.UUID, error) {
	args := m.Called(ctx, req)
	var result uuid.UUID
	if value := args.Get(0); value != nil {
		result = value.(uuid.UUID)
	}
	return result, args.Error(1)
}

func (m *FurnitureUseCaseMock) Update(ctx context.Context, id uuid.UUID, req dto.FurnitureRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *FurnitureUseCaseMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type GarnitureUseCaseMock struct {
	mock.Mock
}

func (m *GarnitureUseCaseMock) GetList(ctx context.Context) ([]dto.GarnitureResponse, error) {
	args := m.Called(ctx)
	var result []dto.GarnitureResponse
	if value := args.Get(0); value != nil {
		result = value.([]dto.GarnitureResponse)
	}
	return result, args.Error(1)
}

func (m *GarnitureUseCaseMock) GetByID(ctx context.Context, id uuid.UUID) (dto.GarnitureResponse, error) {
	args := m.Called(ctx, id)
	var result dto.GarnitureResponse
	if value := args.Get(0); value != nil {
		result = value.(dto.GarnitureResponse)
	}
	return result, args.Error(1)
}

func (m *GarnitureUseCaseMock) Create(ctx context.Context, req dto.GarnitureRequest) (uuid.UUID, error) {
	args := m.Called(ctx, req)
	var result uuid.UUID
	if value := args.Get(0); value != nil {
		result = value.(uuid.UUID)
	}
	return result, args.Error(1)
}

func (m *GarnitureUseCaseMock) Update(ctx context.Context, id uuid.UUID, req dto.GarnitureRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *GarnitureUseCaseMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type ShipmentUseCaseMock struct {
	mock.Mock
}

func (m *ShipmentUseCaseMock) GetList(ctx context.Context) ([]dto.ShipmentResponse, error) {
	args := m.Called(ctx)
	var result []dto.ShipmentResponse
	if value := args.Get(0); value != nil {
		result = value.([]dto.ShipmentResponse)
	}
	return result, args.Error(1)
}

func (m *ShipmentUseCaseMock) GetByID(ctx context.Context, id uuid.UUID) (dto.ShipmentResponse, error) {
	args := m.Called(ctx, id)
	var result dto.ShipmentResponse
	if value := args.Get(0); value != nil {
		result = value.(dto.ShipmentResponse)
	}
	return result, args.Error(1)
}

func (m *ShipmentUseCaseMock) Create(ctx context.Context, customerID uuid.UUID, req dto.ShipmentRequest) (uuid.UUID, error) {
	args := m.Called(ctx, customerID, req)
	var result uuid.UUID
	if value := args.Get(0); value != nil {
		result = value.(uuid.UUID)
	}
	return result, args.Error(1)
}

func (m *ShipmentUseCaseMock) Update(ctx context.Context, id uuid.UUID, req dto.ShipmentRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *ShipmentUseCaseMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type FurnitureRepositoryMock struct {
	mock.Mock
}

func (m *FurnitureRepositoryMock) GetList(ctx context.Context) ([]model.Furniture, error) {
	args := m.Called(ctx)
	var result []model.Furniture
	if value := args.Get(0); value != nil {
		result = value.([]model.Furniture)
	}
	return result, args.Error(1)
}

func (m *FurnitureRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (model.Furniture, error) {
	args := m.Called(ctx, id)
	var result model.Furniture
	if value := args.Get(0); value != nil {
		result = value.(model.Furniture)
	}
	return result, args.Error(1)
}

func (m *FurnitureRepositoryMock) GetByIDTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) (model.Furniture, error) {
	args := m.Called(ctx, db, id)
	var result model.Furniture
	if value := args.Get(0); value != nil {
		result = value.(model.Furniture)
	}
	return result, args.Error(1)
}

func (m *FurnitureRepositoryMock) Create(ctx context.Context, furniture model.Furniture) error {
	args := m.Called(ctx, furniture)
	return args.Error(0)
}

func (m *FurnitureRepositoryMock) CreateTx(ctx context.Context, db repository.DBExecutor, furniture model.Furniture) error {
	args := m.Called(ctx, db, furniture)
	return args.Error(0)
}

func (m *FurnitureRepositoryMock) Update(ctx context.Context, furniture model.Furniture) error {
	args := m.Called(ctx, furniture)
	return args.Error(0)
}

func (m *FurnitureRepositoryMock) UpdateTx(ctx context.Context, db repository.DBExecutor, furniture model.Furniture) error {
	args := m.Called(ctx, db, furniture)
	return args.Error(0)
}

func (m *FurnitureRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *FurnitureRepositoryMock) DeleteTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) error {
	args := m.Called(ctx, db, id)
	return args.Error(0)
}

func (m *FurnitureRepositoryMock) IncreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	args := m.Called(ctx, id, count)
	return args.Error(0)
}

func (m *FurnitureRepositoryMock) IncreaseStockTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID, count int) error {
	args := m.Called(ctx, db, id, count)
	return args.Error(0)
}

func (m *FurnitureRepositoryMock) DecreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	args := m.Called(ctx, id, count)
	return args.Error(0)
}

func (m *FurnitureRepositoryMock) DecreaseStockTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID, count int) error {
	args := m.Called(ctx, db, id, count)
	return args.Error(0)
}

type FurnitureModuleRepositoryMock struct {
	mock.Mock
}

func (m *FurnitureModuleRepositoryMock) CreateMany(ctx context.Context, items []relations.FurnitureModule) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *FurnitureModuleRepositoryMock) CreateManyTx(ctx context.Context, db repository.DBExecutor, items []relations.FurnitureModule) error {
	args := m.Called(ctx, db, items)
	return args.Error(0)
}

func (m *FurnitureModuleRepositoryMock) GetByFurnitureID(ctx context.Context, furnitureID uuid.UUID) ([]relations.FurnitureModule, error) {
	args := m.Called(ctx, furnitureID)
	var result []relations.FurnitureModule
	if value := args.Get(0); value != nil {
		result = value.([]relations.FurnitureModule)
	}
	return result, args.Error(1)
}

func (m *FurnitureModuleRepositoryMock) GetByFurnitureIDTx(ctx context.Context, db repository.DBExecutor, furnitureID uuid.UUID) ([]relations.FurnitureModule, error) {
	args := m.Called(ctx, db, furnitureID)
	var result []relations.FurnitureModule
	if value := args.Get(0); value != nil {
		result = value.([]relations.FurnitureModule)
	}
	return result, args.Error(1)
}

func (m *FurnitureModuleRepositoryMock) DeleteByFurnitureID(ctx context.Context, furnitureID uuid.UUID) error {
	args := m.Called(ctx, furnitureID)
	return args.Error(0)
}

func (m *FurnitureModuleRepositoryMock) DeleteByFurnitureIDTx(ctx context.Context, db repository.DBExecutor, furnitureID uuid.UUID) error {
	args := m.Called(ctx, db, furnitureID)
	return args.Error(0)
}

type GarnitureRepositoryMock struct {
	mock.Mock
}

func (m *GarnitureRepositoryMock) GetList(ctx context.Context) ([]model.Garniture, error) {
	args := m.Called(ctx)
	var result []model.Garniture
	if value := args.Get(0); value != nil {
		result = value.([]model.Garniture)
	}
	return result, args.Error(1)
}

func (m *GarnitureRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (model.Garniture, error) {
	args := m.Called(ctx, id)
	var result model.Garniture
	if value := args.Get(0); value != nil {
		result = value.(model.Garniture)
	}
	return result, args.Error(1)
}

func (m *GarnitureRepositoryMock) GetByIDTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) (model.Garniture, error) {
	args := m.Called(ctx, db, id)
	var result model.Garniture
	if value := args.Get(0); value != nil {
		result = value.(model.Garniture)
	}
	return result, args.Error(1)
}

func (m *GarnitureRepositoryMock) Create(ctx context.Context, garniture model.Garniture) error {
	args := m.Called(ctx, garniture)
	return args.Error(0)
}

func (m *GarnitureRepositoryMock) CreateTx(ctx context.Context, db repository.DBExecutor, garniture model.Garniture) error {
	args := m.Called(ctx, db, garniture)
	return args.Error(0)
}

func (m *GarnitureRepositoryMock) Update(ctx context.Context, garniture model.Garniture) error {
	args := m.Called(ctx, garniture)
	return args.Error(0)
}

func (m *GarnitureRepositoryMock) UpdateTx(ctx context.Context, db repository.DBExecutor, garniture model.Garniture) error {
	args := m.Called(ctx, db, garniture)
	return args.Error(0)
}

func (m *GarnitureRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *GarnitureRepositoryMock) DeleteTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) error {
	args := m.Called(ctx, db, id)
	return args.Error(0)
}

func (m *GarnitureRepositoryMock) IncreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	args := m.Called(ctx, id, count)
	return args.Error(0)
}

func (m *GarnitureRepositoryMock) IncreaseStockTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID, count int) error {
	args := m.Called(ctx, db, id, count)
	return args.Error(0)
}

func (m *GarnitureRepositoryMock) DecreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	args := m.Called(ctx, id, count)
	return args.Error(0)
}

func (m *GarnitureRepositoryMock) DecreaseStockTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID, count int) error {
	args := m.Called(ctx, db, id, count)
	return args.Error(0)
}

type GarnitureFurnitureRepositoryMock struct {
	mock.Mock
}

func (m *GarnitureFurnitureRepositoryMock) CreateMany(ctx context.Context, items []relations.GarnitureFurniture) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *GarnitureFurnitureRepositoryMock) CreateManyTx(ctx context.Context, db repository.DBExecutor, items []relations.GarnitureFurniture) error {
	args := m.Called(ctx, db, items)
	return args.Error(0)
}

func (m *GarnitureFurnitureRepositoryMock) GetByGarnitureID(ctx context.Context, garnitureID uuid.UUID) ([]relations.GarnitureFurniture, error) {
	args := m.Called(ctx, garnitureID)
	var result []relations.GarnitureFurniture
	if value := args.Get(0); value != nil {
		result = value.([]relations.GarnitureFurniture)
	}
	return result, args.Error(1)
}

func (m *GarnitureFurnitureRepositoryMock) GetByGarnitureIDTx(ctx context.Context, db repository.DBExecutor, garnitureID uuid.UUID) ([]relations.GarnitureFurniture, error) {
	args := m.Called(ctx, db, garnitureID)
	var result []relations.GarnitureFurniture
	if value := args.Get(0); value != nil {
		result = value.([]relations.GarnitureFurniture)
	}
	return result, args.Error(1)
}

func (m *GarnitureFurnitureRepositoryMock) DeleteByGarnitureID(ctx context.Context, garnitureID uuid.UUID) error {
	args := m.Called(ctx, garnitureID)
	return args.Error(0)
}

func (m *GarnitureFurnitureRepositoryMock) DeleteByGarnitureIDTx(ctx context.Context, db repository.DBExecutor, garnitureID uuid.UUID) error {
	args := m.Called(ctx, db, garnitureID)
	return args.Error(0)
}

type ShipmentRepositoryMock struct {
	mock.Mock
}

func (m *ShipmentRepositoryMock) GetList(ctx context.Context) ([]model.Shipment, error) {
	args := m.Called(ctx)
	var result []model.Shipment
	if value := args.Get(0); value != nil {
		result = value.([]model.Shipment)
	}
	return result, args.Error(1)
}

func (m *ShipmentRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (model.Shipment, error) {
	args := m.Called(ctx, id)
	var result model.Shipment
	if value := args.Get(0); value != nil {
		result = value.(model.Shipment)
	}
	return result, args.Error(1)
}

func (m *ShipmentRepositoryMock) GetByIDTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) (model.Shipment, error) {
	args := m.Called(ctx, db, id)
	var result model.Shipment
	if value := args.Get(0); value != nil {
		result = value.(model.Shipment)
	}
	return result, args.Error(1)
}

func (m *ShipmentRepositoryMock) Create(ctx context.Context, shipment model.Shipment) error {
	args := m.Called(ctx, shipment)
	return args.Error(0)
}

func (m *ShipmentRepositoryMock) CreateTx(ctx context.Context, db repository.DBExecutor, shipment model.Shipment) error {
	args := m.Called(ctx, db, shipment)
	return args.Error(0)
}

func (m *ShipmentRepositoryMock) Update(ctx context.Context, shipment model.Shipment) error {
	args := m.Called(ctx, shipment)
	return args.Error(0)
}

func (m *ShipmentRepositoryMock) UpdateTx(ctx context.Context, db repository.DBExecutor, shipment model.Shipment) error {
	args := m.Called(ctx, db, shipment)
	return args.Error(0)
}

func (m *ShipmentRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ShipmentRepositoryMock) DeleteTx(ctx context.Context, db repository.DBExecutor, id uuid.UUID) error {
	args := m.Called(ctx, db, id)
	return args.Error(0)
}

type ShipmentGarnitureRepositoryMock struct {
	mock.Mock
}

func (m *ShipmentGarnitureRepositoryMock) CreateMany(ctx context.Context, items []relations.ShipmentGarniture) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *ShipmentGarnitureRepositoryMock) CreateManyTx(ctx context.Context, db repository.DBExecutor, items []relations.ShipmentGarniture) error {
	args := m.Called(ctx, db, items)
	return args.Error(0)
}

func (m *ShipmentGarnitureRepositoryMock) GetByShipmentID(ctx context.Context, shipmentID uuid.UUID) ([]relations.ShipmentGarniture, error) {
	args := m.Called(ctx, shipmentID)
	var result []relations.ShipmentGarniture
	if value := args.Get(0); value != nil {
		result = value.([]relations.ShipmentGarniture)
	}
	return result, args.Error(1)
}

func (m *ShipmentGarnitureRepositoryMock) GetByShipmentIDTx(ctx context.Context, db repository.DBExecutor, shipmentID uuid.UUID) ([]relations.ShipmentGarniture, error) {
	args := m.Called(ctx, db, shipmentID)
	var result []relations.ShipmentGarniture
	if value := args.Get(0); value != nil {
		result = value.([]relations.ShipmentGarniture)
	}
	return result, args.Error(1)
}

func (m *ShipmentGarnitureRepositoryMock) DeleteByShipmentID(ctx context.Context, shipmentID uuid.UUID) error {
	args := m.Called(ctx, shipmentID)
	return args.Error(0)
}

func (m *ShipmentGarnitureRepositoryMock) DeleteByShipmentIDTx(ctx context.Context, db repository.DBExecutor, shipmentID uuid.UUID) error {
	args := m.Called(ctx, db, shipmentID)
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

type SupplierClientMock struct {
	mock.Mock
}

func (m *SupplierClientMock) GetModule(ctx context.Context, moduleID string) (service.ModuleInfo, error) {
	args := m.Called(ctx, moduleID)
	var result service.ModuleInfo
	if value := args.Get(0); value != nil {
		result = value.(service.ModuleInfo)
	}
	return result, args.Error(1)
}

func (m *SupplierClientMock) ReserveModules(ctx context.Context, items []service.ModuleItem) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *SupplierClientMock) ReleaseModules(ctx context.Context, items []service.ModuleItem) error {
	args := m.Called(ctx, items)
	return args.Error(0)
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
