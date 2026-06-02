package handlertests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/pepshot/SoftPlace/services/supplier-service/internal/dto"
	httpserver "github.com/pepshot/SoftPlace/services/supplier-service/internal/transport/http"
	"github.com/pepshot/SoftPlace/services/supplier-service/tests/mocks"
)

const userIDKey = "userID"

func newHandler(t *testing.T) (*httpserver.Handler, *mocks.SupplierUseCaseMock, *mocks.ModuleUseCaseMock, *mocks.MaterialUseCaseMock, *mocks.SupplyUseCaseMock) {
	t.Helper()
	supplier := &mocks.SupplierUseCaseMock{}
	module := &mocks.ModuleUseCaseMock{}
	material := &mocks.MaterialUseCaseMock{}
	supply := &mocks.SupplyUseCaseMock{}
	return httpserver.NewHandler(module, material, supply, supplier, nil), supplier, module, material, supply
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

func setUserID(ctx *gin.Context, id string) { ctx.Set(userIDKey, id) }

func setUserIDValue(ctx *gin.Context, value any) { ctx.Set(userIDKey, value) }

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()
	payload, err := json.Marshal(value)
	require.NoError(t, err)
	return payload
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
