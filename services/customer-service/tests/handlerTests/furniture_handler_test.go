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

func TestListFurnitureHandler(t *testing.T) {
	handler, _, furnitureMock, _, _ := newHandler(t)

	items := []dto.FurnitureResponse{
		{ID: uuid.New().String(), Name: "Chair", Code: "CH-1", Price: 10, StockCount: 2},
	}

	furnitureMock.On("GetList", mock.Anything).Return(items, nil)

	ctx, router, recorder := newTestContext(http.MethodGet, "/furniture", nil)

	router.GET("/furniture", handler.ListFurniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestListFurnitureHandlerError(t *testing.T) {
	handler, _, furnitureMock, _, _ := newHandler(t)

	furnitureMock.On("GetList", mock.Anything).Return(nil, service.ErrExternalService)

	ctx, router, recorder := newTestContext(http.MethodGet, "/furniture", nil)

	router.GET("/furniture", handler.ListFurniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadGateway, recorder.Code)
}

func TestGetFurnitureHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	ctx, router, recorder := newTestContext(http.MethodGet, "/furniture/invalid", nil)

	router.GET("/furniture/:id", handler.GetFurniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetFurnitureHandlerNotFound(t *testing.T) {
	handler, _, furnitureMock, _, _ := newHandler(t)

	itemID := uuid.New()
	furnitureMock.On("GetByID", mock.Anything, itemID).Return(dto.FurnitureResponse{}, service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodGet, "/furniture/"+itemID.String(), nil)

	router.GET("/furniture/:id", handler.GetFurniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestGetFurnitureHandlerSuccess(t *testing.T) {
	handler, _, furnitureMock, _, _ := newHandler(t)

	itemID := uuid.New()
	furnitureMock.On("GetByID", mock.Anything, itemID).Return(dto.FurnitureResponse{
		ID:         itemID.String(),
		Name:       "Chair",
		Code:       "CH-1",
		Price:      10,
		StockCount: 2,
	}, nil)

	ctx, router, recorder := newTestContext(http.MethodGet, "/furniture/"+itemID.String(), nil)

	router.GET("/furniture/:id", handler.GetFurniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestCreateFurnitureHandlerInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	ctx, router, recorder := newTestContext(http.MethodPost, "/furniture", []byte("invalid"))

	router.POST("/furniture", handler.CreateFurniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCreateFurnitureHandlerConflict(t *testing.T) {
	handler, _, furnitureMock, _, _ := newHandler(t)

	requestBody := furnitureRequest()

	furnitureMock.On("Create", mock.Anything, requestBody).Return(uuid.Nil, service.ErrNotEnoughStock)

	ctx, router, recorder := newTestContext(http.MethodPost, "/furniture", mustJSON(t, requestBody))

	router.POST("/furniture", handler.CreateFurniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestCreateFurnitureHandlerSuccess(t *testing.T) {
	handler, _, furnitureMock, _, _ := newHandler(t)

	requestBody := furnitureRequest()

	itemID := uuid.New()
	furnitureMock.On("Create", mock.Anything, requestBody).Return(itemID, nil)

	ctx, router, recorder := newTestContext(http.MethodPost, "/furniture", mustJSON(t, requestBody))

	router.POST("/furniture", handler.CreateFurniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestUpdateFurnitureHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	ctx, router, recorder := newTestContext(http.MethodPut, "/furniture/invalid", mustJSON(t, dto.FurnitureRequest{}))

	router.PUT("/furniture/:id", handler.UpdateFurniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestUpdateFurnitureHandlerInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	itemID := uuid.New()
	ctx, router, recorder := newTestContext(http.MethodPut, "/furniture/"+itemID.String(), []byte("invalid"))

	router.PUT("/furniture/:id", handler.UpdateFurniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestUpdateFurnitureHandlerNotFound(t *testing.T) {
	handler, _, furnitureMock, _, _ := newHandler(t)

	itemID := uuid.New()
	requestBody := furnitureRequest()

	furnitureMock.On("Update", mock.Anything, itemID, requestBody).Return(service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodPut, "/furniture/"+itemID.String(), mustJSON(t, requestBody))

	router.PUT("/furniture/:id", handler.UpdateFurniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestUpdateFurnitureHandlerSuccess(t *testing.T) {
	handler, _, furnitureMock, _, _ := newHandler(t)

	itemID := uuid.New()
	requestBody := furnitureRequest()

	furnitureMock.On("Update", mock.Anything, itemID, requestBody).Return(nil)

	ctx, router, recorder := newTestContext(http.MethodPut, "/furniture/"+itemID.String(), mustJSON(t, requestBody))

	router.PUT("/furniture/:id", handler.UpdateFurniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestDeleteFurnitureHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	ctx, router, recorder := newTestContext(http.MethodDelete, "/furniture/invalid", nil)

	router.DELETE("/furniture/:id", handler.DeleteFurniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestDeleteFurnitureHandlerNotFound(t *testing.T) {
	handler, _, furnitureMock, _, _ := newHandler(t)

	itemID := uuid.New()
	furnitureMock.On("Delete", mock.Anything, itemID).Return(service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodDelete, "/furniture/"+itemID.String(), nil)

	router.DELETE("/furniture/:id", handler.DeleteFurniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestDeleteFurnitureHandlerSuccess(t *testing.T) {
	handler, _, furnitureMock, _, _ := newHandler(t)

	itemID := uuid.New()
	furnitureMock.On("Delete", mock.Anything, itemID).Return(nil)

	ctx, router, recorder := newTestContext(http.MethodDelete, "/furniture/"+itemID.String(), nil)

	router.DELETE("/furniture/:id", handler.DeleteFurniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusNoContent, recorder.Code)
}
