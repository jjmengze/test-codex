package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"log-receiver/pkg/logger"
)

// LoggerMiddleware logs all incoming HTTP requests.
func LoggerMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start time
		startTime := time.Now()

		// Process request
		c.Next()

		// End time
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// Log format
		// Log details in structured key-value format
		logger.InfoM("request log", map[interface{}]interface{}{
			"timestamp": startTime.Format(time.RFC3339),
			"status":    c.Writer.Status(),
			"latency":   latency.String(),
			"clientIP":  c.ClientIP(),
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
		})
	}
}
