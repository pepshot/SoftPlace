package servicetests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/service"
)

func TestCustomerServiceRegister(t *testing.T) {
	svc, repo, jwt := newCustomerService(t)

	req := newRegisterRequest()

	repo.On("ExistsByLogin", context.Background(), req.Login).Return(false, nil)
	repo.On("ExistsByEmail", context.Background(), req.Email).Return(false, nil)
	repo.On("Create", context.Background(), mock.MatchedBy(func(value model.Customer) bool {
		return value.Login == req.Login && value.Email == req.Email && value.Password != ""
	})).Return(nil)
	jwt.On("GenerateToken", mock.AnythingOfType("string"), "customer").Return("token", nil)

	result, err := svc.Register(context.Background(), req)

	require.NoError(t, err)
	require.Equal(t, "token", result.Token)
	require.Equal(t, req.Login, result.Customer.Login)
	require.Equal(t, req.Email, result.Customer.Email)
}

func TestCustomerServiceRegisterPasswordMismatch(t *testing.T) {
	svc, _, _ := newCustomerService(t)

	req := newRegisterRequest()
	req.ConfirmPassword = "wrongpassword"

	_, err := svc.Register(context.Background(), req)

	require.ErrorIs(t, err, service.ErrPasswordMismatch)
}

func TestCustomerServiceRegisterLoginAlreadyUsed(t *testing.T) {
	svc, repo, _ := newCustomerService(t)

	req := newRegisterRequest()

	repo.On("ExistsByLogin", context.Background(), req.Login).Return(true, nil)

	_, err := svc.Register(context.Background(), req)

	require.ErrorIs(t, err, service.ErrLoginAlreadyUsed)
}

func TestCustomerServiceRegisterEmailAlreadyUsed(t *testing.T) {
	svc, repo, _ := newCustomerService(t)

	req := newRegisterRequest()

	repo.On("ExistsByLogin", context.Background(), req.Login).Return(false, nil)
	repo.On("ExistsByEmail", context.Background(), req.Email).Return(true, nil)

	_, err := svc.Register(context.Background(), req)

	require.ErrorIs(t, err, service.ErrEmailAlreadyUsed)
}

func TestCustomerServiceRegisterRepoError(t *testing.T) {
	svc, repo, _ := newCustomerService(t)

	req := newRegisterRequest()

	repo.On("ExistsByLogin", context.Background(), req.Login).Return(false, errors.New("db failure"))

	_, err := svc.Register(context.Background(), req)

	require.Error(t, err)
}

func TestCustomerServiceRegisterCreateError(t *testing.T) {
	svc, repo, _ := newCustomerService(t)

	req := newRegisterRequest()

	repo.On("ExistsByLogin", context.Background(), req.Login).Return(false, nil)
	repo.On("ExistsByEmail", context.Background(), req.Email).Return(false, nil)
	repo.On("Create", context.Background(), mock.AnythingOfType("model.Customer")).Return(errors.New("insert error"))

	_, err := svc.Register(context.Background(), req)

	require.Error(t, err)
}

func TestCustomerServiceRegisterTokenError(t *testing.T) {
	svc, repo, jwt := newCustomerService(t)

	req := newRegisterRequest()

	repo.On("ExistsByLogin", context.Background(), req.Login).Return(false, nil)
	repo.On("ExistsByEmail", context.Background(), req.Email).Return(false, nil)
	repo.On("Create", context.Background(), mock.AnythingOfType("model.Customer")).Return(nil)
	jwt.On("GenerateToken", mock.AnythingOfType("string"), "customer").Return("", errors.New("token error"))

	_, err := svc.Register(context.Background(), req)

	require.Error(t, err)
}

func TestCustomerServiceLogin(t *testing.T) {
	svc, repo, jwt := newCustomerService(t)

	password := "secret123"
	customer := newCustomer(uuid.New(), hashPassword(t, password))

	repo.On("GetByLogin", context.Background(), customer.Login).Return(customer, nil)
	jwt.On("GenerateToken", customer.ID.String(), "customer").Return("token", nil)

	result, err := svc.Login(context.Background(), newLoginRequest(customer.Login, password))

	require.NoError(t, err)
	require.Equal(t, "token", result.Token)
	require.Equal(t, customer.Login, result.Customer.Login)
}

func TestCustomerServiceLoginInvalid(t *testing.T) {
	svc, repo, _ := newCustomerService(t)

	repo.On("GetByLogin", context.Background(), "login").Return(model.Customer{}, errors.New("not found"))

	_, err := svc.Login(context.Background(), newLoginRequest("login", "secret123"))

	require.ErrorIs(t, err, service.ErrInvalidCredentials)
}

func TestCustomerServiceLoginInvalidPassword(t *testing.T) {
	svc, repo, _ := newCustomerService(t)

	customer := newCustomer(uuid.New(), hashPassword(t, "secret123"))

	repo.On("GetByLogin", context.Background(), customer.Login).Return(customer, nil)

	_, err := svc.Login(context.Background(), newLoginRequest(customer.Login, "wrong"))

	require.ErrorIs(t, err, service.ErrInvalidCredentials)
}

func TestCustomerServiceLoginTokenError(t *testing.T) {
	svc, repo, jwt := newCustomerService(t)

	customer := newCustomer(uuid.New(), hashPassword(t, "secret123"))

	repo.On("GetByLogin", context.Background(), customer.Login).Return(customer, nil)
	jwt.On("GenerateToken", customer.ID.String(), "customer").Return("", errors.New("token error"))

	_, err := svc.Login(context.Background(), newLoginRequest(customer.Login, "secret123"))

	require.Error(t, err)
}

func TestCustomerServiceProfileNotFound(t *testing.T) {
	svc, repo, _ := newCustomerService(t)

	customerID := uuid.New()
	repo.On("GetByID", context.Background(), customerID).Return(model.Customer{}, errors.New("not found"))

	_, err := svc.GetProfile(context.Background(), customerID)

	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestCustomerServiceProfileSuccess(t *testing.T) {
	svc, repo, _ := newCustomerService(t)

	customerID := uuid.New()
	customer := newCustomer(customerID, "")

	repo.On("GetByID", context.Background(), customerID).Return(customer, nil)

	result, err := svc.GetProfile(context.Background(), customerID)

	require.NoError(t, err)
	require.Equal(t, customer.ID.String(), result.ID)
	require.Equal(t, customer.Login, result.Login)
}
