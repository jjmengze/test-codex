package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log-receiver/pkg/auth"
	"log-receiver/pkg/logger/slog"
)

func TestCheckHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupContext   func(c *gin.Context)
		expectedStatus int
		expectAbort    bool
	}{
		{
			name: "Valid payload",
			setupContext: func(c *gin.Context) {
				c.Set("productCode", "prod-001")
				c.Set(CtxKeyIDPTokenPayload, &auth.IDPTokenPayload{
					CustomerID: "cust-999",
				})
				c.Request.Header.Set("X-Source-ID", "gateway")
			},
			expectedStatus: http.StatusOK,
			expectAbort:    false,
		},
		{
			name: "Missing payload",
			setupContext: func(c *gin.Context) {
				c.Set("productCode", "prod-002")
			},
			expectedStatus: http.StatusBadRequest,
			expectAbort:    true,
		},
		{
			name: "Invalid payload type",
			setupContext: func(c *gin.Context) {
				c.Set("productCode", "prod-003")
				c.Set(CtxKeyIDPTokenPayload, "not-a-payload")
			},
			expectedStatus: http.StatusInternalServerError,
			expectAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.GetGlobalLogger()

			router := gin.New()
			router.Use(func(c *gin.Context) {
				tt.setupContext(c)
				c.Next()
			})
			router.Use(CheckHeader(logger))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
