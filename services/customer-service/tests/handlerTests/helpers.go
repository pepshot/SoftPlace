package handlertests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pepshot/SoftPlace/services/customer-service/tests/mocks"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/dto"
	httpserver "github.com/pepshot/SoftPlace/services/customer-service/internal/transport/http"
)

const userIDKey = "userID"

func newHandler(t *testing.T) (*httpserver.Handler, *mocks.CustomerUseCaseMock,
	*mocks.FurnitureUseCaseMock, *mocks.GarnitureUseCaseMock,
	*mocks.ShipmentUseCaseMock) {
	t.Helper()

	customer := &mocks.CustomerUseCaseMock{}
	furniture := &mocks.FurnitureUseCaseMock{}
	garniture := &mocks.GarnitureUseCaseMock{}
	shipment := &mocks.ShipmentUseCaseMock{}

	handler := httpserver.NewHandler(customer, furniture, garniture, shipment, nil)
	return handler, customer, furniture, garniture, shipment
}

func newTestContext(method, path string, body []byte) (*gin.Context, *gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, router := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(method, path, bytes.NewReader(body))
	if len(body) > 0 {
		ctx.Request.Header.Set("Content-Type", "application/json")
	}

	return ctx, router, recorder
}

func setUserID(ctx *gin.Context, userID string) {
	ctx.Set(userIDKey, userID)
}

func setUserIDValue(ctx *gin.Context, value any) {
	ctx.Set(userIDKey, value)
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()

	payload, err := json.Marshal(value)
	require.NoError(t, err)
	return payload
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
