package handler

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"

	"log-receiver/internal/usecase"
	"log-receiver/pkg/auth"
)

const ctxKeyIDPTokenPayload = "idpPayload"

var pubKeyAbsPath = os.Getenv("JWT_PUBLIC_KEY_PATH")
var decryptIDPJWTToken = auth.DecryptIDPJWTToken

func init() {
	if pubKeyAbsPath == "" {
		_, currentFile, _, ok := runtime.Caller(0)
		if !ok {
			log.Fatal("Cannot get runtime caller info for public key")
		}
		baseDir := filepath.Dir(currentFile)
		projectRoot := filepath.Join(baseDir, "..", "..")
		pubKeyAbsPath = filepath.Join(projectRoot, "config", "dummy_public_key.pem")
	}
	// 轉為絕對路徑（保險做法）
	absPath, err := filepath.Abs(pubKeyAbsPath)
	if err != nil {
		log.Fatalf("Failed to resolve absolute path: %v", err)
	}
	pubKeyAbsPath = absPath

	// 檢查檔案是否存在
	if _, err := os.Stat(pubKeyAbsPath); os.IsNotExist(err) {
		log.Fatalf("Private key file does not exist at %s", pubKeyAbsPath)
	}
}

// validateIDPToken validates JWT token and productCode.
// It returns the token payload or an HTTP status code and error when validation fails.
func validateIDPToken(c *gin.Context, validator usecase.Validator) (*auth.IDPTokenPayload, int, error) {
	ctx := c.Request.Context()

	authHeader := c.GetHeader("Authorization")
	productCode := c.Param("productCode")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, http.StatusUnauthorized, errors.New("Missing or malformed Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	pubKeyBytes, err := os.ReadFile(pubKeyAbsPath)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	payload, err := decryptIDPJWTToken(tokenString, pubKey)
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

// ValidateTokenController is a middleware that validates JWT token and product code.
func ValidateTokenController(validator usecase.Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, code, err := validateIDPToken(c, validator)
		if err != nil {
			c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		c.Set(ctxKeyIDPTokenPayload, payload)
		c.Next()
	}
}
