package servicetests

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/service"
	"github.com/pepshot/SoftPlace/services/supplier-service/tests/mocks"
	"github.com/pepshot/SoftPlace/shared/config"
	"github.com/pepshot/SoftPlace/shared/logger"
)

func testLogger(t *testing.T) *logger.Logger {
	t.Helper()

	log, err := logger.New(config.LoggerConfig{
		Level:   "error",
		Format:  "text",
		Outputs: []string{"stdout"},
	}, "supplier-service-test")
	require.NoError(t, err)
	return log
}

func newSupplierService(t *testing.T) (*service.SupplierService, *mocks.SupplierRepositoryMock, *mocks.JWTUseCaseMock) {
	t.Helper()
	repo := &mocks.SupplierRepositoryMock{}
	jwt := &mocks.JWTUseCaseMock{}
	return service.NewSupplierService(repo, jwt, testLogger(t)), repo, jwt
}

func newTxManager() *mocks.TransactionManagerMock {
	return &mocks.TransactionManagerMock{}
}

func registerRequest() dto.RegisterSupplierRequest {
	return dto.RegisterSupplierRequest{
		Login:           "supplier-login",
		Password:        "secret123",
		ConfirmPassword: "secret123",
		Email:           "supplier@example.com",
	}
}

func loginRequest(login, password string) dto.LoginSupplierRequest {
	return dto.LoginSupplierRequest{Login: login, Password: password}
}

func supplierModel(id uuid.UUID, password string) models.Supplier {
	return models.Supplier{ID: id, Login: "supplier-login", Password: password, Email: "supplier@example.com"}
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	return string(hash)
}

func moduleRequest() dto.ModuleRequest {
	return dto.ModuleRequest{
		Name: "Module",
		Code: "MOD-1",
		Materials: []dto.ModuleMaterialItem{{
			MaterialID: uuid.New().String(),
			Count:      1,
		}},
	}
}

func materialRequest() dto.MaterialRequest {
	return dto.MaterialRequest{Name: "Material", Code: "MAT-1", Price: 10, StockCount: 5}
}

func supplyRequest() dto.SupplyRequest {
	return dto.SupplyRequest{
		Code: "SUP-1",
		Date: "2026-05-20",
		Materials: []dto.SupplyMaterialItem{{
			MaterialID: uuid.New().String(),
			Count:      1,
		}},
	}
}

func sampleModule(id uuid.UUID) models.Module {
	return models.Module{ID: id, Name: "Module", Code: "MOD-1", Price: 10, StockCount: 2}
}

func sampleMaterial(id uuid.UUID) models.Material {
	return models.Material{ID: id, Name: "Material", Code: "MAT-1", Price: 5, StockCount: 10}
}

func sampleSupply(id, supplierID uuid.UUID) models.Supply {
	return models.Supply{ID: id, SupplierID: supplierID, Code: "SUP-1", Date: time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC), Price: 25}
}
