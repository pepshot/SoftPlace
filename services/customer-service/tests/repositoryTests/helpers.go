package repositorytests

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/model"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/model/relations"
	"github.com/pepshot/SoftPlace/services/customer-service/internal/repository"
	"github.com/pepshot/SoftPlace/shared/config"
	"github.com/pepshot/SoftPlace/shared/logger"
)

const (
	sqlSelectCustomer  = "SELECT id, login, password, email FROM customers WHERE id = $1;"
	sqlSelectFromLogin = "SELECT id, login, password, email FROM customers WHERE login = $1;"
	sqlExistFromLogin  = "SELECT EXISTS(SELECT id FROM customers WHERE login = $1);"
	sqlExistFromEmail  = "SELECT EXISTS(SELECT id FROM customers WHERE email = $1);"
	sqlInsertCustomer  = "INSERT INTO customers"

	sqlSelectFurniture     = "SELECT id, name, code, price, stock_count FROM furniture ORDER BY name;"
	sqlSelectFurnitureByID = "SELECT id, name, code, price, stock_count FROM furniture WHERE id = $1;"
	sqlIncreaseFurniture   = "UPDATE furniture SET stock_count = stock_count + $2 WHERE id = $1;"
	sqlDecreaseFurniture   = "UPDATE furniture SET stock_count = stock_count - $2 WHERE id = $1 AND stock_count >= $2;"
	sqlDeleteFurniture     = "DELETE FROM furniture"

	sqlSelectFurnitureModule = "SELECT furniture_id, module_id, count FROM furniture_modules WHERE furniture_id = $1;"
	sqlInsertFurnitureModule = "INSERT INTO furniture_modules"
	sqlDeleteFurnitureModule = "DELETE FROM furniture_modules WHERE furniture_id = $1;"

	sqlSelectGarniture       = "SELECT id, name, code, price, stock_count FROM garnitures ORDER BY name;"
	sqlSelectGarnitureFromID = "SELECT id, name, code, price, stock_count FROM garnitures WHERE id = $1;"
	sqlInsertGarniture       = "INSERT INTO garnitures"
	sqlUpdateGarniture       = "UPDATE garnitures"
	sqlIncreaseGarniture     = "UPDATE garnitures SET stock_count = stock_count + $2 WHERE id = $1;"
	sqlDecreaseGarniture     = "UPDATE garnitures SET stock_count = stock_count - $2 WHERE id = $1 AND stock_count >= $2;"
	sqlDeleteGarniture       = "DELETE FROM garnitures"

	sqlSelectGarnitureFurniture = "SELECT garniture_id, furniture_id, count FROM garniture_furniture WHERE garniture_id = $1;"
	sqlInsertGarnitureFurniture = "INSERT INTO garniture_furniture"
	sqlDeleteGarnitureFurniture = "DELETE FROM garniture_furniture WHERE garniture_id = $1;"

	sqlSelectShipment       = "SELECT id, customer_id, code, date, price FROM shipments ORDER BY date DESC;"
	sqlSelectShipmentFromID = "SELECT id, customer_id, code, date, price FROM shipments WHERE id = $1;"
	sqlInsertShipment       = "INSERT INTO shipments"
	sqlUpdateShipment       = "UPDATE shipments"
	sqlDeleteShipment       = "DELETE FROM shipments"

	sqlSelectShipmentGarniture = "SELECT shipment_id, garniture_id, count FROM shipment_garnitures WHERE shipment_id = $1;"
	sqlInsertShipmentGarniture = "INSERT INTO shipment_garnitures"
	sqlDeleteShipmentGarniture = "DELETE FROM shipment_garnitures WHERE shipment_id = $1;"
)

type looseQueryMatcher struct{}

func (m looseQueryMatcher) Match(expectedSQL, actualSQL string) error {
	normalize := func(sql string) string {
		sql = strings.TrimSpace(sql)
		sql = strings.ReplaceAll(sql, `\+`, `+`)
		return strings.Join(strings.Fields(sql), " ")
	}

	expectedSQL = normalize(expectedSQL)
	actualSQL = normalize(actualSQL)

	if expectedSQL == actualSQL || strings.Contains(actualSQL, expectedSQL) {
		return nil
	}

	return fmt.Errorf("actual sql: %q does not equal to expected %q", actualSQL, expectedSQL)
}

func testLogger(t *testing.T) *logger.Logger {
	t.Helper()

	log, err := logger.New(config.LoggerConfig{
		Level:   "error",
		Format:  "text",
		Outputs: []string{"stdout"},
	}, "customer-repo-test")
	require.NoError(t, err)
	return log
}

func newMockPool(t *testing.T) pgxmock.PgxPoolIface {
	t.Helper()

	pool, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(looseQueryMatcher{}))
	require.NoError(t, err)
	return pool
}

func newRepository(t *testing.T, pool pgxmock.PgxPoolIface) *repository.CustomerRepository {
	t.Helper()

	return repository.NewCustomerRepository(pool, testLogger(t))
}

func newFurnitureRepository(t *testing.T, pool pgxmock.PgxPoolIface) *repository.FurnitureRepository {
	t.Helper()

	repo := repository.NewFurnitureRepository(pool, testLogger(t))
	return &repo
}

func newGarnitureRepository(t *testing.T, pool pgxmock.PgxPoolIface) *repository.GarnitureRepository {
	t.Helper()

	return repository.NewGarnitureRepository(pool, testLogger(t))
}

func newShipmentRepository(t *testing.T, pool pgxmock.PgxPoolIface) *repository.ShipmentRepository {
	t.Helper()

	return repository.NewShipmentRepository(pool, testLogger(t))
}

func newFurnitureModuleRepository(t *testing.T, pool pgxmock.PgxPoolIface) *repository.FurnitureModuleRepository {
	t.Helper()

	return repository.NewFurnitureModuleRepository(pool, testLogger(t))
}

func newGarnitureFurnitureRepository(t *testing.T, pool pgxmock.PgxPoolIface) *repository.GarnitureFurnitureRepository {
	t.Helper()

	return repository.NewGarnitureFurnitureRepository(pool, testLogger(t))
}

func newShipmentGarnitureRepository(t *testing.T, pool pgxmock.PgxPoolIface) *repository.ShipmentGarnitureRepository {
	t.Helper()

	return repository.NewShipmentGarnitureRepository(pool, testLogger(t))
}

func newCustomer(id uuid.UUID) model.Customer {
	return model.Customer{
		ID:       id,
		Login:    "login",
		Password: "hash",
		Email:    "user@example.com",
	}
}

func sampleFurniture(id uuid.UUID) model.Furniture {
	return model.Furniture{
		ID:         id,
		Name:       "Chair",
		Code:       "CH-1",
		Price:      10,
		StockCount: 2,
	}
}

func sampleGarniture(id uuid.UUID) model.Garniture {
	return model.Garniture{
		ID:         id,
		Name:       "Set",
		Code:       "GS-1",
		Price:      100,
		StockCount: 1,
	}
}

func sampleShipment(id uuid.UUID, customerID uuid.UUID) model.Shipment {
	return model.Shipment{
		ID:         id,
		CustomerID: customerID,
		Code:       "SH-1",
		Date:       time.Date(2026, 5, 9, 0, 0, 0, 0, time.UTC),
		Price:      200,
	}
}

func customerRows(customer model.Customer) *pgxmock.Rows {
	return pgxmock.NewRows([]string{"id", "login", "password", "email"}).
		AddRow(customer.ID, customer.Login, customer.Password, customer.Email)
}

func furnitureRows(items ...model.Furniture) *pgxmock.Rows {
	rows := pgxmock.NewRows([]string{"id", "name", "code", "price", "stock_count"})
	for _, item := range items {
		rows.AddRow(item.ID, item.Name, item.Code, item.Price, item.StockCount)
	}
	return rows
}

func garnitureRows(items ...model.Garniture) *pgxmock.Rows {
	rows := pgxmock.NewRows([]string{"id", "name", "code", "price", "stock_count"})
	for _, item := range items {
		rows.AddRow(item.ID, item.Name, item.Code, item.Price, item.StockCount)
	}
	return rows
}

func shipmentRows(items ...model.Shipment) *pgxmock.Rows {
	rows := pgxmock.NewRows([]string{"id", "customer_id", "code", "date", "price"})
	for _, item := range items {
		rows.AddRow(item.ID, item.CustomerID, item.Code, item.Date, item.Price)
	}
	return rows
}

func furnitureModuleRows(items ...relations.FurnitureModule) *pgxmock.Rows {
	rows := pgxmock.NewRows([]string{"furniture_id", "module_id", "count"})
	for _, item := range items {
		rows.AddRow(item.FurnitureID, item.ModuleID, item.Count)
	}
	return rows
}

func garnitureFurnitureRows(items ...relations.GarnitureFurniture) *pgxmock.Rows {
	rows := pgxmock.NewRows([]string{"garniture_id", "furniture_id", "count"})
	for _, item := range items {
		rows.AddRow(item.GarnitureID, item.FurnitureID, item.Count)
	}
	return rows
}

func shipmentGarnitureRows(items ...relations.ShipmentGarniture) *pgxmock.Rows {
	rows := pgxmock.NewRows([]string{"shipment_id", "garniture_id", "count"})
	for _, item := range items {
		rows.AddRow(item.ShipmentID, item.GarnitureID, item.Count)
	}
	return rows
}
