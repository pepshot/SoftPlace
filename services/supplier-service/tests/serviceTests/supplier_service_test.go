package servicetests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/service"
)

func TestSupplierServiceRegister(t *testing.T) {
	svc, repo, jwt := newSupplierService(t)
	req := registerRequest()

	repo.On("ExistsByLogin", context.Background(), req.Login).Return(false, nil)
	repo.On("ExistsByEmail", context.Background(), req.Email).Return(false, nil)
	repo.On("Create", context.Background(), mock.MatchedBy(func(item models.Supplier) bool {
		return item.Login == req.Login && item.Email == req.Email && item.Password != ""
	})).Return(nil)
	jwt.On("GenerateToken", mock.AnythingOfType("string"), "supplier").Return("token", nil)

	result, err := svc.Register(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, "token", result.Token)
	require.Equal(t, req.Login, result.Supplier.Login)
	require.Equal(t, req.Email, result.Supplier.Email)
}

func TestSupplierServiceRegisterPasswordMismatch(t *testing.T) {
	svc, _, _ := newSupplierService(t)
	req := registerRequest()
	req.ConfirmPassword = "wrong"

	_, err := svc.Register(context.Background(), req)
	require.ErrorIs(t, err, service.ErrPasswordMismatch)
}

func TestSupplierServiceRegisterLoginAlreadyUsed(t *testing.T) {
	svc, repo, _ := newSupplierService(t)
	req := registerRequest()
	repo.On("ExistsByLogin", context.Background(), req.Login).Return(true, nil)

	_, err := svc.Register(context.Background(), req)
	require.ErrorIs(t, err, service.ErrLoginAlreadyUsed)
}

func TestSupplierServiceRegisterEmailAlreadyUsed(t *testing.T) {
	svc, repo, _ := newSupplierService(t)
	req := registerRequest()
	repo.On("ExistsByLogin", context.Background(), req.Login).Return(false, nil)
	repo.On("ExistsByEmail", context.Background(), req.Email).Return(true, nil)

	_, err := svc.Register(context.Background(), req)
	require.ErrorIs(t, err, service.ErrEmailAlreadyUsed)
}

func TestSupplierServiceRegisterRepoError(t *testing.T) {
	svc, repo, _ := newSupplierService(t)
	req := registerRequest()
	repo.On("ExistsByLogin", context.Background(), req.Login).Return(false, errors.New("db failure"))

	_, err := svc.Register(context.Background(), req)
	require.Error(t, err)
}

func TestSupplierServiceRegisterCreateError(t *testing.T) {
	svc, repo, _ := newSupplierService(t)
	req := registerRequest()
	repo.On("ExistsByLogin", context.Background(), req.Login).Return(false, nil)
	repo.On("ExistsByEmail", context.Background(), req.Email).Return(false, nil)
	repo.On("Create", context.Background(), mock.AnythingOfType("models.Supplier")).Return(errors.New("insert error"))

	_, err := svc.Register(context.Background(), req)
	require.Error(t, err)
}

func TestSupplierServiceRegisterTokenError(t *testing.T) {
	svc, repo, jwt := newSupplierService(t)
	req := registerRequest()
	repo.On("ExistsByLogin", context.Background(), req.Login).Return(false, nil)
	repo.On("ExistsByEmail", context.Background(), req.Email).Return(false, nil)
	repo.On("Create", context.Background(), mock.AnythingOfType("models.Supplier")).Return(nil)
	jwt.On("GenerateToken", mock.AnythingOfType("string"), "supplier").Return("", errors.New("token error"))

	_, err := svc.Register(context.Background(), req)
	require.Error(t, err)
}

func TestSupplierServiceLoginSuccess(t *testing.T) {
	svc, repo, jwt := newSupplierService(t)
	password := "secret123"
	supplier := supplierModel(uuid.New(), hashPassword(t, password))

	repo.On("GetByLogin", context.Background(), supplier.Login).Return(supplier, nil)
	jwt.On("GenerateToken", supplier.ID.String(), "supplier").Return("token", nil)

	result, err := svc.Login(context.Background(), loginRequest(supplier.Login, password))
	require.NoError(t, err)
	require.Equal(t, "token", result.Token)
	require.Equal(t, supplier.Login, result.Supplier.Login)
}

func TestSupplierServiceLoginInvalid(t *testing.T) {
	svc, repo, _ := newSupplierService(t)
	repo.On("GetByLogin", context.Background(), "supplier-login").Return(models.Supplier{}, errors.New("not found"))

	_, err := svc.Login(context.Background(), loginRequest("supplier-login", "secret123"))
	require.ErrorIs(t, err, service.ErrInvalidCredentials)
}

func TestSupplierServiceLoginInvalidPassword(t *testing.T) {
	svc, repo, _ := newSupplierService(t)
	supplier := supplierModel(uuid.New(), hashPassword(t, "secret123"))
	repo.On("GetByLogin", context.Background(), supplier.Login).Return(supplier, nil)

	_, err := svc.Login(context.Background(), loginRequest(supplier.Login, "wrong"))
	require.ErrorIs(t, err, service.ErrInvalidCredentials)
}

func TestSupplierServiceLoginTokenError(t *testing.T) {
	svc, repo, jwt := newSupplierService(t)
	supplier := supplierModel(uuid.New(), hashPassword(t, "secret123"))
	repo.On("GetByLogin", context.Background(), supplier.Login).Return(supplier, nil)
	jwt.On("GenerateToken", supplier.ID.String(), "supplier").Return("", errors.New("token error"))

	_, err := svc.Login(context.Background(), loginRequest(supplier.Login, "secret123"))
	require.Error(t, err)
}

func TestSupplierServiceGetProfile(t *testing.T) {
	svc, repo, _ := newSupplierService(t)
	supplierID := uuid.New()
	supplier := supplierModel(supplierID, "")
	repo.On("GetByID", context.Background(), supplierID).Return(supplier, nil)

	result, err := svc.GetProfile(context.Background(), supplierID)
	require.NoError(t, err)
	require.Equal(t, supplierID.String(), result.ID)
}

func TestSupplierServiceGetProfileNotFound(t *testing.T) {
	svc, repo, _ := newSupplierService(t)
	supplierID := uuid.New()
	repo.On("GetByID", context.Background(), supplierID).Return(models.Supplier{}, errors.New("not found"))

	_, err := svc.GetProfile(context.Background(), supplierID)
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestSupplierServiceGetProfileSuccess(t *testing.T) {
	svc, repo, _ := newSupplierService(t)
	supplierID := uuid.New()
	supplier := supplierModel(supplierID, "")
	repo.On("GetByID", context.Background(), supplierID).Return(supplier, nil)

	result, err := svc.GetProfile(context.Background(), supplierID)
	require.NoError(t, err)
	require.Equal(t, supplierID.String(), result.ID)
	require.Equal(t, supplier.Login, result.Login)
}
