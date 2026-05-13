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

func TestListGarnitureHandler(t *testing.T) {
	handler, _, _, garnitureMock, _ := newHandler(t)

	items := []dto.GarnitureResponse{
		{ID: uuid.New().String(), Name: "Set", Code: "GS-1", Price: 100, StockCount: 1},
	}

	garnitureMock.On("GetList", mock.Anything).Return(items, nil)

	ctx, router, recorder := newTestContext(http.MethodGet, "/garnitures", nil)

	router.GET("/garnitures", handler.ListGarnitures)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestListGarnitureHandlerError(t *testing.T) {
	handler, _, _, garnitureMock, _ := newHandler(t)

	garnitureMock.On("GetList", mock.Anything).Return(nil, service.ErrExternalService)

	ctx, router, recorder := newTestContext(http.MethodGet, "/garnitures", nil)

	router.GET("/garnitures", handler.ListGarnitures)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadGateway, recorder.Code)
}

func TestGetGarnitureHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	ctx, router, recorder := newTestContext(http.MethodGet, "/garnitures/invalid", nil)

	router.GET("/garnitures/:id", handler.GetGarniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetGarnitureHandlerNotFound(t *testing.T) {
	handler, _, _, garnitureMock, _ := newHandler(t)

	itemID := uuid.New()
	garnitureMock.On("GetByID", mock.Anything, itemID).Return(dto.GarnitureResponse{}, service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodGet, "/garnitures/"+itemID.String(), nil)

	router.GET("/garnitures/:id", handler.GetGarniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestGetGarnitureHandlerSuccess(t *testing.T) {
	handler, _, _, garnitureMock, _ := newHandler(t)

	itemID := uuid.New()
	garnitureMock.On("GetByID", mock.Anything, itemID).Return(dto.GarnitureResponse{
		ID:         itemID.String(),
		Name:       "Set",
		Code:       "GS-1",
		Price:      100,
		StockCount: 1,
	}, nil)

	ctx, router, recorder := newTestContext(http.MethodGet, "/garnitures/"+itemID.String(), nil)

	router.GET("/garnitures/:id", handler.GetGarniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestCreateGarnitureHandlerInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	ctx, router, recorder := newTestContext(http.MethodPost, "/garnitures", []byte("invalid"))

	router.POST("/garnitures", handler.CreateGarniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCreateGarnitureHandlerConflict(t *testing.T) {
	handler, _, _, garnitureMock, _ := newHandler(t)

	requestBody := garnitureRequest()

	garnitureMock.On("Create", mock.Anything, requestBody).Return(uuid.Nil, service.ErrNotEnoughStock)

	ctx, router, recorder := newTestContext(http.MethodPost, "/garnitures", mustJSON(t, requestBody))

	router.POST("/garnitures", handler.CreateGarniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestCreateGarnitureHandlerSuccess(t *testing.T) {
	handler, _, _, garnitureMock, _ := newHandler(t)

	requestBody := garnitureRequest()

	itemID := uuid.New()
	garnitureMock.On("Create", mock.Anything, requestBody).Return(itemID, nil)

	ctx, router, recorder := newTestContext(http.MethodPost, "/garnitures", mustJSON(t, requestBody))

	router.POST("/garnitures", handler.CreateGarniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestUpdateGarnitureHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	ctx, router, recorder := newTestContext(http.MethodPut, "/garnitures/invalid", mustJSON(t, dto.GarnitureRequest{}))

	router.PUT("/garnitures/:id", handler.UpdateGarniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestUpdateGarnitureHandlerInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	itemID := uuid.New()
	ctx, router, recorder := newTestContext(http.MethodPut, "/garnitures/"+itemID.String(), []byte("invalid"))

	router.PUT("/garnitures/:id", handler.UpdateGarniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestUpdateGarnitureHandlerNotFound(t *testing.T) {
	handler, _, _, garnitureMock, _ := newHandler(t)

	itemID := uuid.New()
	requestBody := garnitureRequest()

	garnitureMock.On("Update", mock.Anything, itemID, requestBody).Return(service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodPut, "/garnitures/"+itemID.String(), mustJSON(t, requestBody))

	router.PUT("/garnitures/:id", handler.UpdateGarniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestUpdateGarnitureHandlerSuccess(t *testing.T) {
	handler, _, _, garnitureMock, _ := newHandler(t)

	itemID := uuid.New()
	requestBody := garnitureRequest()

	garnitureMock.On("Update", mock.Anything, itemID, requestBody).Return(nil)

	ctx, router, recorder := newTestContext(http.MethodPut, "/garnitures/"+itemID.String(), mustJSON(t, requestBody))

	router.PUT("/garnitures/:id", handler.UpdateGarniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestDeleteGarnitureHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)

	ctx, router, recorder := newTestContext(http.MethodDelete, "/garnitures/invalid", nil)

	router.DELETE("/garnitures/:id", handler.DeleteGarniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestDeleteGarnitureHandlerNotFound(t *testing.T) {
	handler, _, _, garnitureMock, _ := newHandler(t)

	itemID := uuid.New()
	garnitureMock.On("Delete", mock.Anything, itemID).Return(service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodDelete, "/garnitures/"+itemID.String(), nil)

	router.DELETE("/garnitures/:id", handler.DeleteGarniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestDeleteGarnitureHandlerSuccess(t *testing.T) {
	handler, _, _, garnitureMock, _ := newHandler(t)

	itemID := uuid.New()
	garnitureMock.On("Delete", mock.Anything, itemID).Return(nil)

	ctx, router, recorder := newTestContext(http.MethodDelete, "/garnitures/"+itemID.String(), nil)

	router.DELETE("/garnitures/:id", handler.DeleteGarniture)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusNoContent, recorder.Code)
}
