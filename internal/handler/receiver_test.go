package handler

import (
	"bytes"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	usecaseMock "log-receiver/mock/internal_/usecase"
	"log-receiver/pkg/auth"
	"log-receiver/pkg/logger/slog"
)

func injectTestContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("{productCode}", "testProduct")
		c.Set("traceID", "testTrace")
		c.Set("customerID", "testCustomer")
		c.Set("encoding", "utf-8")
		c.Set("subType", "testType")
		c.Set("sourceID", "testSource")
		c.Next()
	}
}
func TestActivityLogRequestSize(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_PUBLIC_KEY_PATH", "../../config/dummy_public_key.pem")

	type testContext struct {
		t             *testing.T
		mockReceiver  *usecaseMock.Receiver
		mockValidator *usecaseMock.Validator
		router        *gin.Engine
	}

	tests := []struct {
		name           string
		bodySizeMB     int
		expectedStatus int
		expectContains string
		setupMock      func(ctx *testContext)
	}{
		{
			name:           "OK request under 30MB",
			bodySizeMB:     5,
			expectedStatus: http.StatusOK,
			expectContains: "Success",
			setupMock: func(ctx *testContext) {
				ctx.mockReceiver.On("PutData", testifymock.Anything, testifymock.Anything).Return(nil)
				ctx.mockValidator.On("Validate", testifymock.Anything, "test").Return(true, nil)
			},
		},
		{
			name:           "Request over 30MB should be rejected",
			bodySizeMB:     31,
			expectedStatus: http.StatusRequestEntityTooLarge,
			expectContains: "Request body too large",
			setupMock: func(ctx *testContext) {
				ctx.mockValidator.On("Validate", testifymock.Anything, "test").Return(true, nil)
				// 不 mock Receiver，因為不會觸發
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			mockValidator := usecaseMock.NewValidator(t)
			mockReceiver := usecaseMock.NewReceiver(t)
			mockLogger := slog.GetGlobalLogger()

			ctx := &testContext{
				t:             t,
				mockReceiver:  mockReceiver,
				mockValidator: mockValidator,
				router:        router,
			}
			ctx.router.Use(injectTestContext())

			// 偽造 token 驗證
			decryptIDPJWTToken = func(string, *rsa.PublicKey) (*auth.IDPTokenPayload, error) {
				return &auth.IDPTokenPayload{ProducerProductID: "test"}, nil
			}
			defer func() { decryptIDPJWTToken = auth.DecryptIDPJWTToken }()

			// 呼叫各案例自定義的 mock 設定
			if tt.setupMock != nil {
				tt.setupMock(ctx)
			}

			NewReceiverService(mockLogger, ctx.router, ctx.mockReceiver, ctx.mockValidator, false)

			body := strings.Repeat("A", tt.bodySizeMB*1024*1024)

			req := httptest.NewRequest(http.MethodPost, "/activity_log/test", bytes.NewBufferString(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer dummy")

			rec := httptest.NewRecorder()
			ctx.router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectContains != "" {
				assert.Contains(t, rec.Body.String(), tt.expectContains)
			}
		})
	}
}
