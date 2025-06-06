package auth

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	jwt "github.com/golang-jwt/jwt"
)

const (
	xLogRProductID = "clr"
)

var priKeyAbsPath = os.Getenv("JWT_PRIVATE_KEY_PATH") //default "config/dummy_private_key.pem"

func init() {
	// 如果沒設環境變數，就預設從這個檔案的位置推導出路徑
	if priKeyAbsPath == "" {
		// 找出當前 go 檔案的絕對路徑
		_, currentFile, _, ok := runtime.Caller(0)
		if !ok {
			log.Fatal("Cannot get runtime caller info")
		}
		// 假設 config 資料夾在專案根目錄下
		baseDir := filepath.Dir(currentFile)              // 例如 log-receiver.go/cmd/server
		projectRoot := filepath.Join(baseDir, "..", "..") // 回到 log-receiver.go/
		priKeyAbsPath = filepath.Join(projectRoot, "config", "dummy_private_key.pem")
	}

	absPath, err := filepath.Abs(priKeyAbsPath)
	if err != nil {
		log.Fatalf("Failed to resolve absolute path: %v", err)
	}
	priKeyAbsPath = absPath

	if _, err := os.Stat(priKeyAbsPath); os.IsNotExist(err) {
		log.Fatalf("Private key file does not exist at %s", priKeyAbsPath)
	}

	log.Println("Using private key path:", priKeyAbsPath)
}

// GenJWTToken generate jwt token for testing
func GenIDPJWTToken(ppid, cpid, cid, uid string, et int64) string {
	verifyBytes, err := os.ReadFile(priKeyAbsPath)
	if err != nil {
		log.Fatalf("Fail to decrypt JWT token: %v", err)
	}

	priKey, err := jwt.ParseRSAPrivateKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatalf("Fail to decrypt JWT token: %v", err)
	}

	claims := jwt.MapClaims{
		"ppid": ppid,
		"cpid": cpid,
		"cid":  cid,
		"uid":  uid,
		"it":   time.Now().UTC().Unix(),
		"et":   et,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(priKey)
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}
	return tokenString
}

func DecryptJwtToken(tokenString string, rsaPublicKey *rsa.PublicKey) (*TokenPayload, error) {
	payload := &TokenPayload{}

	token, err := validateToken(tokenString, rsaPublicKey)
	if err != nil {
		return payload, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if val, ok := claims["prod_id"]; ok {
			payload.ProductId = val.(string)
		}
		if val, ok := claims["company_id"]; ok {
			payload.CompanyId = val.(string)
		}
		if val, ok := claims["computer_id"]; ok {
			payload.ComputerId = val.(string)
		}
	}
	return payload, err
}

func DecryptIDPJWTToken(tokenString string, rsaPublicKey *rsa.PublicKey) (*IDPTokenPayload, error) {
	payload := &IDPTokenPayload{}

	token, err := validateToken(tokenString, rsaPublicKey)
	if err != nil {
		return payload, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if val, ok := claims["ppid"]; ok {
			payload.ProducerProductID = val.(string)
		}
		if val, ok := claims["cpid"]; ok {
			payload.ConsumerProductID = val.(string)
		}
		if val, ok := claims["cid"]; ok {
			payload.CustomerID = val.(string)
		}
		if val, ok := claims["uid"]; ok {
			payload.UserID = val.(string)
		}
		if val, ok := claims["pl"]; ok {
			payload.Payload = val.(string)
		}
		if val, ok := claims["it"]; ok {
			payload.IssueTime = int64(val.(float64))
		}
		if val, ok := claims["et"]; ok {
			payload.ExpiredTime = int64(val.(float64))
		}
	}

	if payload.ExpiredTime > 0 && payload.ExpiredTime < time.Now().Unix() {
		return payload, jwt.NewValidationError("Token is expired", jwt.ValidationErrorExpired)
	}

	return payload, err
}

func validateToken(tokenString string, rsaPublicKey *rsa.PublicKey) (*jwt.Token, error) {
	// validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return rsaPublicKey, nil
	})

	// error handling
	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			log.Fatalf("Invalid JWT Token: %v", err)
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			log.Fatalf("Expired JWT Token: %v", err)
		}
	}

	return token, err
}
