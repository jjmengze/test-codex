package handler

import (
	"crypto/rsa"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	usecaseMock "log-receiver/mock/internal_/usecase"
	"log-receiver/pkg/auth"
)

func TestValidateIDPTokenHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// prepare environment for reading public key
	os.Setenv("JWT_PUBLIC_KEY_PATH", "../../config/dummy_public_key.pem")

	tests := []struct {
		name         string
		header       string
		tokenPayload *auth.IDPTokenPayload
		validatorOK  bool
		validatorErr error
		expectedCode int
		expectedErr  bool
	}{
		{
			name:         "missing authorization header",
			header:       "",
			expectedCode: http.StatusUnauthorized,
			expectedErr:  true,
		},
		{
			name:         "token product code mismatch",
			header:       "Bearer dummy",
			tokenPayload: &auth.IDPTokenPayload{ProducerProductID: "foo"},
			validatorOK:  true,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  true,
		},
		{
			name:         "unsupported product",
			header:       "Bearer dummy",
			tokenPayload: &auth.IDPTokenPayload{ProducerProductID: "sao"},
			validatorOK:  false,
			expectedCode: http.StatusNotAcceptable,
			expectedErr:  true,
		},
		{
			name:         "valid token",
			header:       "Bearer dummy",
			tokenPayload: &auth.IDPTokenPayload{ProducerProductID: "sao"},
			validatorOK:  true,
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock decrypt function
			decryptIDPJWTToken = func(token string, key *rsa.PublicKey) (*auth.IDPTokenPayload, error) {
				if tt.tokenPayload == nil {
					return nil, errors.New("fail")
				}
				return tt.tokenPayload, nil
			}
			defer func() { decryptIDPJWTToken = auth.DecryptIDPJWTToken }()

			// mock validator
			mVal := usecaseMock.NewValidator(t)
			if tt.header != "" {
				mVal.On("Validate", testifymock.Anything, tt.tokenPayload.ProducerProductID).Return(tt.validatorOK, tt.validatorErr)
			}

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodGet, "/activity_log/sao", nil)
			c.Params = gin.Params{{Key: "productCode", Value: "sao"}}
			if tt.header != "" {
				c.Request.Header.Set("Authorization", tt.header)
			}

			payload, code, err := validateIDPToken(c, mVal)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedCode, code)
			if code == http.StatusOK {
				assert.Equal(t, tt.tokenPayload, payload)
			}
		})
	}
}
