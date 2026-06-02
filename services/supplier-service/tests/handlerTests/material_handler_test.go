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

func TestListMaterialsHandler(t *testing.T) {
	handler, _, _, materialMock, _ := newHandler(t)
	items := []dto.MaterialResponse{{ID: uuid.New().String(), Name: "Material", Code: "MAT-1", Price: 10, StockCount: 2}}
	materialMock.On("GetList", mock.Anything).Return(items, nil)

	ctx, router, recorder := newTestContext(http.MethodGet, "/materials", nil)
	router.GET("/materials", handler.ListMaterials)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestListMaterialsHandlerError(t *testing.T) {
	handler, _, _, materialMock, _ := newHandler(t)
	materialMock.On("GetList", mock.Anything).Return(nil, errors.New("db error"))

	ctx, router, recorder := newTestContext(http.MethodGet, "/materials", nil)
	router.GET("/materials", handler.ListMaterials)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestGetMaterialHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, router, recorder := newTestContext(http.MethodGet, "/materials/bad", nil)
	router.GET("/materials/:id", handler.GetMaterial)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCreateMaterialHandlerInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, router, recorder := newTestContext(http.MethodPost, "/materials", []byte("invalid"))
	router.POST("/materials", handler.CreateMaterial)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCreateMaterialHandlerConflict(t *testing.T) {
	handler, _, _, materialMock, _ := newHandler(t)
	req := materialRequest()
	materialMock.On("Create", mock.Anything, req).Return(uuid.Nil, service.ErrCodeAlreadyUsed)

	ctx, router, recorder := newTestContext(http.MethodPost, "/materials", mustJSON(t, req))
	router.POST("/materials", handler.CreateMaterial)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestCreateMaterialHandlerSuccess(t *testing.T) {
	handler, _, _, materialMock, _ := newHandler(t)
	req := materialRequest()
	id := uuid.New()
	materialMock.On("Create", mock.Anything, req).Return(id, nil)

	ctx, router, recorder := newTestContext(http.MethodPost, "/materials", mustJSON(t, req))
	router.POST("/materials", handler.CreateMaterial)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestGetMaterialHandlerNotFound(t *testing.T) {
	handler, _, _, materialMock, _ := newHandler(t)
	id := uuid.New()
	materialMock.On("GetByID", mock.Anything, id).Return(dto.MaterialResponse{}, service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodGet, "/materials/"+id.String(), nil)
	router.GET("/materials/:id", handler.GetMaterial)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestGetMaterialHandlerSuccess(t *testing.T) {
	handler, _, _, materialMock, _ := newHandler(t)
	id := uuid.New()
	materialMock.On("GetByID", mock.Anything, id).Return(dto.MaterialResponse{ID: id.String(), Name: "Material", Code: "MAT-1", Price: 10, StockCount: 2}, nil)

	ctx, router, recorder := newTestContext(http.MethodGet, "/materials/"+id.String(), nil)
	router.GET("/materials/:id", handler.GetMaterial)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestUpdateMaterialHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, router, recorder := newTestContext(http.MethodPut, "/materials/bad", mustJSON(t, dto.MaterialRequest{}))
	router.PUT("/materials/:id", handler.UpdateMaterial)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestUpdateMaterialHandlerInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	id := uuid.New()
	ctx, router, recorder := newTestContext(http.MethodPut, "/materials/"+id.String(), []byte("invalid"))
	router.PUT("/materials/:id", handler.UpdateMaterial)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestUpdateMaterialHandlerNotFound(t *testing.T) {
	handler, _, _, materialMock, _ := newHandler(t)
	id := uuid.New()
	req := materialRequest()
	materialMock.On("Update", mock.Anything, id, req).Return(service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodPut, "/materials/"+id.String(), mustJSON(t, req))
	router.PUT("/materials/:id", handler.UpdateMaterial)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestUpdateMaterialHandlerSuccess(t *testing.T) {
	handler, _, _, materialMock, _ := newHandler(t)
	id := uuid.New()
	req := materialRequest()
	materialMock.On("Update", mock.Anything, id, req).Return(nil)

	ctx, router, recorder := newTestContext(http.MethodPut, "/materials/"+id.String(), mustJSON(t, req))
	router.PUT("/materials/:id", handler.UpdateMaterial)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestDeleteMaterialHandlerNotFound(t *testing.T) {
	handler, _, _, materialMock, _ := newHandler(t)
	id := uuid.New()
	materialMock.On("Delete", mock.Anything, id).Return(service.ErrNotFound)

	ctx, router, recorder := newTestContext(http.MethodDelete, "/materials/"+id.String(), nil)
	router.DELETE("/materials/:id", handler.DeleteMaterial)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestDeleteMaterialHandlerSuccess(t *testing.T) {
	handler, _, _, materialMock, _ := newHandler(t)
	id := uuid.New()
	materialMock.On("Delete", mock.Anything, id).Return(nil)

	ctx, router, recorder := newTestContext(http.MethodDelete, "/materials/"+id.String(), nil)
	router.DELETE("/materials/:id", handler.DeleteMaterial)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestDeleteMaterialHandlerInvalidID(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, router, recorder := newTestContext(http.MethodDelete, "/materials/bad", nil)
	router.DELETE("/materials/:id", handler.DeleteMaterial)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}
