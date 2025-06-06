//go:build mock
// +build mock

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var EdlProductIDRegionMap = map[string]string{
	"us-east-1":      "nabu",
	"us-west-2":      "nabu",
	"eu-central-1":   "emea",
	"ap-northeast-1": "japan",
}

func init() {
	initTestPublicKeyMap()
}

type KeyPair struct {
	privateKey rsa.PrivateKey
	publicKey  rsa.PublicKey
}

func prepareProductIdKeyPair(productIdToKeyPair map[string]KeyPair, productId string) (rsa.PrivateKey, rsa.PublicKey) {
	keyPair, ok := productIdToKeyPair[productId]
	if ok {
		return keyPair.privateKey, keyPair.publicKey
	} else {
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		pub := &key.PublicKey
		if err != nil {
			panic(err)
		}
		productIdToKeyPair[productId] = KeyPair{
			privateKey: *key,
			publicKey:  *pub,
		}
		return *key, *pub
	}
}

func prepareTestJWT(producerProductId string, et int64, key rsa.PrivateKey, pub rsa.PublicKey, additionalClaims map[string]interface{}) (string, string, rsa.PublicKey) {
	bytes, err := x509.MarshalPKIXPublicKey(&pub)
	if err != nil {
		panic(err)
	}
	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: bytes,
		},
	)

	// setup JWT with dynamic claims
	claims := jwt.MapClaims{
		"ppid": producerProductId,
		"cpid": "clr",
		"cid":  "customerID",
		"uid":  "userID",
		"it":   time.Now().UTC().Unix(),
		"et":   et,
	}

	// Add additional claims
	for key, value := range additionalClaims {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(&key)
	if err != nil {
		log.Printf("Sign JWT failed for ppid %s: %v", producerProductId, err)
		panic(err)
	}
	return tokenString, string(publicKeyPEM), pub
}

func initTestPublicKeyMap() {
	argumentSet := []struct {
		jwtEnvKey        string
		payloadProductId string
		keyProductId     string
		et               int64
		prodAcceptable   bool
		additionalClaims map[string]interface{}
	}{
		{"SAO_EXPIRED_TOKEN", "sao", "sao", 1704134213, true, nil},
		{"SAO_NEVER_EXPIRED_TOKEN", "sao", "sao", 4600012345, true, nil},
		{"SDS_NEVER_EXPIRED_TOKEN", "sds", "sds", 4600012345, true, nil},
		{"SXX_NEVER_EXPIRED_TOKEN", "sxx", "sao", 4600012345, false, nil}, // sxx is not acceptable, the test case requires sign it by sao key
	}

	var productIdToKeyPair = make(map[string]KeyPair)

	for _, args := range argumentSet {
		key, pub := prepareProductIdKeyPair(productIdToKeyPair, args.keyProductId)
		tokenString, publicKeyPEMString, publicKey := prepareTestJWT(args.payloadProductId, args.et, key, pub, args.additionalClaims) // set Environment Variable
		log.Printf("Set JWT to environment variable: jwtEnvKey=%s, tokenString=%s", args.jwtEnvKey, tokenString)
		err := os.Setenv(args.jwtEnvKey, tokenString)
		if err != nil {
			log.Printf("Set ENV failed for jwtEnvKey %s: %v", args.jwtEnvKey, err)
			panic(err)
		}

		// construct map
		if args.prodAcceptable {
			publicKeyMap[args.payloadProductId] = &RSAPublicKey{
				PEM:       publicKeyPEMString,
				VerifyKey: &publicKey,
			}
		}
	}
}
