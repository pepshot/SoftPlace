package servicetests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/pepshot/SoftPlace/services/customer-service/tests/mocks"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/service"
	"github.com/pepshot/SoftPlace/shared/config"
	"github.com/pepshot/SoftPlace/shared/logger"
)

func testLogger(t *testing.T) *logger.Logger {
	t.Helper()

	log, err := logger.New(config.LoggerConfig{
		Level:   "error",
		Format:  "text",
		Outputs: []string{"stdout"},
	}, "customer-service-test")
	require.NoError(t, err)
	return log
}

func newCustomerService(t *testing.T) (*service.CustomerService, *mocks.CustomerRepositoryMock, *mocks.JWTUseCaseMock) {
	t.Helper()

	repo := &mocks.CustomerRepositoryMock{}
	jwt := &mocks.JWTUseCaseMock{}
	svc := service.NewCustomerService(repo, jwt, testLogger(t))

	return svc, repo, jwt
}

func sampleRegisterRequest() dto.RegisterCustomerRequest {
	return dto.RegisterCustomerRequest{
		Login:           "login",
		Password:        "secret123",
		ConfirmPassword: "secret123",
		Email:           "user@example.com",
	}
}

func sampleLoginRequest(login, password string) dto.LoginCustomerRequest {
	return dto.LoginCustomerRequest{
		Login:    login,
		Password: password,
	}
}

func sampleCustomer(id uuid.UUID, password string) model.Customer {
	return model.Customer{
		ID:       id,
		Login:    "login",
		Password: password,
		Email:    "user@example.com",
	}
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	return string(hash)
}

func newTxManager() *mocks.TransactionManagerMock {
	return &mocks.TransactionManagerMock{}
}

func furnitureRequest() dto.FurnitureRequest {
	return dto.FurnitureRequest{
		Name: "Chair",
		Code: "CH-1",
		Modules: []dto.FurnitureModuleItem{{
			ModuleID: uuid.New().String(),
			Count:    1,
		}},
	}
}

func garnitureRequest() dto.GarnitureRequest {
	return dto.GarnitureRequest{
		Name: "Set",
		Code: "GS-1",
		Furniture: []dto.GarnitureFurnitureItem{{
			FurnitureID: uuid.New().String(),
			Count:       1,
		}},
	}
}

func shipmentRequest() dto.ShipmentRequest {
	return dto.ShipmentRequest{
		Code: "SH-1",
		Date: "2026-05-09",
		Garnitures: []dto.ShipmentGarnitureItem{{
			GarnitureID: uuid.New().String(),
			Count:       1,
		}},
	}
}

func moduleInfo(moduleID string, price float64, stock int) service.ModuleInfo {
	return service.ModuleInfo{
		ID:         moduleID,
		Name:       "Module",
		Code:       "M-1",
		Price:      price,
		StockCount: stock,
	}
}

func newRegisterRequest() dto.RegisterCustomerRequest {
	return sampleRegisterRequest()
}

func newLoginRequest(login, password string) dto.LoginCustomerRequest {
	return sampleLoginRequest(login, password)
}

func newCustomer(id uuid.UUID, password string) model.Customer {
	return sampleCustomer(id, password)
}
