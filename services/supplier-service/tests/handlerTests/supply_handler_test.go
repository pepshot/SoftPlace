package handlertests

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/service"
)

func TestListSuppliesHandler(t *testing.T) {
	handler, _, _, _, supplyMock := newHandler(t)
	items := []dto.SupplyResponse{{ID: uuid.New().String(), SupplierID: uuid.New().String(), Code: "SUP-1", Date: "2026-05-20", Price: 10}}
	supplyMock.On("GetList", mock.Anything).Return(items, nil)

	ctx, router, recorder := newTestContext(http.MethodGet, "/supplies", nil)
	router.GET("/supplies", handler.ListSupplies)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestListSuppliesHandlerError(t *testing.T) {
	handler, _, _, _, supplyMock := newHandler(t)
	supplyMock.On("GetList", mock.Anything).Return(nil, errors.New("db error"))

	ctx, router, recorder := newTestContext(http.MethodGet, "/supplies", nil)
	router.GET("/supplies", handler.ListSupplies)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestGetSupplyHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, router, recorder := newTestContext(http.MethodGet, "/supplies/bad", nil)
	router.GET("/supplies/:id", handler.GetSupply)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCreateSupplyHandlerUnauthorized(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, _, recorder := newTestContext(http.MethodPost, "/supplies", mustJSON(t, supplyRequest()))
	handler.CreateSupply(ctx)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestCreateSupplyHandlerInvalidUserIDType(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, _, recorder := newTestContext(http.MethodPost, "/supplies", mustJSON(t, supplyRequest()))
	setUserIDValue(ctx, 123)
	handler.CreateSupply(ctx)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestCreateSupplyHandlerInvalidUserIDValue(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, _, recorder := newTestContext(http.MethodPost, "/supplies", mustJSON(t, supplyRequest()))
	setUserID(ctx, "not-a-uuid")
	handler.CreateSupply(ctx)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestCreateSupplyHandlerInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, _, recorder := newTestContext(http.MethodPost, "/supplies", []byte("invalid"))
	setUserID(ctx, uuid.New().String())
	handler.CreateSupply(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCreateSupplyHandlerConflict(t *testing.T) {
	handler, _, _, _, supplyMock := newHandler(t)
	req := supplyRequest()
	supplierID := uuid.New()
	supplyMock.On("Create", mock.Anything, supplierID, req).Return(uuid.Nil, service.ErrNotEnoughStock)

	ctx, _, recorder := newTestContext(http.MethodPost, "/supplies", mustJSON(t, req))
	setUserID(ctx, supplierID.String())
	handler.CreateSupply(ctx)
	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestCreateSupplyHandlerSuccess(t *testing.T) {
	handler, _, _, _, supplyMock := newHandler(t)
	req := supplyRequest()
	supplierID := uuid.New()
	id := uuid.New()
	supplyMock.On("Create", mock.Anything, supplierID, req).Return(id, nil)

	ctx, _, recorder := newTestContext(http.MethodPost, "/supplies", mustJSON(t, req))
	setUserID(ctx, supplierID.String())
	handler.CreateSupply(ctx)
	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestGetSupplyHandlerNotFound(t *testing.T) {
	handler, _, _, _, supplyMock := newHandler(t)
	id := uuid.New()
	supplyMock.On("GetByID", mock.Anything, id).Return(dto.SupplyResponse{}, service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodGet, "/supplies/"+id.String(), nil)
	router.GET("/supplies/:id", handler.GetSupply)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestGetSupplyHandlerSuccess(t *testing.T) {
	handler, _, _, _, supplyMock := newHandler(t)
	id := uuid.New()
	supplyMock.On("GetByID", mock.Anything, id).Return(dto.SupplyResponse{ID: id.String(), SupplierID: uuid.New().String(), Code: "SUP-1", Date: "2026-05-20", Price: 10}, nil)

	ctx, router, recorder := newTestContext(http.MethodGet, "/supplies/"+id.String(), nil)
	router.GET("/supplies/:id", handler.GetSupply)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestUpdateSupplyHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, router, recorder := newTestContext(http.MethodPut, "/supplies/bad", mustJSON(t, dto.SupplyRequest{}))
	router.PUT("/supplies/:id", handler.UpdateSupply)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestUpdateSupplyHandlerInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	id := uuid.New()
	ctx, router, recorder := newTestContext(http.MethodPut, "/supplies/"+id.String(), []byte("invalid"))
	router.PUT("/supplies/:id", handler.UpdateSupply)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestUpdateSupplyHandlerNotFound(t *testing.T) {
	handler, _, _, _, supplyMock := newHandler(t)
	id := uuid.New()
	req := supplyRequest()
	supplyMock.On("Update", mock.Anything, id, req).Return(service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodPut, "/supplies/"+id.String(), mustJSON(t, req))
	router.PUT("/supplies/:id", handler.UpdateSupply)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestUpdateSupplyHandlerSuccess(t *testing.T) {
	handler, _, _, _, supplyMock := newHandler(t)
	id := uuid.New()
	req := supplyRequest()
	supplyMock.On("Update", mock.Anything, id, req).Return(nil)

	ctx, router, recorder := newTestContext(http.MethodPut, "/supplies/"+id.String(), mustJSON(t, req))
	router.PUT("/supplies/:id", handler.UpdateSupply)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestDeleteSupplyHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, router, recorder := newTestContext(http.MethodDelete, "/supplies/bad", nil)
	router.DELETE("/supplies/:id", handler.DeleteSupply)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestDeleteSupplyHandlerNotFound(t *testing.T) {
	handler, _, _, _, supplyMock := newHandler(t)
	id := uuid.New()
	supplyMock.On("Delete", mock.Anything, id).Return(service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodDelete, "/supplies/"+id.String(), nil)
	router.DELETE("/supplies/:id", handler.DeleteSupply)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestDeleteSupplyHandlerSuccess(t *testing.T) {
	handler, _, _, _, supplyMock := newHandler(t)
	id := uuid.New()
	supplyMock.On("Delete", mock.Anything, id).Return(nil)

	ctx, router, recorder := newTestContext(http.MethodDelete, "/supplies/"+id.String(), nil)
	router.DELETE("/supplies/:id", handler.DeleteSupply)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusNoContent, recorder.Code)
}
