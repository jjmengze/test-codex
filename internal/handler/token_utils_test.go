package handler

import (
	"context"
	"crypto/rsa"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log-receiver/pkg/logger/slog"

	"log-receiver/internal/usecase"
	"log-receiver/pkg/auth"
)

type mockValidator struct {
	shouldPass bool
	err        error
}

func (m *mockValidator) Validate(ctx context.Context, productCode string) (bool, error) {
	return m.shouldPass, m.err
}

func mockTokenDecryptorFail(tokenString string, pubKey *rsa.PublicKey) (*auth.IDPTokenPayload, error) {
	return nil, errors.New("invalid token")
}

func mockTokenDecryptorSuccess(tokenString string, pubKey *rsa.PublicKey) (*auth.IDPTokenPayload, error) {
	return &auth.IDPTokenPayload{ProducerProductID: "product-123"}, nil
}

func TestValidateIDPToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	originalDecrypt := decryptIDPJWTToken
	defer func() {
		decryptIDPJWTToken = originalDecrypt
	}()

	tests := []struct {
		name               string
		authHeader         string
		productCode        string
		mockDecryptor      func(string, *rsa.PublicKey) (*auth.IDPTokenPayload, error)
		validator          usecase.Validator
		expectedStatusCode int
		expectedErr        string
	}{
		{
			name:               "missing authorization header",
			authHeader:         "",
			productCode:        "product-123",
			mockDecryptor:      mockTokenDecryptorSuccess,
			validator:          &mockValidator{shouldPass: true},
			expectedStatusCode: http.StatusUnauthorized,
			expectedErr:        "missing or malformed Authorization header",
		},
		{
			name:               "invalid token format",
			authHeader:         "InvalidTokenFormat",
			productCode:        "product-123",
			mockDecryptor:      mockTokenDecryptorSuccess,
			validator:          &mockValidator{shouldPass: true},
			expectedStatusCode: http.StatusUnauthorized,
			expectedErr:        "missing or malformed Authorization header",
		},
		{
			name:               "token decryption fails",
			authHeader:         "Bearer invalidtoken",
			productCode:        "product-123",
			mockDecryptor:      mockTokenDecryptorFail,
			validator:          &mockValidator{shouldPass: true},
			expectedStatusCode: http.StatusUnauthorized,
			expectedErr:        "invalid token",
		},
		{
			name:        "product ID mismatch",
			authHeader:  "Bearer validtoken",
			productCode: "mismatch",
			mockDecryptor: func(token string, pubKey *rsa.PublicKey) (*auth.IDPTokenPayload, error) {
				return &auth.IDPTokenPayload{ProducerProductID: "product-123"}, nil
			},
			validator:          &mockValidator{shouldPass: true},
			expectedStatusCode: http.StatusUnauthorized,
			expectedErr:        "token product code does not match request path",
		},
		{
			name:               "validation fails",
			authHeader:         "Bearer validtoken",
			productCode:        "product-123",
			mockDecryptor:      mockTokenDecryptorSuccess,
			validator:          &mockValidator{shouldPass: false},
			expectedStatusCode: http.StatusNotAcceptable,
			expectedErr:        "product not supported",
		},
		{
			name:               "valid case",
			authHeader:         "Bearer validtoken",
			productCode:        "product-123",
			mockDecryptor:      mockTokenDecryptorSuccess,
			validator:          &mockValidator{shouldPass: true},
			expectedStatusCode: http.StatusOK,
			expectedErr:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decryptIDPJWTToken = tt.mockDecryptor

			req := httptest.NewRequest("GET", "/api/v2/"+tt.productCode, nil)
			req.Header.Set("Authorization", tt.authHeader)
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			logger := slog.GetGlobalLogger()
			ctx.Request = req
			ctx.Params = []gin.Param{
				{Key: "productCode", Value: tt.productCode},
			}

			payload, statusCode, err := validateIDPToken(logger, ctx, tt.validator, true)

			assert.Equal(t, tt.expectedStatusCode, statusCode)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, payload)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, payload)
				assert.Equal(t, "product-123", payload.ProducerProductID)
			}
		})
	}
}
