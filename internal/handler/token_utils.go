package handler

import (
	"crypto/rsa"
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
	"log-receiver/pkg/aws"
	"log-receiver/pkg/logger"
)

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
func validateIDPToken(logger logger.Logger, c *gin.Context, validator usecase.Validator, isTestPem bool) (*auth.IDPTokenPayload, int, error) {
	ctx := c.Request.Context()

	authHeader := c.GetHeader("Authorization")
	productCode := c.Param("productCode")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		logger.WithContext(ctx).ErrorF("invalid authorization header: %s", authHeader)
		return nil, http.StatusUnauthorized, errors.New("missing or malformed Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	var pubKey *rsa.PublicKey
	var err error
	if isTestPem {
		pubKeyBytes, err := os.ReadFile(pubKeyAbsPath)
		if err != nil {
			logger.WithContext(ctx).ErrorF("failed to read public key file: %s", pubKeyAbsPath)
			return nil, http.StatusInternalServerError, err
		}
		pubKey, err = jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
		if err != nil {
			logger.WithContext(ctx).ErrorF("failed to parse public key file: %s", pubKeyAbsPath)
			return nil, http.StatusInternalServerError, err
		}
	} else {
		pubKey, err = aws.GetVerifyKeyByProductCode(productCode)
		if err != nil {
			logger.WithContext(ctx).ErrorF("invalid product code: %s", productCode)
			return nil, http.StatusUnauthorized, errors.New("invalid product code")
		}
	}

	payload, err := decryptIDPJWTToken(c.Request.Context(), logger, tokenString, pubKey)
	if err != nil {
		logger.WithContext(ctx).ErrorF("failed to decrypt token: %s", tokenString)
		return nil, http.StatusUnauthorized, err
	}

	if payload.ProducerProductID != productCode {
		logger.WithContext(ctx).ErrorF("invalid product code: %s", productCode)
		return nil, http.StatusUnauthorized, errors.New("token product code does not match request path")
	}

	ok, err := validator.Validate(ctx, payload.ProducerProductID)
	if err != nil {
		logger.WithContext(ctx).ErrorF("failed to validate token payload: %s", payload.ProducerProductID)
		return nil, http.StatusInternalServerError, err
	}
	if !ok {
		logger.WithContext(ctx).ErrorF("invalid product code: %s", payload.ProducerProductID)
		return nil, http.StatusNotAcceptable, errors.New("product not supported")
	}

	return payload, http.StatusOK, nil
}
