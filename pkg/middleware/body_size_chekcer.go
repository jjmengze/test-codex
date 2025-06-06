package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO Implement this middleware to validate and extract the request body
// This function should:
// 1. Read the request body
// 2. Check the data is less then 30 MB, otherwise return an corresponding error code.
const maxSize = 30 << 20 // 30 MB
func CheckRequestBody() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Limit the size of the request body reader
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

		// Read the body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "Request body too large",
			})
			return
		}
		// Replace the body for other middlewares/handlers that might want to read it
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		c.Next()
	}
}
