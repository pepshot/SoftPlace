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

func TestListModulesHandler(t *testing.T) {
	handler, _, moduleMock, _, _ := newHandler(t)
	items := []dto.ModuleResponse{{ID: uuid.New().String(), Name: "Module", Code: "MOD-1", Price: 10, StockCount: 2}}
	moduleMock.On("GetList", mock.Anything).Return(items, nil)

	ctx, router, recorder := newTestContext(http.MethodGet, "/modules", nil)
	router.GET("/modules", handler.ListModules)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestListModulesHandlerError(t *testing.T) {
	handler, _, moduleMock, _, _ := newHandler(t)
	moduleMock.On("GetList", mock.Anything).Return(nil, errors.New("db error"))

	ctx, router, recorder := newTestContext(http.MethodGet, "/modules", nil)
	router.GET("/modules", handler.ListModules)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestGetModuleHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, router, recorder := newTestContext(http.MethodGet, "/modules/bad", nil)
	router.GET("/modules/:id", handler.GetModule)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetModuleHandlerNotFound(t *testing.T) {
	handler, _, moduleMock, _, _ := newHandler(t)
	id := uuid.New()
	moduleMock.On("GetByID", mock.Anything, id).Return(dto.ModuleResponse{}, service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodGet, "/modules/"+id.String(), nil)
	router.GET("/modules/:id", handler.GetModule)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestGetModuleHandlerSuccess(t *testing.T) {
	handler, _, moduleMock, _, _ := newHandler(t)
	id := uuid.New()
	moduleMock.On("GetByID", mock.Anything, id).Return(dto.ModuleResponse{ID: id.String(), Name: "Module", Code: "MOD-1", Price: 10, StockCount: 2}, nil)

	ctx, router, recorder := newTestContext(http.MethodGet, "/modules/"+id.String(), nil)
	router.GET("/modules/:id", handler.GetModule)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestCreateModuleHandlerInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, router, recorder := newTestContext(http.MethodPost, "/modules", []byte("invalid"))
	router.POST("/modules", handler.CreateModule)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCreateModuleHandlerConflict(t *testing.T) {
	handler, _, moduleMock, _, _ := newHandler(t)
	req := moduleRequest()
	moduleMock.On("Create", mock.Anything, req).Return(uuid.Nil, service.ErrNotEnoughStock)

	ctx, router, recorder := newTestContext(http.MethodPost, "/modules", mustJSON(t, req))
	router.POST("/modules", handler.CreateModule)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestCreateModuleHandlerSuccess(t *testing.T) {
	handler, _, moduleMock, _, _ := newHandler(t)
	req := moduleRequest()
	id := uuid.New()
	moduleMock.On("Create", mock.Anything, req).Return(id, nil)

	ctx, router, recorder := newTestContext(http.MethodPost, "/modules", mustJSON(t, req))
	router.POST("/modules", handler.CreateModule)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestUpdateModuleHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, router, recorder := newTestContext(http.MethodPut, "/modules/bad", mustJSON(t, dto.ModuleRequest{}))
	router.PUT("/modules/:id", handler.UpdateModule)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestUpdateModuleHandlerInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	id := uuid.New()
	ctx, router, recorder := newTestContext(http.MethodPut, "/modules/"+id.String(), []byte("invalid"))
	router.PUT("/modules/:id", handler.UpdateModule)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestUpdateModuleHandlerNotFound(t *testing.T) {
	handler, _, moduleMock, _, _ := newHandler(t)
	id := uuid.New()
	req := moduleRequest()
	moduleMock.On("Update", mock.Anything, id, req).Return(service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodPut, "/modules/"+id.String(), mustJSON(t, req))
	router.PUT("/modules/:id", handler.UpdateModule)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestUpdateModuleHandlerSuccess(t *testing.T) {
	handler, _, moduleMock, _, _ := newHandler(t)
	id := uuid.New()
	req := moduleRequest()
	moduleMock.On("Update", mock.Anything, id, req).Return(nil)

	ctx, router, recorder := newTestContext(http.MethodPut, "/modules/"+id.String(), mustJSON(t, req))
	router.PUT("/modules/:id", handler.UpdateModule)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestDeleteModuleHandlerNotFound(t *testing.T) {
	handler, _, moduleMock, _, _ := newHandler(t)
	id := uuid.New()
	moduleMock.On("Delete", mock.Anything, id).Return(service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodDelete, "/modules/"+id.String(), nil)
	router.DELETE("/modules/:id", handler.DeleteModule)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestDeleteModuleHandlerSuccess(t *testing.T) {
	handler, _, moduleMock, _, _ := newHandler(t)
	id := uuid.New()
	moduleMock.On("Delete", mock.Anything, id).Return(nil)

	ctx, router, recorder := newTestContext(http.MethodDelete, "/modules/"+id.String(), nil)
	router.DELETE("/modules/:id", handler.DeleteModule)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusNoContent, recorder.Code)
}
