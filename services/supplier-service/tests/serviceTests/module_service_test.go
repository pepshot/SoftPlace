package servicetests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/service"
	"github.com/pepshot/SoftPlace/services/supplier-service/tests/mocks"
)

func TestModuleServiceGetListSuccess(t *testing.T) {
	moduleRepo := &mocks.ModuleRepositoryMock{}
	svc := service.NewModuleService(moduleRepo, &mocks.ModuleMaterialRepositoryMock{}, &mocks.MaterialRepositoryMock{}, newTxManager(), testLogger(t))
	moduleRepo.On("GetList", mock.Anything).Return([]models.Module{sampleModule(uuid.New())}, nil)

	result, err := svc.GetList(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestModuleServiceCreateEmptyComposition(t *testing.T) {
	svc := service.NewModuleService(&mocks.ModuleRepositoryMock{}, &mocks.ModuleMaterialRepositoryMock{}, &mocks.MaterialRepositoryMock{}, newTxManager(), testLogger(t))

	_, err := svc.Create(context.Background(), dto.ModuleRequest{})
	require.ErrorIs(t, err, service.ErrEmptyComposition)
}

func TestModuleServiceCreateInvalidCount(t *testing.T) {
	moduleRepo := &mocks.ModuleRepositoryMock{}
	materialRepo := &mocks.MaterialRepositoryMock{}
	moduleMaterialRepo := &mocks.ModuleMaterialRepositoryMock{}
	tx := newTxManager()
	svc := service.NewModuleService(moduleRepo, moduleMaterialRepo, materialRepo, tx, testLogger(t))
	req := moduleRequest()
	req.Materials[0].Count = 0

	moduleRepo.On("ExistsByCode", mock.Anything, req.Code).Return(false, nil)
	tx.On("RunInTx", mock.Anything).Return(nil)

	_, err := svc.Create(context.Background(), req)
	require.ErrorIs(t, err, service.ErrInvalidCount)
}

func TestModuleServiceCreateCodeAlreadyUsed(t *testing.T) {
	moduleRepo := &mocks.ModuleRepositoryMock{}
	svc := service.NewModuleService(moduleRepo, &mocks.ModuleMaterialRepositoryMock{}, &mocks.MaterialRepositoryMock{}, newTxManager(), testLogger(t))
	req := moduleRequest()
	moduleRepo.On("ExistsByCode", mock.Anything, req.Code).Return(true, nil)

	_, err := svc.Create(context.Background(), req)
	require.ErrorIs(t, err, service.ErrCodeAlreadyUsed)
}

func TestModuleServiceCreateSuccess(t *testing.T) {
	moduleRepo := &mocks.ModuleRepositoryMock{}
	materialRepo := &mocks.MaterialRepositoryMock{}
	moduleMaterialRepo := &mocks.ModuleMaterialRepositoryMock{}
	tx := newTxManager()
	svc := service.NewModuleService(moduleRepo, moduleMaterialRepo, materialRepo, tx, testLogger(t))
	req := moduleRequest()
	materialID := uuid.MustParse(req.Materials[0].MaterialID)

	moduleRepo.On("ExistsByCode", mock.Anything, req.Code).Return(false, nil)
	tx.On("RunInTx", mock.Anything).Return(nil)
	materialRepo.On("GetByIDTx", mock.Anything, mock.Anything, materialID).Return(sampleMaterial(materialID), nil)
	materialRepo.On("DecreaseStockTx", mock.Anything, mock.Anything, materialID, req.Materials[0].Count).Return(nil)
	moduleRepo.On("CreateTx", mock.Anything, mock.Anything, mock.AnythingOfType("models.Module")).Return(nil)
	moduleMaterialRepo.On("CreateManyTx", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	id, err := svc.Create(context.Background(), req)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, id)
}

func TestModuleServiceReserveInvalidComposition(t *testing.T) {
	svc := service.NewModuleService(&mocks.ModuleRepositoryMock{}, &mocks.ModuleMaterialRepositoryMock{}, &mocks.MaterialRepositoryMock{}, newTxManager(), testLogger(t))
	err := svc.ReserveModules(context.Background(), nil)
	require.ErrorIs(t, err, service.ErrInvalidComposition)
}

func TestModuleServiceDeleteError(t *testing.T) {
	moduleRepo := &mocks.ModuleRepositoryMock{}
	moduleMaterialRepo := &mocks.ModuleMaterialRepositoryMock{}
	materialRepo := &mocks.MaterialRepositoryMock{}
	tx := newTxManager()
	svc := service.NewModuleService(moduleRepo, moduleMaterialRepo, materialRepo, tx, testLogger(t))
	moduleID := uuid.New()

	tx.On("RunInTx", mock.Anything).Return(nil)
	moduleMaterialRepo.On("GetByModuleIDTx", mock.Anything, mock.Anything, moduleID).Return(nil, errors.New("relations error"))

	err := svc.Delete(context.Background(), moduleID)
	require.Error(t, err)
}
