package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	mockusecase "log-receiver/mock/internal_/usecase"
	"log-receiver/pkg/auth"
	"log-receiver/pkg/logger/slog"
	"log-receiver/pkg/middleware"
)

func TestValidateToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		param    string
		ctxValue interface{}
		// 下面兩個欄位給 mock 回傳用
		validatorOK  bool
		validatorErr error

		wantStatus     int
		wantBodySubstr string
	}{
		{
			name:           "Missing payload",
			param:          "prod",
			ctxValue:       nil,
			wantStatus:     http.StatusInternalServerError,
			wantBodySubstr: "missing ID P Token Payload",
		},
		{
			name:           "Invalid type",
			param:          "prod",
			ctxValue:       "not a payload",
			wantStatus:     http.StatusInternalServerError,
			wantBodySubstr: "can't convert to IDP Token Payload",
		},
		{
			name:           "Success",
			param:          "prod",
			ctxValue:       &auth.IDPTokenPayload{ProducerProductID: "prod"},
			validatorOK:    true,
			validatorErr:   nil,
			wantStatus:     http.StatusOK,
			wantBodySubstr: `{"payload":{"cpid":"","ppid":"prod","cid":"","uid":"","pl":"","it":0,"et":0}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 準備 HTTP recorder 與 Gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", fmt.Sprintf("/validate_token/%s", tc.param), nil)

			// 將 payload 放到 context
			if tc.ctxValue != nil {
				c.Set(middleware.CtxKeyIDPTokenPayload, tc.ctxValue)
			}
			c.Params = gin.Params{{Key: "productCode", Value: tc.param}}

			// 用 Mockey 生成的 mock
			validatorMock := mockusecase.NewValidator(t)

			svc := validateService{
				logger:           slog.GetGlobalLogger(),
				usecaseValidator: validatorMock,
				isTestPem:        true,
			}

			// 執行 handler
			svc.validateToken(c)

			// 驗證 HTTP status & body
			assert.Equal(t, tc.wantStatus, w.Code, "status code")
			assert.Contains(t, w.Body.String(), tc.wantBodySubstr, "body substring")

		})
	}
}
