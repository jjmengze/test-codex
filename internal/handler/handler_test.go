package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"log-receiver/internal/usecase"
	"log-receiver/pkg/logger/slog"
)

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set up
	r := gin.New()
	mockLogger := slog.GetGlobalLogger()
	var mockReceiver usecase.Receiver // can be nil unless you need it
	var mockValidator usecase.Validator

	// Init handler with route
	NewHttpHandler(mockLogger, r, mockReceiver, mockValidator, "test/pem", true)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// Serve request
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}
