package repositorytests

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/models"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/repository"
	"github.com/pepshot/SoftPlace/shared/config"
	"github.com/pepshot/SoftPlace/shared/logger"
)

const (
	sqlSelectSupplierByID    = "SELECT id, login, password, email FROM suppliers WHERE id = $1;"
	sqlSelectSupplierByLogin = "SELECT id, login, password, email FROM suppliers WHERE login = $1;"
	sqlExistsSupplierLogin   = "SELECT EXISTS(SELECT id FROM suppliers WHERE login = $1);"
	sqlExistsSupplierEmail   = "SELECT EXISTS(SELECT id FROM suppliers WHERE email = $1);"
	sqlInsertSupplier        = "INSERT INTO suppliers"

	sqlSelectModuleByID   = "SELECT id, name, code, price, stock_count FROM modules WHERE id = $1;"
	sqlSelectModuleByCode = "SELECT id, name, code, price, stock_count FROM modules WHERE code = $1;"
	sqlSelectModules      = "SELECT id, name, code, price, stock_count FROM modules ORDER BY name;"
	sqlExistsModuleCode   = "SELECT EXISTS(SELECT id FROM modules WHERE code = $1);"
	sqlInsertModule       = "INSERT INTO modules"
	sqlUpdateModule       = "UPDATE modules"
	sqlDeleteModule       = "DELETE FROM modules"
	sqlIncreaseModule     = "UPDATE modules SET stock_count = stock_count + $2 WHERE id = $1;"
	sqlDecreaseModule     = "UPDATE modules SET stock_count = stock_count - $2 WHERE id = $1 AND stock_count >= $2;"

	sqlSelectMaterialByID = "SELECT id, name, code, price, stock_count FROM materials WHERE id = $1;"
	sqlSelectMaterials    = "SELECT id, name, code, price, stock_count FROM materials ORDER BY name;"
	sqlExistsMaterialCode = "SELECT EXISTS(SELECT id FROM materials WHERE code = $1);"
	sqlInsertMaterial     = "INSERT INTO materials"
	sqlUpdateMaterial     = "UPDATE materials"
	sqlDeleteMaterial     = "DELETE FROM materials"
	sqlIncreaseMaterial   = "UPDATE materials SET stock_count = stock_count + $2 WHERE id = $1;"
	sqlDecreaseMaterial   = "UPDATE materials SET stock_count = stock_count - $2 WHERE id = $1 AND stock_count >= $2;"

	sqlSelectSupplyByID = "SELECT id, supplier_id, code, date, price FROM supplies WHERE id = $1;"
	sqlSelectSupplies   = "SELECT id, supplier_id, code, date, price FROM supplies ORDER BY date DESC;"
	sqlExistsSupplyCode = "SELECT EXISTS(SELECT id FROM supplies WHERE code = $1);"
	sqlInsertSupply     = "INSERT INTO supplies"
	sqlUpdateSupply     = "UPDATE supplies"
	sqlDeleteSupply     = "DELETE FROM supplies"
)

type looseQueryMatcher struct{}

func (m looseQueryMatcher) Match(expectedSQL, actualSQL string) error {
	expectedSQL = strings.TrimSpace(strings.Join(strings.Fields(strings.ReplaceAll(expectedSQL, "\n", " ")), " "))
	actualSQL = strings.TrimSpace(strings.Join(strings.Fields(strings.ReplaceAll(actualSQL, "\n", " ")), " "))
	if expectedSQL == actualSQL || strings.Contains(actualSQL, expectedSQL) {
		return nil
	}
	return fmt.Errorf("actual sql: %q does not equal expected %q", actualSQL, expectedSQL)
}

func testLogger(t *testing.T) *logger.Logger {
	t.Helper()
	log, err := logger.New(config.LoggerConfig{Level: "error", Format: "text", Outputs: []string{"stdout"}}, "supplier-repo-test")
	require.NoError(t, err)
	return log
}

func newMockPool(t *testing.T) pgxmock.PgxPoolIface {
	t.Helper()
	pool, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(looseQueryMatcher{}))
	require.NoError(t, err)
	return pool
}

func newSupplierRepo(t *testing.T, pool pgxmock.PgxPoolIface) *repository.SupplierRepository {
	t.Helper()
	return repository.NewSupplierRepository(pool)
}

func newMaterialRepo(t *testing.T, pool pgxmock.PgxPoolIface) *repository.MaterialRepository {
	t.Helper()
	return repository.NewMaterialRepository(pool, testLogger(t))
}

func newModuleRepository(t *testing.T, pool pgxmock.PgxPoolIface) *repository.ModuleRepository {
	t.Helper()
	return repository.NewModuleRepository(pool, testLogger(t))
}

func newSupplyRepository(t *testing.T, pool pgxmock.PgxPoolIface) *repository.SupplyRepository {
	t.Helper()
	return repository.NewSupplyRepository(pool, testLogger(t))
}

func sampleSupplier(id uuid.UUID) models.Supplier {
	return models.Supplier{ID: id, Login: "supplier", Password: "hash", Email: "supplier@example.com"}
}

func sampleModule(id uuid.UUID) models.Module {
	return models.Module{ID: id, Name: "Module", Code: "MOD-1", Price: 10, StockCount: 2}
}

func sampleMaterial(id uuid.UUID) models.Material {
	return models.Material{ID: id, Name: "Material", Code: "MAT-1", Price: 10, StockCount: 5}
}

func sampleSupply(id, supplierID uuid.UUID) models.Supply {
	return models.Supply{ID: id, SupplierID: supplierID, Code: "SUP-1", Date: time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC), Price: 20}
}

func supplierRows(item models.Supplier) *pgxmock.Rows {
	return pgxmock.NewRows([]string{"id", "login", "password", "email"}).AddRow(item.ID, item.Login, item.Password, item.Email)
}

func materialRows(item models.Material) *pgxmock.Rows {
	return pgxmock.NewRows([]string{"id", "name", "code", "price", "stock_count"}).AddRow(item.ID, item.Name, item.Code, item.Price, item.StockCount)
}

func moduleRows(item models.Module) *pgxmock.Rows {
	return pgxmock.NewRows([]string{"id", "name", "code", "price", "stock_count"}).AddRow(item.ID, item.Name, item.Code, item.Price, item.StockCount)
}

func supplyRows(item models.Supply) *pgxmock.Rows {
	return pgxmock.NewRows([]string{"id", "supplier_id", "code", "date", "price"}).AddRow(item.ID, item.SupplierID, item.Code, item.Date, item.Price)
}
