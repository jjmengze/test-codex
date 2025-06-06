package kinesis

import (
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

type Kinesis interface {
	kinesisiface.KinesisAPI
}

type client struct {
	*kinesis.Kinesis
}

func NewClient(k *kinesis.Kinesis) Kinesis {
	return client{k}
}
