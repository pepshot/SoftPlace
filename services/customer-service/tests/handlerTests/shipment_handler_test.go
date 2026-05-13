package handlertests

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/service"
)

func TestListShipmentHandler(t *testing.T) {
	handler, _, _, _, shipmentMock := newHandler(t)

	items := []dto.ShipmentResponse{
		{ID: uuid.New().String(), Code: "SH-1", Date: "2026-05-09", Price: 200},
	}

	shipmentMock.On("GetList", mock.Anything).Return(items, nil)

	ctx, router, recorder := newTestContext(http.MethodGet, "/shipments", nil)

	router.GET("/shipments", handler.ListShipments)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestListShipmentHandlerError(t *testing.T) {
	handler, _, _, _, shipmentMock := newHandler(t)

	shipmentMock.On("GetList", mock.Anything).Return(nil, service.ErrExternalService)

	ctx, router, recorder := newTestContext(http.MethodGet, "/shipments", nil)

	router.GET("/shipments", handler.ListShipments)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadGateway, recorder.Code)
}

func TestGetShipmentHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	ctx, router, recorder := newTestContext(http.MethodGet, "/shipments/invalid", nil)

	router.GET("/shipments/:id", handler.GetShipment)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetShipmentHandlerNotFound(t *testing.T) {
	handler, _, _, _, shipmentMock := newHandler(t)

	itemID := uuid.New()
	shipmentMock.On("GetByID", mock.Anything, itemID).Return(dto.ShipmentResponse{}, service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodGet, "/shipments/"+itemID.String(), nil)

	router.GET("/shipments/:id", handler.GetShipment)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestGetShipmentHandlerSuccess(t *testing.T) {
	handler, _, _, _, shipmentMock := newHandler(t)

	itemID := uuid.New()
	shipmentMock.On("GetByID", mock.Anything, itemID).Return(dto.ShipmentResponse{
		ID:    itemID.String(),
		Code:  "SH-1",
		Date:  "2026-05-09",
		Price: 200,
	}, nil)

	ctx, router, recorder := newTestContext(http.MethodGet, "/shipments/"+itemID.String(), nil)

	router.GET("/shipments/:id", handler.GetShipment)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestCreateShipmentHandlerUnauthorized(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	ctx, _, recorder := newTestContext(http.MethodPost, "/shipments", mustJSON(t, shipmentRequest()))

	handler.CreateShipment(ctx)

	require.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestCreateShipmentHandlerInvalidUserIDType(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	ctx, _, recorder := newTestContext(http.MethodPost, "/shipments", mustJSON(t, shipmentRequest()))
	setUserIDValue(ctx, 123)

	handler.CreateShipment(ctx)

	require.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestCreateShipmentHandlerInvalidUserIDValue(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	ctx, _, recorder := newTestContext(http.MethodPost, "/shipments", mustJSON(t, shipmentRequest()))
	setUserID(ctx, "not-a-uuid")

	handler.CreateShipment(ctx)

	require.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestCreateShipmentHandlerInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	ctx, _, recorder := newTestContext(http.MethodPost, "/shipments", []byte("invalid"))
	setUserID(ctx, uuid.New().String())

	handler.CreateShipment(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCreateShipmentHandlerConflict(t *testing.T) {
	handler, _, _, _, shipmentMock := newHandler(t)

	requestBody := shipmentRequest()
	customerID := uuid.New()

	shipmentMock.On("Create", mock.Anything, customerID, requestBody).Return(uuid.Nil, service.ErrNotEnoughStock)

	ctx, _, recorder := newTestContext(http.MethodPost, "/shipments", mustJSON(t, requestBody))
	setUserID(ctx, customerID.String())

	handler.CreateShipment(ctx)

	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestCreateShipmentHandlerSuccess(t *testing.T) {
	handler, _, _, _, shipmentMock := newHandler(t)

	requestBody := shipmentRequest()
	customerID := uuid.New()
	shipmentID := uuid.New()

	shipmentMock.On("Create", mock.Anything, customerID, requestBody).Return(shipmentID, nil)

	ctx, _, recorder := newTestContext(http.MethodPost, "/shipments", mustJSON(t, requestBody))
	setUserID(ctx, customerID.String())

	handler.CreateShipment(ctx)

	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestUpdateShipmentHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	ctx, router, recorder := newTestContext(http.MethodPut, "/shipments/invalid", mustJSON(t, dto.ShipmentRequest{}))

	router.PUT("/shipments/:id", handler.UpdateShipment)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestUpdateShipmentHandlerInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	itemID := uuid.New()
	ctx, router, recorder := newTestContext(http.MethodPut, "/shipments/"+itemID.String(), []byte("invalid"))

	router.PUT("/shipments/:id", handler.UpdateShipment)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestUpdateShipmentHandlerNotFound(t *testing.T) {
	handler, _, _, _, shipmentMock := newHandler(t)

	itemID := uuid.New()
	requestBody := shipmentRequest()

	shipmentMock.On("Update", mock.Anything, itemID, requestBody).Return(service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodPut, "/shipments/"+itemID.String(), mustJSON(t, requestBody))

	router.PUT("/shipments/:id", handler.UpdateShipment)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestUpdateShipmentHandlerSuccess(t *testing.T) {
	handler, _, _, _, shipmentMock := newHandler(t)

	itemID := uuid.New()
	requestBody := shipmentRequest()

	shipmentMock.On("Update", mock.Anything, itemID, requestBody).Return(nil)

	ctx, router, recorder := newTestContext(http.MethodPut, "/shipments/"+itemID.String(), mustJSON(t, requestBody))

	router.PUT("/shipments/:id", handler.UpdateShipment)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestDeleteShipmentHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	ctx, router, recorder := newTestContext(http.MethodDelete, "/shipments/invalid", nil)

	router.DELETE("/shipments/:id", handler.DeleteShipment)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestDeleteShipmentHandlerNotFound(t *testing.T) {
	handler, _, _, _, shipmentMock := newHandler(t)

	itemID := uuid.New()
	shipmentMock.On("Delete", mock.Anything, itemID).Return(service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodDelete, "/shipments/"+itemID.String(), nil)

	router.DELETE("/shipments/:id", handler.DeleteShipment)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestDeleteShipmentHandlerSuccess(t *testing.T) {
	handler, _, _, _, shipmentMock := newHandler(t)

	itemID := uuid.New()
	shipmentMock.On("Delete", mock.Anything, itemID).Return(nil)

	ctx, router, recorder := newTestContext(http.MethodDelete, "/shipments/"+itemID.String(), nil)

	router.DELETE("/shipments/:id", handler.DeleteShipment)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusNoContent, recorder.Code)
}
