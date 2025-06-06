package repo

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	awsKinesis "github.com/aws/aws-sdk-go/service/kinesis"

	"log-receiver/pkg/aws/kinesis"
	"log-receiver/pkg/logger"
)

type Publisher interface {
	PutDataToStream(ctx context.Context, data []byte) error
}

var _ Publisher = &publisher{}

type publisher struct {
	client kinesis.Kinesis
	logger logger.Logger
}

func NewPublisher(logger logger.Logger, client kinesis.Kinesis) Publisher {
	return publisher{
		client: client,
		logger: logger,
	}
}

// TODO: Implement this function to put data into a Kinesis stream
// This function should:
// 1. Create a PutRecord request for the Kinesis stream
// 2. Set the stream name, partition key, and data
// 3. Send the request to the Kinesis service
// 4. Handle any errors and return them
func (p publisher) PutDataToStream(ctx context.Context, data []byte) error {

	input := &awsKinesis.PutRecordInput{
		Data:         data,
		StreamName:   aws.String("jason-orientation-test"),
		PartitionKey: aws.String("test"), // 可用 userID、UUID 等作為 key
	}
	output, err := p.client.PutRecord(input)
	if err != nil {
		p.logger.ErrorF("failed to put record to Kinesis: %w", err)
		return err
	}

	// For now, just log that we would send data and return nil
	p.logger.InfoF("put record to Kinesis succeeded %v", *output)
	return nil
}
