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

func TestGarnitureServiceGetListSuccess(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	garnitureFurnitureRepo := &mocks.GarnitureFurnitureRepositoryMock{}
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewGarnitureService(garnitureRepo, garnitureFurnitureRepo, furnitureRepo, txManager, testLogger(t))

	items := []model.Garniture{{ID: uuid.New(), Name: "Set", Code: "GS-1", Price: 100, StockCount: 1}}
	garnitureRepo.On("GetList", mock.Anything).Return(items, nil)

	result, err := svc.GetList(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestGarnitureServiceGetByIDNotFound(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	svc := service.NewGarnitureService(garnitureRepo, &mocks.GarnitureFurnitureRepositoryMock{}, &mocks.FurnitureRepositoryMock{}, newTxManager(), testLogger(t))

	garnitureRepo.On("GetByID", mock.Anything, mock.Anything).Return(model.Garniture{}, errors.New("not found"))

	_, err := svc.GetByID(context.Background(), uuid.New())
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestGarnitureServiceGetByIDSuccess(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	svc := service.NewGarnitureService(garnitureRepo, &mocks.GarnitureFurnitureRepositoryMock{}, &mocks.FurnitureRepositoryMock{}, newTxManager(), testLogger(t))

	id := uuid.New()
	garnitureRepo.On("GetByID", mock.Anything, id).Return(model.Garniture{ID: id, Name: "Set", Code: "GS-1", Price: 100, StockCount: 1}, nil)

	result, err := svc.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id.String(), result.ID)
}

func TestGarnitureServiceCreateEmpty(t *testing.T) {
	svc := service.NewGarnitureService(&mocks.GarnitureRepositoryMock{}, &mocks.GarnitureFurnitureRepositoryMock{}, &mocks.FurnitureRepositoryMock{}, newTxManager(), testLogger(t))

	_, err := svc.Create(context.Background(), dto.GarnitureRequest{})
	require.ErrorIs(t, err, service.ErrEmptyComposition)
}

func TestGarnitureServiceCreateInvalidCount(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	garnitureFurnitureRepo := &mocks.GarnitureFurnitureRepositoryMock{}
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewGarnitureService(garnitureRepo, garnitureFurnitureRepo, furnitureRepo, txManager, testLogger(t))

	req := garnitureRequest()
	req.Furniture[0].Count = 0
	txManager.On("RunInTx", mock.Anything).Return(nil)

	_, err := svc.Create(context.Background(), req)
	require.ErrorIs(t, err, service.ErrInvalidCount)
}

func TestGarnitureServiceCreateInvalidID(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	garnitureFurnitureRepo := &mocks.GarnitureFurnitureRepositoryMock{}
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewGarnitureService(garnitureRepo, garnitureFurnitureRepo, furnitureRepo, txManager, testLogger(t))

	req := garnitureRequest()
	req.Furniture[0].FurnitureID = "bad"
	txManager.On("RunInTx", mock.Anything).Return(nil)

	_, err := svc.Create(context.Background(), req)
	require.ErrorIs(t, err, service.ErrInvalidID)
}

func TestGarnitureServiceCreateNotFound(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	garnitureFurnitureRepo := &mocks.GarnitureFurnitureRepositoryMock{}
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewGarnitureService(garnitureRepo, garnitureFurnitureRepo, furnitureRepo, txManager, testLogger(t))

	req := garnitureRequest()
	furnitureID := uuid.MustParse(req.Furniture[0].FurnitureID)
	furnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, furnitureID).Return(model.Furniture{}, errors.New("not found"))
	txManager.On("RunInTx", mock.Anything).Return(nil)

	_, err := svc.Create(context.Background(), req)
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestGarnitureServiceCreateNotEnoughStock(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	garnitureFurnitureRepo := &mocks.GarnitureFurnitureRepositoryMock{}
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewGarnitureService(garnitureRepo, garnitureFurnitureRepo, furnitureRepo, txManager, testLogger(t))

	req := garnitureRequest()
	furnitureID := uuid.MustParse(req.Furniture[0].FurnitureID)
	furnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, furnitureID).Return(model.Furniture{ID: furnitureID, Price: 10, StockCount: 0}, nil)
	txManager.On("RunInTx", mock.Anything).Return(nil)

	_, err := svc.Create(context.Background(), req)
	require.ErrorIs(t, err, service.ErrNotEnoughStock)
}

func TestGarnitureServiceCreateDecreaseStockError(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	garnitureFurnitureRepo := &mocks.GarnitureFurnitureRepositoryMock{}
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewGarnitureService(garnitureRepo, garnitureFurnitureRepo, furnitureRepo, txManager, testLogger(t))

	req := garnitureRequest()
	furnitureID := uuid.MustParse(req.Furniture[0].FurnitureID)
	furnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, furnitureID).Return(model.Furniture{ID: furnitureID, Price: 10, StockCount: 5}, nil)
	furnitureRepo.On("DecreaseStockTx", mock.Anything, mock.Anything, furnitureID, req.Furniture[0].Count).Return(errors.New("decrease error"))
	txManager.On("RunInTx", mock.Anything).Return(nil)

	_, err := svc.Create(context.Background(), req)
	require.ErrorIs(t, err, service.ErrNotEnoughStock)
}

func TestGarnitureServiceCreateSuccess(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	garnitureFurnitureRepo := &mocks.GarnitureFurnitureRepositoryMock{}
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewGarnitureService(garnitureRepo, garnitureFurnitureRepo, furnitureRepo, txManager, testLogger(t))

	req := garnitureRequest()
	furnitureID := uuid.MustParse(req.Furniture[0].FurnitureID)
	furnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, furnitureID).Return(model.Furniture{ID: furnitureID, Price: 10, StockCount: 5}, nil)
	furnitureRepo.On("DecreaseStockTx", mock.Anything, mock.Anything, furnitureID, req.Furniture[0].Count).Return(nil)

	garnitureRepo.On("CreateTx", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	garnitureFurnitureRepo.On("CreateManyTx", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	txManager.On("RunInTx", mock.Anything).Return(nil)

	id, err := svc.Create(context.Background(), req)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, id)
}

func TestGarnitureServiceUpdateEmpty(t *testing.T) {
	svc := service.NewGarnitureService(&mocks.GarnitureRepositoryMock{}, &mocks.GarnitureFurnitureRepositoryMock{}, &mocks.FurnitureRepositoryMock{}, newTxManager(), testLogger(t))

	err := svc.Update(context.Background(), uuid.New(), dto.GarnitureRequest{})
	require.ErrorIs(t, err, service.ErrEmptyComposition)
}

func TestGarnitureServiceUpdateInvalidCount(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	garnitureFurnitureRepo := &mocks.GarnitureFurnitureRepositoryMock{}
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewGarnitureService(garnitureRepo, garnitureFurnitureRepo, furnitureRepo, txManager, testLogger(t))

	req := garnitureRequest()
	req.Furniture[0].Count = 0
	id := uuid.New()
	txManager.On("RunInTx", mock.Anything).Return(nil)
	garnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, id).Return(model.Garniture{ID: id, StockCount: 1}, nil)
	garnitureFurnitureRepo.On("GetByGarnitureIDTx", mock.Anything, mock.Anything, id).Return(nil, nil)

	err := svc.Update(context.Background(), id, req)
	require.ErrorIs(t, err, service.ErrInvalidCount)
}

func TestGarnitureServiceUpdateInvalidID(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	garnitureFurnitureRepo := &mocks.GarnitureFurnitureRepositoryMock{}
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewGarnitureService(garnitureRepo, garnitureFurnitureRepo, furnitureRepo, txManager, testLogger(t))

	req := garnitureRequest()
	req.Furniture[0].FurnitureID = "bad"
	id := uuid.New()
	txManager.On("RunInTx", mock.Anything).Return(nil)
	garnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, id).Return(model.Garniture{ID: id, StockCount: 1}, nil)
	garnitureFurnitureRepo.On("GetByGarnitureIDTx", mock.Anything, mock.Anything, id).Return(nil, nil)

	err := svc.Update(context.Background(), id, req)
	require.ErrorIs(t, err, service.ErrInvalidID)
}

func TestGarnitureServiceUpdateRelationsError(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	garnitureFurnitureRepo := &mocks.GarnitureFurnitureRepositoryMock{}
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewGarnitureService(garnitureRepo, garnitureFurnitureRepo, furnitureRepo, txManager, testLogger(t))

	id := uuid.New()
	garnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, id).Return(model.Garniture{ID: id, StockCount: 1}, nil)
	garnitureFurnitureRepo.On("GetByGarnitureIDTx", mock.Anything, mock.Anything, id).Return(nil, errors.New("relations error"))
	txManager.On("RunInTx", mock.Anything).Return(nil)

	err := svc.Update(context.Background(), id, garnitureRequest())
	require.Error(t, err)
}

func TestGarnitureServiceUpdateNotEnoughStock(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	garnitureFurnitureRepo := &mocks.GarnitureFurnitureRepositoryMock{}
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewGarnitureService(garnitureRepo, garnitureFurnitureRepo, furnitureRepo, txManager, testLogger(t))

	id := uuid.New()
	req := garnitureRequest()
	req.Furniture[0].Count = 2
	furnitureID := uuid.MustParse(req.Furniture[0].FurnitureID)

	garnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, id).Return(model.Garniture{ID: id, StockCount: 1}, nil)
	garnitureFurnitureRepo.On("GetByGarnitureIDTx", mock.Anything, mock.Anything, id).Return([]relations.GarnitureFurniture{{GarnitureID: id, FurnitureID: furnitureID, Count: 1}}, nil)
	furnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, furnitureID).Return(model.Furniture{ID: furnitureID, Price: 10, StockCount: 2}, nil).Once()
	furnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, furnitureID).Return(model.Furniture{ID: furnitureID, Price: 10, StockCount: 0}, nil).Once()
	txManager.On("RunInTx", mock.Anything).Return(nil)

	err := svc.Update(context.Background(), id, req)
	require.ErrorIs(t, err, service.ErrNotEnoughStock)
}

func TestGarnitureServiceUpdateNotFound(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	garnitureFurnitureRepo := &mocks.GarnitureFurnitureRepositoryMock{}
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewGarnitureService(garnitureRepo, garnitureFurnitureRepo, furnitureRepo, txManager, testLogger(t))

	id := uuid.New()
	garnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, id).Return(model.Garniture{}, errors.New("not found"))
	txManager.On("RunInTx", mock.Anything).Return(nil)

	err := svc.Update(context.Background(), id, garnitureRequest())
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestGarnitureServiceDeleteSuccess(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	garnitureFurnitureRepo := &mocks.GarnitureFurnitureRepositoryMock{}
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewGarnitureService(garnitureRepo, garnitureFurnitureRepo, furnitureRepo, txManager, testLogger(t))

	garnitureID := uuid.New()
	relationsItems := []relations.GarnitureFurniture{{GarnitureID: garnitureID, FurnitureID: uuid.New(), Count: 1}}

	garnitureFurnitureRepo.On("GetByGarnitureIDTx", mock.Anything, mock.Anything, garnitureID).Return(relationsItems, nil)
	furnitureRepo.On("IncreaseStockTx", mock.Anything, mock.Anything, relationsItems[0].FurnitureID, relationsItems[0].Count).Return(nil)
	garnitureFurnitureRepo.On("DeleteByGarnitureIDTx", mock.Anything, mock.Anything, garnitureID).Return(nil)
	garnitureRepo.On("DeleteTx", mock.Anything, mock.Anything, garnitureID).Return(nil)
	txManager.On("RunInTx", mock.Anything).Return(nil)

	err := svc.Delete(context.Background(), garnitureID)
	require.NoError(t, err)
}

func TestGarnitureServiceDeleteRelationsError(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	garnitureFurnitureRepo := &mocks.GarnitureFurnitureRepositoryMock{}
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewGarnitureService(garnitureRepo, garnitureFurnitureRepo, furnitureRepo, txManager, testLogger(t))

	id := uuid.New()
	txManager.On("RunInTx", mock.Anything).Return(nil)
	garnitureFurnitureRepo.On("GetByGarnitureIDTx", mock.Anything, mock.Anything, id).Return(nil, errors.New("relations error"))

	err := svc.Delete(context.Background(), id)
	require.Error(t, err)
}

func TestGarnitureServiceDeleteTxError(t *testing.T) {
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	garnitureFurnitureRepo := &mocks.GarnitureFurnitureRepositoryMock{}
	furnitureRepo := &mocks.FurnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewGarnitureService(garnitureRepo, garnitureFurnitureRepo, furnitureRepo, txManager, testLogger(t))

	id := uuid.New()
	txManager.On("RunInTx", mock.Anything).Return(nil)
	garnitureFurnitureRepo.On("GetByGarnitureIDTx", mock.Anything, mock.Anything, id).Return([]relations.GarnitureFurniture{{GarnitureID: id, FurnitureID: uuid.New(), Count: 1}}, nil)
	furnitureRepo.On("IncreaseStockTx", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	garnitureFurnitureRepo.On("DeleteByGarnitureIDTx", mock.Anything, mock.Anything, id).Return(nil)
	garnitureRepo.On("DeleteTx", mock.Anything, mock.Anything, id).Return(errors.New("delete error"))

	err := svc.Delete(context.Background(), id)
	require.Error(t, err)
}
