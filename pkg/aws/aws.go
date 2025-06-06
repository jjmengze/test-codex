package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/ssm"

	jwt "github.com/golang-jwt/jwt"
)

var (
	awsSession    *session.Session
	ssmClient     *ssm.SSM
	kinesisClient *kinesis.Kinesis
)

func init() {
	initClients()
	initPublicKeyMap()
}

var publicKeySSMParam = map[string]string{
	"sds": "/SAE/XDR/LogReceiver/Auth/SDS/PubKey", // Deep Security SaaS
	"sao": "/SAE/XDR/LogReceiver/Auth/SAO/PubKey", // Apex One SaaS
}

func NewKinesisClient() *kinesis.Kinesis {
	return kinesisClient
}

func initClients() {
	//only for test
	cfg := aws.NewConfig().WithRegion("us-west-2")
	awsSession = session.Must(session.NewSession(cfg))
	kinesisClient = kinesis.New(awsSession)
	ssmClient = ssm.New(awsSession)
}

func initPublicKeyMap() {
	for prodID, ssmParam := range publicKeySSMParam {
		publicKeyPEM, err := getSSMParameter(ssmParam, true)
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
}

func getSSMParameter(key string, withDecryption bool) (string, error) {
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
