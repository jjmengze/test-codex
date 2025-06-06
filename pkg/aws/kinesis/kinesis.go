package kinesis

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"log-receiver/pkg/logger"
)

type Kinesis interface {
	kinesisiface.KinesisAPI
}

type client struct {
	*kinesis.Kinesis
}

func NewClient(ctx context.Context, logger logger.Logger, k *kinesis.Kinesis) (Kinesis, error) {
	_, err := k.ListStreams(&kinesis.ListStreamsInput{
		Limit: aws.Int64(1),
	})
	if err != nil {
		err := fmt.Errorf("failed to verify Kinesis permissions: %w", err)
		logger.WithContext(ctx).ErrorF(err.Error())
		return nil, err
	}

	return client{k}, nil
}
