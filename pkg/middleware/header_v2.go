package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func AssignV2Header() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("token", getJWTToken(c.GetHeader("Authorization")))
		c.Set("productCode", c.Param("productCode"))
		c.Set("encoding", c.GetHeader("Content-Encoding"))
		c.Set("bytes", c.Request.ContentLength)
		parts := strings.Split(strings.Trim(c.Request.URL.Path, "/"), "/")
		logType := parts[len(parts)-1]
		c.Set("logType", logType)
		c.Set("subType", c.Param("subType"))
		c.Next()
	}
}
