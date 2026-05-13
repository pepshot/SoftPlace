package servicetests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/pepshot/SoftPlace/services/customer-service/tests/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model/relations"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/service"
)

func TestFurnitureServiceGetListSuccess(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	items := []model.Furniture{{ID: uuid.New(), Name: "Chair", Code: "CH-1", Price: 10, StockCount: 2}}
	furnitureRepo.On("GetList", mock.Anything).Return(items, nil)

	result, err := svc.GetList(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	furnitureRepo.AssertExpectations(t)
}

func TestFurnitureServiceGetListError(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	furnitureRepo.On("GetList", mock.Anything).Return(nil, errors.New("db error"))

	_, err := svc.GetList(context.Background())
	require.Error(t, err)
}

func TestFurnitureServiceGetByIDError(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	itemID := uuid.New()
	furnitureRepo.On("GetByID", mock.Anything, itemID).Return(model.Furniture{}, errors.New("not found"))

	_, err := svc.GetByID(context.Background(), itemID)
	require.Error(t, err)
}

func TestFurnitureServiceCreateEmptyModules(t *testing.T) {
	svc := service.NewFurnitureService(&mocks.FurnitureRepositoryMock{}, &mocks.FurnitureModuleRepositoryMock{}, &mocks.SupplierClientMock{}, newTxManager(), testLogger(t))

	_, err := svc.Create(context.Background(), dto.FurnitureRequest{Name: "Chair", Code: "CH-1"})
	require.ErrorIs(t, err, service.ErrEmptyComposition)
}

func TestFurnitureServiceCreateInvalidCount(t *testing.T) {
	svc := service.NewFurnitureService(&mocks.FurnitureRepositoryMock{}, &mocks.FurnitureModuleRepositoryMock{}, &mocks.SupplierClientMock{}, newTxManager(), testLogger(t))

	req := furnitureRequest()
	req.Modules[0].Count = 0

	_, err := svc.Create(context.Background(), req)
	require.ErrorIs(t, err, service.ErrInvalidCount)
}

func TestFurnitureServiceCreateInvalidModuleID(t *testing.T) {
	svc := service.NewFurnitureService(&mocks.FurnitureRepositoryMock{}, &mocks.FurnitureModuleRepositoryMock{}, &mocks.SupplierClientMock{}, newTxManager(), testLogger(t))

	req := furnitureRequest()
	req.Modules[0].ModuleID = "bad"

	_, err := svc.Create(context.Background(), req)
	require.ErrorIs(t, err, service.ErrInvalidID)
}

func TestFurnitureServiceCreateSupplierError(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	req := furnitureRequest()
	supplier.On("GetModule", mock.Anything, req.Modules[0].ModuleID).Return(service.ModuleInfo{}, errors.New("fail"))

	_, err := svc.Create(context.Background(), req)
	require.Error(t, err)
}

func TestFurnitureServiceCreateNotEnoughStock(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	req := furnitureRequest()
	supplier.On("GetModule", mock.Anything, req.Modules[0].ModuleID).
		Return(moduleInfo(req.Modules[0].ModuleID, 5, 0), nil)

	_, err := svc.Create(context.Background(), req)
	require.ErrorIs(t, err, service.ErrNotEnoughStock)
}

func TestFurnitureServiceCreateReserveError(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	req := furnitureRequest()
	supplier.On("GetModule", mock.Anything, req.Modules[0].ModuleID).
		Return(moduleInfo(req.Modules[0].ModuleID, 5, 10), nil)
	supplier.On("ReserveModules", mock.Anything, mock.Anything).Return(errors.New("reserve error"))

	_, err := svc.Create(context.Background(), req)
	require.Error(t, err)
}

func TestFurnitureServiceCreateTxErrorReleases(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	req := furnitureRequest()
	supplier.On("GetModule", mock.Anything, req.Modules[0].ModuleID).
		Return(moduleInfo(req.Modules[0].ModuleID, 5, 10), nil)
	supplier.On("ReserveModules", mock.Anything, mock.Anything).Return(nil)
	supplier.On("ReleaseModules", mock.Anything, mock.Anything).Return(nil)

	furnitureRepo.On("CreateTx", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("insert error"))
	furnitureModuleRepo.On("CreateManyTx", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	txManager.On("RunInTx", mock.Anything).Return(nil)

	_, err := svc.Create(context.Background(), req)
	require.Error(t, err)
}

func TestFurnitureServiceCreateSuccess(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	req := furnitureRequest()
	supplier.On("GetModule", mock.Anything, req.Modules[0].ModuleID).
		Return(moduleInfo(req.Modules[0].ModuleID, 5, 10), nil)
	supplier.On("ReserveModules", mock.Anything, mock.Anything).Return(nil)

	furnitureRepo.On("CreateTx", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	furnitureModuleRepo.On("CreateManyTx", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	txManager.On("RunInTx", mock.Anything).Return(nil)

	id, err := svc.Create(context.Background(), req)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, id)
}

func TestFurnitureServiceUpdateEmptyModules(t *testing.T) {
	svc := service.NewFurnitureService(&mocks.FurnitureRepositoryMock{}, &mocks.FurnitureModuleRepositoryMock{}, &mocks.SupplierClientMock{}, newTxManager(), testLogger(t))

	err := svc.Update(context.Background(), uuid.New(), dto.FurnitureRequest{})
	require.ErrorIs(t, err, service.ErrEmptyComposition)
}

func TestFurnitureServiceUpdateRelationsError(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	furnitureModuleRepo.On("GetByFurnitureID", mock.Anything, mock.Anything).Return(nil, errors.New("fail"))

	err := svc.Update(context.Background(), uuid.New(), furnitureRequest())
	require.Error(t, err)
}

func TestFurnitureServiceUpdateSuccess(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	itemID := uuid.New()
	req := furnitureRequest()
	moduleUUID := uuid.MustParse(req.Modules[0].ModuleID)

	furnitureModuleRepo.On("GetByFurnitureID", mock.Anything, itemID).Return([]relations.FurnitureModule{{
		FurnitureID: itemID,
		ModuleID:    moduleUUID,
		Count:       1,
	}}, nil)

	supplier.On("GetModule", mock.Anything, req.Modules[0].ModuleID).
		Return(moduleInfo(req.Modules[0].ModuleID, 5, 10), nil)
	supplier.On("ReserveModules", mock.Anything, mock.Anything).Return(nil)
	supplier.On("ReleaseModules", mock.Anything, mock.Anything).Return(nil)

	furnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, itemID).Return(model.Furniture{ID: itemID, StockCount: 1}, nil)
	furnitureRepo.On("UpdateTx", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	furnitureModuleRepo.On("DeleteByFurnitureIDTx", mock.Anything, mock.Anything, itemID).Return(nil)
	furnitureModuleRepo.On("CreateManyTx", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	txManager.On("RunInTx", mock.Anything).Return(nil)

	err := svc.Update(context.Background(), itemID, req)
	require.NoError(t, err)
}

func TestFurnitureServiceDeleteSuccess(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	itemID := uuid.New()
	furnitureModuleRepo.On("GetByFurnitureID", mock.Anything, itemID).Return([]relations.FurnitureModule{{
		FurnitureID: itemID,
		ModuleID:    uuid.New(),
		Count:       1,
	}}, nil)
	furnitureModuleRepo.On("DeleteByFurnitureIDTx", mock.Anything, mock.Anything, itemID).Return(nil)
	furnitureRepo.On("DeleteTx", mock.Anything, mock.Anything, itemID).Return(nil)
	supplier.On("ReleaseModules", mock.Anything, mock.Anything).Return(nil)
	txManager.On("RunInTx", mock.Anything).Return(nil)

	err := svc.Delete(context.Background(), itemID)
	require.NoError(t, err)
}

func TestFurnitureServiceCreateSupplierGetModuleError(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	req := furnitureRequest()
	supplier.On("GetModule", mock.Anything, req.Modules[0].ModuleID).Return(service.ModuleInfo{}, errors.New("supplier error"))

	_, err := svc.Create(context.Background(), req)
	require.Error(t, err)
}

func TestFurnitureServiceUpdateInvalidCount(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	req := furnitureRequest()
	req.Modules[0].Count = 0
	furnitureModuleRepo.On("GetByFurnitureID", mock.Anything, mock.Anything).Return(nil, nil)

	err := svc.Update(context.Background(), uuid.New(), req)
	require.ErrorIs(t, err, service.ErrInvalidCount)
}

func TestFurnitureServiceUpdateInvalidID(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	req := furnitureRequest()
	req.Modules[0].ModuleID = "bad"
	furnitureModuleRepo.On("GetByFurnitureID", mock.Anything, mock.Anything).Return(nil, nil)

	err := svc.Update(context.Background(), uuid.New(), req)
	require.ErrorIs(t, err, service.ErrInvalidID)
}

func TestFurnitureServiceUpdateReserveError(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	itemID := uuid.New()
	req := furnitureRequest()
	moduleID := req.Modules[0].ModuleID

	furnitureModuleRepo.On("GetByFurnitureID", mock.Anything, itemID).Return([]relations.FurnitureModule{{FurnitureID: itemID, ModuleID: uuid.MustParse(moduleID), Count: 1}}, nil)
	supplier.On("GetModule", mock.Anything, moduleID).Return(moduleInfo(moduleID, 5, 10), nil)
	supplier.On("ReserveModules", mock.Anything, mock.Anything).Return(errors.New("reserve error"))

	err := svc.Update(context.Background(), itemID, req)
	require.ErrorIs(t, err, service.ErrExternalService)
}

func TestFurnitureServiceUpdateNotFoundInTx(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	itemID := uuid.New()
	req := furnitureRequest()
	moduleID := req.Modules[0].ModuleID

	furnitureModuleRepo.On("GetByFurnitureID", mock.Anything, itemID).Return([]relations.FurnitureModule{{FurnitureID: itemID, ModuleID: uuid.MustParse(moduleID), Count: 1}}, nil)
	supplier.On("GetModule", mock.Anything, moduleID).Return(moduleInfo(moduleID, 5, 10), nil)
	supplier.On("ReserveModules", mock.Anything, mock.Anything).Return(nil)
	supplier.On("ReleaseModules", mock.Anything, mock.Anything).Return(nil)
	txManager.On("RunInTx", mock.Anything).Return(nil)
	furnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, itemID).Return(model.Furniture{}, errors.New("not found"))

	err := svc.Update(context.Background(), itemID, req)
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestFurnitureServiceDeleteRelationsError(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	itemID := uuid.New()
	furnitureModuleRepo.On("GetByFurnitureID", mock.Anything, itemID).Return(nil, errors.New("relations error"))

	err := svc.Delete(context.Background(), itemID)
	require.Error(t, err)
}

func TestFurnitureServiceDeleteTxError(t *testing.T) {
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	furnitureModuleRepo := &mocks.FurnitureModuleRepositoryMock{}
	supplier := &mocks.SupplierClientMock{}
	txManager := newTxManager()

	svc := service.NewFurnitureService(furnitureRepo, furnitureModuleRepo, supplier, txManager, testLogger(t))

	itemID := uuid.New()
	furnitureModuleRepo.On("GetByFurnitureID", mock.Anything, itemID).Return([]relations.FurnitureModule{{FurnitureID: itemID, ModuleID: uuid.New(), Count: 1}}, nil)
	txManager.On("RunInTx", mock.Anything).Return(nil)
	furnitureModuleRepo.On("DeleteByFurnitureIDTx", mock.Anything, mock.Anything, itemID).Return(errors.New("delete relation error"))

	err := svc.Delete(context.Background(), itemID)
	require.Error(t, err)
}
