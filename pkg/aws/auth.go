package aws

import (
	"crypto/rsa"
)

type NonSupportAuthProductCodeError struct {
}

type NonSupportFeatureError struct {
}

func (e NonSupportAuthProductCodeError) Error() string {
	return "Non Support ProductCode"
}

func (e NonSupportFeatureError) Error() string {
	return "Non Support Feature"
}

var publicKeyMap = make(map[string]*RSAPublicKey)

type RSAPublicKey struct {
	PEM       string
	VerifyKey *rsa.PublicKey
}

func GetVerifyKeyByProductCode(productCode string) (*rsa.PublicKey, error) {
	if rsaPublicKey, ok := publicKeyMap[productCode]; ok {
		return rsaPublicKey.VerifyKey, nil
	}
	return nil, NonSupportAuthProductCodeError{}
}
