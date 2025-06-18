package aws

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"

	jwt "github.com/golang-jwt/jwt"
)

var publicKeySSMParam = map[string]string{
	"sds": "/SAE/XDR/LogReceiver/Auth/SDS/PubKey", // Deep Security SaaS
	"sao": "/SAE/XDR/LogReceiver/Auth/SAO/PubKey", // Apex One SaaS
}
var publicKeyMapInitOnce sync.Once

func NewKinesisClient(awsSession *session.Session) *kinesis.Kinesis {
	return kinesis.New(awsSession)
}

func NewSsmClient(awsSession *session.Session) *ssm.SSM {
	return ssm.New(awsSession)
}
func NewSession(cfg *aws.Config) (*session.Session, error) {
	return session.NewSession(cfg)
}

func InitPublicKeyMap(ssmClient ssmiface.SSMAPI) {
	publicKeyMapInitOnce.Do(func() {
		for prodID, ssmParam := range publicKeySSMParam {
			publicKeyPEM, err := getSSMParameter(ssmClient, ssmParam, true)
			if err != nil {
				log.Printf("Get Public Key from SSM failed: %v", err)
				panic(err)
			}

			verifyKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPEM))
			if err != nil {
				log.Printf("Parse Public Key with RSA Public key Form failed: %v", err)
				panic(err)
			}

			publicKeyMap[prodID] = &RSAPublicKey{
				PEM:       publicKeyPEM,
				VerifyKey: verifyKey,
			}
		}
	})
}

func getSSMParameter(ssmClient ssmiface.SSMAPI, key string, withDecryption bool) (string, error) {
	param, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(withDecryption),
	})
	if err != nil {
		log.Printf("Get SSM Failed: %v", err)
		return "", err
	}
	value := aws.StringValue(param.Parameter.Value)
	return value, err
}
