package middleware

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAssignTraceID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		inputTraceID    string
		expectSameValue bool
	}{
		{
			name:            "Has X-TRACE-ID header",
			inputTraceID:    "abc-123",
			expectSameValue: true,
		},
		{
			name:            "No X-TRACE-ID header",
			inputTraceID:    "",
			expectSameValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(AssignTraceID())

			// 偵測 middleware 是否有設值進 context
			router.GET("/test", func(c *gin.Context) {
				traceID, exists := c.Get("traceID")
				assert.True(t, exists)

				c.JSON(http.StatusOK, gin.H{"traceID": traceID})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.inputTraceID != "" {
				req.Header.Set("X-TRACE-ID", tt.inputTraceID)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			// 驗證 traceID
			body := w.Body.String()
			if tt.expectSameValue {
				assert.Contains(t, body, tt.inputTraceID)
			} else {
				// 判斷是否是 UUID 格式
				uuidRegex := regexp.MustCompile(`"[a-f0-9\-]{36}"`)
				assert.Regexp(t, uuidRegex, body)
			}
		})
	}
}
