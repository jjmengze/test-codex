package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"log-receiver/pkg/auth"
	"log-receiver/pkg/logger"
)

func CheckHeader(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		productCode := c.MustGet("productCode").(string)
		sourceID := c.GetHeader("X-Source-ID") //if from gateway

		val, exists := c.Get(CtxKeyIDPTokenPayload)
		if !exists {
			err := fmt.Errorf("missing key in context: %s", CtxKeyIDPTokenPayload)
			logger.WarnF("error getting key  %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		idPTokenPayload, ok := val.(*auth.IDPTokenPayload)
		if !ok {
			err := fmt.Errorf("invalid payload type: %T", val)
			logger.ErrorF("error getting payload %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.Set("customerID", idPTokenPayload.CustomerID)
		c.Set("productCode", productCode)
		c.Set("sourceID", sourceID)
	}
}
