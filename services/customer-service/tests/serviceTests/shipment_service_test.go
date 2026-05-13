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

func TestShipmentServiceGetListSuccess(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	shipmentGarnitureRepo := &mocks.ShipmentGarnitureRepositoryMock{}
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewShipmentService(shipmentRepo, shipmentGarnitureRepo, garnitureRepo, txManager, testLogger(t))

	items := []model.Shipment{{ID: uuid.New(), CustomerID: uuid.New(), Code: "SH-1"}}
	shipmentRepo.On("GetList", mock.Anything).Return(items, nil)

	result, err := svc.GetList(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestShipmentServiceGetByIDNotFound(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	svc := service.NewShipmentService(shipmentRepo, &mocks.ShipmentGarnitureRepositoryMock{}, &mocks.GarnitureRepositoryMock{}, newTxManager(), testLogger(t))

	shipmentRepo.On("GetByID", mock.Anything, mock.Anything).Return(model.Shipment{}, errors.New("not found"))

	_, err := svc.GetByID(context.Background(), uuid.New())
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestShipmentServiceGetByIDSuccess(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	svc := service.NewShipmentService(shipmentRepo, &mocks.ShipmentGarnitureRepositoryMock{}, &mocks.GarnitureRepositoryMock{}, newTxManager(), testLogger(t))

	id := uuid.New()
	shipmentRepo.On("GetByID", mock.Anything, id).Return(model.Shipment{ID: id, CustomerID: uuid.New(), Code: "SH-1"}, nil)

	result, err := svc.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id.String(), result.ID)
}

func TestShipmentServiceCreateEmpty(t *testing.T) {
	svc := service.NewShipmentService(&mocks.ShipmentRepositoryMock{}, &mocks.ShipmentGarnitureRepositoryMock{}, &mocks.GarnitureRepositoryMock{}, newTxManager(), testLogger(t))

	_, err := svc.Create(context.Background(), uuid.New(), dto.ShipmentRequest{})
	require.ErrorIs(t, err, service.ErrEmptyComposition)
}

func TestShipmentServiceCreateInvalidDate(t *testing.T) {
	svc := service.NewShipmentService(&mocks.ShipmentRepositoryMock{}, &mocks.ShipmentGarnitureRepositoryMock{}, &mocks.GarnitureRepositoryMock{}, newTxManager(), testLogger(t))

	req := shipmentRequest()
	req.Date = "bad"

	_, err := svc.Create(context.Background(), uuid.New(), req)
	require.ErrorIs(t, err, service.ErrInvalidDate)
}

func TestShipmentServiceCreateInvalidCount(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	shipmentGarnitureRepo := &mocks.ShipmentGarnitureRepositoryMock{}
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewShipmentService(shipmentRepo, shipmentGarnitureRepo, garnitureRepo, txManager, testLogger(t))

	req := shipmentRequest()
	req.Garnitures[0].Count = 0
	txManager.On("RunInTx", mock.Anything).Return(nil)

	_, err := svc.Create(context.Background(), uuid.New(), req)
	require.ErrorIs(t, err, service.ErrInvalidCount)
}

func TestShipmentServiceCreateInvalidID(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	shipmentGarnitureRepo := &mocks.ShipmentGarnitureRepositoryMock{}
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewShipmentService(shipmentRepo, shipmentGarnitureRepo, garnitureRepo, txManager, testLogger(t))

	req := shipmentRequest()
	req.Garnitures[0].GarnitureID = "bad"
	txManager.On("RunInTx", mock.Anything).Return(nil)

	_, err := svc.Create(context.Background(), uuid.New(), req)
	require.ErrorIs(t, err, service.ErrInvalidID)
}

func TestShipmentServiceCreateNotFound(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	shipmentGarnitureRepo := &mocks.ShipmentGarnitureRepositoryMock{}
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewShipmentService(shipmentRepo, shipmentGarnitureRepo, garnitureRepo, txManager, testLogger(t))

	req := shipmentRequest()
	garnitureID := uuid.MustParse(req.Garnitures[0].GarnitureID)
	garnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, garnitureID).Return(model.Garniture{}, errors.New("not found"))
	txManager.On("RunInTx", mock.Anything).Return(nil)

	_, err := svc.Create(context.Background(), uuid.New(), req)
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestShipmentServiceCreateNotEnoughStock(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	shipmentGarnitureRepo := &mocks.ShipmentGarnitureRepositoryMock{}
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewShipmentService(shipmentRepo, shipmentGarnitureRepo, garnitureRepo, txManager, testLogger(t))

	req := shipmentRequest()
	garnitureID := uuid.MustParse(req.Garnitures[0].GarnitureID)
	garnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, garnitureID).Return(model.Garniture{ID: garnitureID, Price: 100, StockCount: 0}, nil)
	txManager.On("RunInTx", mock.Anything).Return(nil)

	_, err := svc.Create(context.Background(), uuid.New(), req)
	require.ErrorIs(t, err, service.ErrNotEnoughStock)
}

func TestShipmentServiceCreateDecreaseStockError(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	shipmentGarnitureRepo := &mocks.ShipmentGarnitureRepositoryMock{}
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewShipmentService(shipmentRepo, shipmentGarnitureRepo, garnitureRepo, txManager, testLogger(t))

	req := shipmentRequest()
	garnitureID := uuid.MustParse(req.Garnitures[0].GarnitureID)
	garnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, garnitureID).Return(model.Garniture{ID: garnitureID, Price: 100, StockCount: 5}, nil)
	garnitureRepo.On("DecreaseStockTx", mock.Anything, mock.Anything, garnitureID, req.Garnitures[0].Count).Return(errors.New("decrease error"))
	txManager.On("RunInTx", mock.Anything).Return(nil)

	_, err := svc.Create(context.Background(), uuid.New(), req)
	require.ErrorIs(t, err, service.ErrNotEnoughStock)
}

func TestShipmentServiceCreateSuccess(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	shipmentGarnitureRepo := &mocks.ShipmentGarnitureRepositoryMock{}
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewShipmentService(shipmentRepo, shipmentGarnitureRepo, garnitureRepo, txManager, testLogger(t))

	req := shipmentRequest()
	garnitureID := uuid.MustParse(req.Garnitures[0].GarnitureID)
	garnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, garnitureID).Return(model.Garniture{ID: garnitureID, Price: 100, StockCount: 5}, nil)
	garnitureRepo.On("DecreaseStockTx", mock.Anything, mock.Anything, garnitureID, req.Garnitures[0].Count).Return(nil)

	shipmentRepo.On("CreateTx", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	shipmentGarnitureRepo.On("CreateManyTx", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	txManager.On("RunInTx", mock.Anything).Return(nil)

	id, err := svc.Create(context.Background(), uuid.New(), req)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, id)
}

func TestShipmentServiceUpdateEmpty(t *testing.T) {
	svc := service.NewShipmentService(&mocks.ShipmentRepositoryMock{}, &mocks.ShipmentGarnitureRepositoryMock{}, &mocks.GarnitureRepositoryMock{}, newTxManager(), testLogger(t))

	err := svc.Update(context.Background(), uuid.New(), dto.ShipmentRequest{})
	require.ErrorIs(t, err, service.ErrEmptyComposition)
}

func TestShipmentServiceUpdateInvalidCount(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	shipmentGarnitureRepo := &mocks.ShipmentGarnitureRepositoryMock{}
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewShipmentService(shipmentRepo, shipmentGarnitureRepo, garnitureRepo, txManager, testLogger(t))

	req := shipmentRequest()
	req.Garnitures[0].Count = 0
	id := uuid.New()
	txManager.On("RunInTx", mock.Anything).Return(nil)
	shipmentRepo.On("GetByIDTx", mock.Anything, mock.Anything, id).Return(model.Shipment{ID: id, CustomerID: uuid.New()}, nil)
	shipmentGarnitureRepo.On("GetByShipmentIDTx", mock.Anything, mock.Anything, id).Return(nil, nil)

	err := svc.Update(context.Background(), id, req)
	require.ErrorIs(t, err, service.ErrInvalidCount)
}

func TestShipmentServiceUpdateInvalidID(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	shipmentGarnitureRepo := &mocks.ShipmentGarnitureRepositoryMock{}
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewShipmentService(shipmentRepo, shipmentGarnitureRepo, garnitureRepo, txManager, testLogger(t))

	req := shipmentRequest()
	req.Garnitures[0].GarnitureID = "bad"
	id := uuid.New()
	txManager.On("RunInTx", mock.Anything).Return(nil)
	shipmentRepo.On("GetByIDTx", mock.Anything, mock.Anything, id).Return(model.Shipment{ID: id, CustomerID: uuid.New()}, nil)
	shipmentGarnitureRepo.On("GetByShipmentIDTx", mock.Anything, mock.Anything, id).Return(nil, nil)

	err := svc.Update(context.Background(), id, req)
	require.ErrorIs(t, err, service.ErrInvalidID)
}

func TestShipmentServiceUpdateRelationsError(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	shipmentGarnitureRepo := &mocks.ShipmentGarnitureRepositoryMock{}
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewShipmentService(shipmentRepo, shipmentGarnitureRepo, garnitureRepo, txManager, testLogger(t))

	id := uuid.New()
	shipmentRepo.On("GetByIDTx", mock.Anything, mock.Anything, id).Return(model.Shipment{ID: id, CustomerID: uuid.New()}, nil)
	shipmentGarnitureRepo.On("GetByShipmentIDTx", mock.Anything, mock.Anything, id).Return(nil, errors.New("relations error"))
	txManager.On("RunInTx", mock.Anything).Return(nil)

	err := svc.Update(context.Background(), id, shipmentRequest())
	require.Error(t, err)
}

func TestShipmentServiceUpdateNotEnoughStock(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	shipmentGarnitureRepo := &mocks.ShipmentGarnitureRepositoryMock{}
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewShipmentService(shipmentRepo, shipmentGarnitureRepo, garnitureRepo, txManager, testLogger(t))

	id := uuid.New()
	req := shipmentRequest()
	req.Garnitures[0].Count = 2
	garnitureID := uuid.MustParse(req.Garnitures[0].GarnitureID)

	shipmentRepo.On("GetByIDTx", mock.Anything, mock.Anything, id).Return(model.Shipment{ID: id, CustomerID: uuid.New()}, nil)
	shipmentGarnitureRepo.On("GetByShipmentIDTx", mock.Anything, mock.Anything, id).Return([]relations.ShipmentGarniture{{ShipmentID: id, GarnitureID: garnitureID, Count: 1}}, nil)
	garnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, garnitureID).Return(model.Garniture{ID: garnitureID, Price: 100, StockCount: 2}, nil).Once()
	garnitureRepo.On("GetByIDTx", mock.Anything, mock.Anything, garnitureID).Return(model.Garniture{ID: garnitureID, Price: 100, StockCount: 0}, nil).Once()
	txManager.On("RunInTx", mock.Anything).Return(nil)

	err := svc.Update(context.Background(), id, req)
	require.ErrorIs(t, err, service.ErrNotEnoughStock)
}

func TestShipmentServiceUpdateInvalidDate(t *testing.T) {
	svc := service.NewShipmentService(&mocks.ShipmentRepositoryMock{}, &mocks.ShipmentGarnitureRepositoryMock{}, &mocks.GarnitureRepositoryMock{}, newTxManager(), testLogger(t))

	req := shipmentRequest()
	req.Date = "bad"

	err := svc.Update(context.Background(), uuid.New(), req)
	require.ErrorIs(t, err, service.ErrInvalidDate)
}

func TestShipmentServiceUpdateNotFound(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	shipmentGarnitureRepo := &mocks.ShipmentGarnitureRepositoryMock{}
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewShipmentService(shipmentRepo, shipmentGarnitureRepo, garnitureRepo, txManager, testLogger(t))

	id := uuid.New()
	shipmentRepo.On("GetByIDTx", mock.Anything, mock.Anything, id).Return(model.Shipment{}, errors.New("not found"))
	txManager.On("RunInTx", mock.Anything).Return(nil)

	err := svc.Update(context.Background(), id, shipmentRequest())
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestShipmentServiceDeleteSuccess(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	shipmentGarnitureRepo := &mocks.ShipmentGarnitureRepositoryMock{}
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewShipmentService(shipmentRepo, shipmentGarnitureRepo, garnitureRepo, txManager, testLogger(t))

	shipmentID := uuid.New()
	relationsItems := []relations.ShipmentGarniture{{ShipmentID: shipmentID, GarnitureID: uuid.New(), Count: 1}}

	shipmentGarnitureRepo.On("GetByShipmentIDTx", mock.Anything, mock.Anything, shipmentID).Return(relationsItems, nil)
	garnitureRepo.On("IncreaseStockTx", mock.Anything, mock.Anything, relationsItems[0].GarnitureID, relationsItems[0].Count).Return(nil)
	shipmentGarnitureRepo.On("DeleteByShipmentIDTx", mock.Anything, mock.Anything, shipmentID).Return(nil)
	shipmentRepo.On("DeleteTx", mock.Anything, mock.Anything, shipmentID).Return(nil)
	txManager.On("RunInTx", mock.Anything).Return(nil)

	err := svc.Delete(context.Background(), shipmentID)
	require.NoError(t, err)
}

func TestShipmentServiceDeleteRelationsError(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	shipmentGarnitureRepo := &mocks.ShipmentGarnitureRepositoryMock{}
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewShipmentService(shipmentRepo, shipmentGarnitureRepo, garnitureRepo, txManager, testLogger(t))

	id := uuid.New()
	txManager.On("RunInTx", mock.Anything).Return(nil)
	shipmentGarnitureRepo.On("GetByShipmentIDTx", mock.Anything, mock.Anything, id).Return(nil, errors.New("relations error"))

	err := svc.Delete(context.Background(), id)
	require.Error(t, err)
}

func TestShipmentServiceDeleteTxError(t *testing.T) {
	shipmentRepo := &mocks.ShipmentRepositoryMock{}
	shipmentGarnitureRepo := &mocks.ShipmentGarnitureRepositoryMock{}
	garnitureRepo := &mocks.GarnitureRepositoryMock{}
	txManager := newTxManager()

	svc := service.NewShipmentService(shipmentRepo, shipmentGarnitureRepo, garnitureRepo, txManager, testLogger(t))

	id := uuid.New()
	txManager.On("RunInTx", mock.Anything).Return(nil)
	shipmentGarnitureRepo.On("GetByShipmentIDTx", mock.Anything, mock.Anything, id).Return([]relations.ShipmentGarniture{{ShipmentID: id, GarnitureID: uuid.New(), Count: 1}}, nil)
	garnitureRepo.On("IncreaseStockTx", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	shipmentGarnitureRepo.On("DeleteByShipmentIDTx", mock.Anything, mock.Anything, id).Return(nil)
	shipmentRepo.On("DeleteTx", mock.Anything, mock.Anything, id).Return(errors.New("delete error"))

	err := svc.Delete(context.Background(), id)
	require.Error(t, err)
}
