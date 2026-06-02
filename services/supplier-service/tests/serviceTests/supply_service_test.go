package servicetests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models/relations"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/service"
	"github.com/pepshot/SoftPlace/services/supplier-service/tests/mocks"
)

func TestSupplyServiceGetListSuccess(t *testing.T) {
	supplyRepo := &mocks.SupplyRepositoryMock{}
	svc := service.NewSupplyService(supplyRepo, &mocks.SupplyMaterialRepositoryMock{}, &mocks.MaterialRepositoryMock{}, newTxManager(), testLogger(t))
	supplyRepo.On("GetList", mock.Anything).Return([]models.Supply{sampleSupply(uuid.New(), uuid.New())}, nil)

	result, err := svc.GetList(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestSupplyServiceCreateInvalidDate(t *testing.T) {
	svc := service.NewSupplyService(&mocks.SupplyRepositoryMock{}, &mocks.SupplyMaterialRepositoryMock{}, &mocks.MaterialRepositoryMock{}, newTxManager(), testLogger(t))
	req := supplyRequest()
	req.Date = "bad"

	_, err := svc.Create(context.Background(), uuid.New(), req)
	require.ErrorIs(t, err, service.ErrInvalidDate)
}

func TestSupplyServiceCreateCodeAlreadyUsed(t *testing.T) {
	supplyRepo := &mocks.SupplyRepositoryMock{}
	svc := service.NewSupplyService(supplyRepo, &mocks.SupplyMaterialRepositoryMock{}, &mocks.MaterialRepositoryMock{}, newTxManager(), testLogger(t))
	req := supplyRequest()
	supplyRepo.On("ExistsByCode", mock.Anything, req.Code).Return(true, nil)

	_, err := svc.Create(context.Background(), uuid.New(), req)
	require.ErrorIs(t, err, service.ErrCodeAlreadyUsed)
}

func TestSupplyServiceCreateSuccess(t *testing.T) {
	supplyRepo := &mocks.SupplyRepositoryMock{}
	supplyMaterialRepo := &mocks.SupplyMaterialRepositoryMock{}
	materialRepo := &mocks.MaterialRepositoryMock{}
	tx := newTxManager()
	svc := service.NewSupplyService(supplyRepo, supplyMaterialRepo, materialRepo, tx, testLogger(t))
	req := supplyRequest()
	materialID := uuid.MustParse(req.Materials[0].MaterialID)
	supplierID := uuid.New()

	supplyRepo.On("ExistsByCode", mock.Anything, req.Code).Return(false, nil)
	tx.On("RunInTx", mock.Anything).Return(nil)
	materialRepo.On("GetByIDTx", mock.Anything, mock.Anything, materialID).Return(sampleMaterial(materialID), nil)
	materialRepo.On("IncreaseStockTx", mock.Anything, mock.Anything, materialID, req.Materials[0].Count).Return(nil)
	supplyRepo.On("CreateTx", mock.Anything, mock.Anything, mock.AnythingOfType("models.Supply")).Return(nil)
	supplyMaterialRepo.On("CreateManyTx", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	id, err := svc.Create(context.Background(), supplierID, req)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, id)
}

func TestSupplyServiceUpdateNotFound(t *testing.T) {
	supplyRepo := &mocks.SupplyRepositoryMock{}
	supplyMaterialRepo := &mocks.SupplyMaterialRepositoryMock{}
	materialRepo := &mocks.MaterialRepositoryMock{}
	tx := newTxManager()
	svc := service.NewSupplyService(supplyRepo, supplyMaterialRepo, materialRepo, tx, testLogger(t))
	supplyID := uuid.New()
	req := supplyRequest()

	tx.On("RunInTx", mock.Anything).Return(nil)
	supplyRepo.On("GetByIDTx", mock.Anything, mock.Anything, supplyID).Return(models.Supply{}, errors.New("not found"))

	err := svc.Update(context.Background(), supplyID, req)
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestSupplyServiceDeleteSuccess(t *testing.T) {
	supplyRepo := &mocks.SupplyRepositoryMock{}
	supplyMaterialRepo := &mocks.SupplyMaterialRepositoryMock{}
	materialRepo := &mocks.MaterialRepositoryMock{}
	tx := newTxManager()
	svc := service.NewSupplyService(supplyRepo, supplyMaterialRepo, materialRepo, tx, testLogger(t))
	supplyID := uuid.New()
	materialID := uuid.New()

	tx.On("RunInTx", mock.Anything).Return(nil)
	supplyMaterialRepo.On("GetBySupplyIDTx", mock.Anything, mock.Anything, supplyID).Return([]relations.SupplyMaterial{{SupplyID: supplyID, MaterialID: materialID, Count: 1}}, nil)
	materialRepo.On("GetByIDTx", mock.Anything, mock.Anything, materialID).Return(sampleMaterial(materialID), nil)
	materialRepo.On("DecreaseStockTx", mock.Anything, mock.Anything, materialID, 1).Return(nil)
	supplyMaterialRepo.On("DeleteBySupplyIDTx", mock.Anything, mock.Anything, supplyID).Return(nil)
	supplyRepo.On("DeleteTx", mock.Anything, mock.Anything, supplyID).Return(nil)

	err := svc.Delete(context.Background(), supplyID)
	require.NoError(t, err)
}
