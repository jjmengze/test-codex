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

// TODO: Implement this function to properly wrap the input data in a Wrapper struct
// This function should:
// 1. Create and initialize a Wrapper struct with data from the input
// 2. Set appropriate fields in the Wrapper (productCode, traceID, customerID, etc.)
// 3. Set the payload field with the raw data
// 4. Marshal the Wrapper to protobuf format
// 5. Send the data to Kinesis stream
//func (p *publisher) PutDataToStream(ctx context.Context, data []byte) error {
//
//	// This is a minimal implementation that should be expanded by new employees
//
//	// Create a basic wrapper - new employees should properly initialize this
//	dataWrap := &protobuf.Wrapper{
//		Payload:     input.rawData,
//		ProductCode: aws.String(input.productCode),
//		TraceId:     aws.String(input.traceID),
//		CustomerId:  aws.String(input.customerID),
//		Encoding:    aws.String(input.encoding),
//		SourceId:    aws.String(input.sourceID),
//		DataSchema:  protobuf.Wrapper_COMMON.Enum(), // 若你有其他 schema，也可改這邊
//		ReceivedTime: func() *uint64 {
//			ms := uint64(time.Now().UnixMilli())
//			return &ms
//		}(),
//	}
//
//	// Marshal the wrapper to protobuf format
//	data, err := proto.Marshal(dataWrap)
//	if err != nil {
//		log.Printf("Encode Protobuf Failed: %v", err)
//		return err
//	}
//
//	// Send the data to Kinesis stream
//	err = cloud.PutDataToStream(data)
//
//	return err
//}
