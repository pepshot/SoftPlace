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

func TestRegisterHandlerSuccess(t *testing.T) {
	handler, supplierMock, _, _, _ := newHandler(t)
	requestBody := dto.RegisterSupplierRequest{Login: "supplier", Password: "secret123", ConfirmPassword: "secret123", Email: "supplier@example.com"}

	supplierMock.On("Register", mock.Anything, requestBody).Return(dto.AuthSupplierResponse{Token: "token", Supplier: dto.SupplierResponse{ID: uuid.New().String(), Login: requestBody.Login, Email: requestBody.Email}}, nil)

	ctx, router, recorder := newTestContext(http.MethodPost, "/auth/register", mustJSON(t, requestBody))
	router.POST("/auth/register", handler.Register)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestRegisterHandlerInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, router, recorder := newTestContext(http.MethodPost, "/auth/register", []byte("invalid"))
	router.POST("/auth/register", handler.Register)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestRegisterHandlerConflict(t *testing.T) {
	handler, supplierMock, _, _, _ := newHandler(t)
	requestBody := dto.RegisterSupplierRequest{Login: "supplier", Password: "secret123", ConfirmPassword: "secret123", Email: "supplier@example.com"}
	supplierMock.On("Register", mock.Anything, requestBody).Return(dto.AuthSupplierResponse{}, service.ErrLoginAlreadyUsed)

	ctx, router, recorder := newTestContext(http.MethodPost, "/auth/register", mustJSON(t, requestBody))
	router.POST("/auth/register", handler.Register)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestRegisterHandlerInternalError(t *testing.T) {
	handler, supplierMock, _, _, _ := newHandler(t)
	requestBody := dto.RegisterSupplierRequest{Login: "supplier", Password: "secret123", ConfirmPassword: "secret123", Email: "supplier@example.com"}
	supplierMock.On("Register", mock.Anything, requestBody).Return(dto.AuthSupplierResponse{}, errors.New("fail"))

	ctx, router, recorder := newTestContext(http.MethodPost, "/auth/register", mustJSON(t, requestBody))
	router.POST("/auth/register", handler.Register)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestLoginHandlerSuccess(t *testing.T) {
	handler, supplierMock, _, _, _ := newHandler(t)
	requestBody := dto.LoginSupplierRequest{Login: "supplier", Password: "secret123"}
	supplierMock.On("Login", mock.Anything, requestBody).Return(dto.AuthSupplierResponse{Token: "token", Supplier: dto.SupplierResponse{ID: uuid.New().String(), Login: requestBody.Login, Email: "supplier@example.com"}}, nil)

	ctx, router, recorder := newTestContext(http.MethodPost, "/auth/login", mustJSON(t, requestBody))
	router.POST("/auth/login", handler.Login)
	router.HandleContext(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestLoginHandlerInvalidJSON(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, router, recorder := newTestContext(http.MethodPost, "/auth/login", []byte("invalid"))
	router.POST("/auth/login", handler.Login)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestLoginHandlerUnauthorized(t *testing.T) {
	handler, supplierMock, _, _, _ := newHandler(t)
	requestBody := dto.LoginSupplierRequest{Login: "supplier", Password: "secret123"}
	supplierMock.On("Login", mock.Anything, requestBody).Return(dto.AuthSupplierResponse{}, service.ErrInvalidCredentials)

	ctx, router, recorder := newTestContext(http.MethodPost, "/auth/login", mustJSON(t, requestBody))
	router.POST("/auth/login", handler.Login)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestLoginHandlerInternalError(t *testing.T) {
	handler, supplierMock, _, _, _ := newHandler(t)
	requestBody := dto.LoginSupplierRequest{Login: "supplier", Password: "secret123"}
	supplierMock.On("Login", mock.Anything, requestBody).Return(dto.AuthSupplierResponse{}, errors.New("fail"))

	ctx, router, recorder := newTestContext(http.MethodPost, "/auth/login", mustJSON(t, requestBody))
	router.POST("/auth/login", handler.Login)
	router.HandleContext(ctx)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestProfileHandlerUnauthorized(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, _, recorder := newTestContext(http.MethodGet, "/profile", nil)
	handler.Profile(ctx)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestProfileHandlerInvalidUserIDType(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, _, recorder := newTestContext(http.MethodGet, "/profile", nil)
	setUserIDValue(ctx, 123)
	handler.Profile(ctx)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestProfileHandlerInvalidUserIDValue(t *testing.T) {
	handler, _, _, _, _ := newHandler(t)
	ctx, _, recorder := newTestContext(http.MethodGet, "/profile", nil)
	setUserID(ctx, "not-a-uuid")
	handler.Profile(ctx)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestProfileHandlerNotFound(t *testing.T) {
	handler, supplierMock, _, _, _ := newHandler(t)
	supplierID := uuid.New()
	supplierMock.On("GetProfile", mock.Anything, supplierID).Return(dto.SupplierResponse{}, service.ErrNotFound)

	ctx, _, recorder := newTestContext(http.MethodGet, "/profile", nil)
	setUserID(ctx, supplierID.String())
	handler.Profile(ctx)
	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestProfileHandlerSuccess(t *testing.T) {
	handler, supplierMock, _, _, _ := newHandler(t)
	supplierID := uuid.New()
	supplierMock.On("GetProfile", mock.Anything, supplierID).Return(dto.SupplierResponse{ID: supplierID.String(), Login: "supplier", Email: "supplier@example.com"}, nil)

	ctx, _, recorder := newTestContext(http.MethodGet, "/profile", nil)
	setUserID(ctx, supplierID.String())
	handler.Profile(ctx)
	require.Equal(t, http.StatusOK, recorder.Code)
}
