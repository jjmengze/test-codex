package middleware

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log-receiver/internal/usecase"
	"log-receiver/pkg/auth"
	"log-receiver/pkg/aws"
)

const CtxKeyIDPTokenPayload = "idpPayload"

func getJWTToken(token string) string {
	return strings.TrimPrefix(token, "Bearer ")
}

// ValidateTokenController is a middleware that validates JWT token and product code.
func ValidateTokenController(validator usecase.Validator, pubKeyAbsPath string, isTestPem bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var publicKey *rsa.PublicKey
		var err error
		if isTestPem {
			pubKeyAbsPath = ""
			productCode := c.Param("productCode")
			publicKey, err = aws.GetVerifyKeyByProductCode(productCode)
			if err != nil {
				err := fmt.Errorf("cannot get public key for product code %s", productCode)
				c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
				return
			}
		}
		payload, code, err := validateIDPToken(c, pubKeyAbsPath, publicKey, validator)
		if err != nil {
			c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		c.Set(CtxKeyIDPTokenPayload, payload)
		c.Next()
	}
}

// validateIDPToken validates JWT token and productCode.
// It returns the token payload or an HTTP status code and error when validation fails.
func validateIDPToken(c *gin.Context, pubKeyAbsPath string, publicKey *rsa.PublicKey, validator usecase.Validator) (*auth.IDPTokenPayload, int, error) {
	ctx := c.Request.Context()

	authHeader := c.GetHeader("Authorization")
	productCode := c.Param("productCode")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, http.StatusUnauthorized, errors.New("missing or malformed Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if pubKeyAbsPath != "" && publicKey == nil {
		pubKeyBytes, err := os.ReadFile(pubKeyAbsPath)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		publicKey, err = jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	payload, err := auth.DecryptIDPJWTToken(tokenString, publicKey)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	if payload.ProducerProductID != productCode {
		return nil, http.StatusUnauthorized, errors.New("token product code does not match request path")
	}

	ok, err := validator.Validate(ctx, payload.ProducerProductID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if !ok {
		return nil, http.StatusNotAcceptable, errors.New("product not supported")
	}

	return payload, http.StatusOK, nil
}
