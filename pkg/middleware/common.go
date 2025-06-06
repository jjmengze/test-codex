package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AssignTraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-TRACE-ID")
		if traceID == "" {
			traceID = genUUID()
		}
		c.Set("traceID", traceID)
	}
}

func genUUID() string {
	return uuid.Must(uuid.NewRandom()).String()
}
