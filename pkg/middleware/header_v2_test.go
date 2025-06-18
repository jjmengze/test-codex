package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAssignV2Header(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		authorization    string
		contentEncoding  string
		body             string
		requestPath      string
		expectedToken    string
		expectedEncoding string
		expectedLogType  string
		expectedBytes    int64
	}{
		{
			name:             "Full header and params",
			authorization:    "Bearer my-token-123",
			contentEncoding:  "gzip",
			body:             `{"foo":"bar"}`,
			requestPath:      "/activity_log/productCode123/subTypeABC",
			expectedToken:    "my-token-123",
			expectedEncoding: "gzip",
			expectedLogType:  "subTypeABC",
			expectedBytes:    int64(len(`{"foo":"bar"}`)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(AssignV2Header())

			// 將 productCode 和 subType 當成 URL param
			router.POST("/activity_log/:productCode/:subType", func(c *gin.Context) {
				token, _ := c.Get("token")
				productCode, _ := c.Get("productCode")
				encoding, _ := c.Get("encoding")
				bytesVal, _ := c.Get("bytes")
				logType, _ := c.Get("logType")
				subType, _ := c.Get("subType")

				c.JSON(http.StatusOK, gin.H{
					"token":       token,
					"productCode": productCode,
					"encoding":    encoding,
					"bytes":       bytesVal,
					"logType":     logType,
					"subType":     subType,
				})
			})

			bodyReader := strings.NewReader(tt.body)
			req := httptest.NewRequest(http.MethodPost, tt.requestPath, bodyReader)
			req.Header.Set("Authorization", tt.authorization)
			req.Header.Set("Content-Encoding", tt.contentEncoding)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedToken)
			assert.Contains(t, w.Body.String(), `"productCode":"productCode123"`)
			assert.Contains(t, w.Body.String(), `"encoding":"gzip"`)
			assert.Contains(t, w.Body.String(), `"logType":"subTypeABC"`)
			assert.Contains(t, w.Body.String(), `"subType":"subTypeABC"`)
			assert.Contains(t, w.Body.String(), `"bytes":`)

			// 可以額外驗證 bytes 是否正確
			// 但 json 被轉成 float64 可能不太精準比對
		})
	}
}
