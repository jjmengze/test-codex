package middleware

import "github.com/gin-gonic/gin"

func CheckHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		productCode := c.MustGet("productCode").(string)
		//todo maybe this is mock block?
		customerID := "789"
		sourceID := "ghi"

		c.Set("customerID", customerID)
		c.Set("productCode", productCode)
		c.Set("sourceID", sourceID)
	}
}
