package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCheckRequestBodyMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		bodySize       int
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Body under limit",
			bodySize:       (10 << 20), // 10 MB
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"OK"}`,
		},
		{
			name:           "Body at limit",
			bodySize:       (30 << 20), // 30 MB
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"OK"}`,
		},
		{
			name:           "Body over limit",
			bodySize:       (31 << 20), // 31 MB
			expectedStatus: http.StatusRequestEntityTooLarge,
			expectedBody:   `{"error":"Request body too large"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create dummy body of requested size
			body := bytes.Repeat([]byte("a"), tt.bodySize)

			// Setup router with middleware and a dummy handler
			router := gin.New()
			router.Use(CheckRequestBody())
			router.POST("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "OK"})
			})

			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
