package usecase

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"google.golang.org/protobuf/proto"
	"log-receiver/internal/domain/entity"
	"log-receiver/internal/repo"
	"log-receiver/pkg/logger"
	"log-receiver/pkg/protobuf"
)

type Receiver interface {
	PutData(ctx context.Context, input entity.PutDataInput) error
}

type receiver struct {
	logger    logger.Logger
	publisher repo.Publisher
}

func NewReceiver(logger logger.Logger, publisher repo.Publisher) Receiver {
	return receiver{logger: logger, publisher: publisher}
}

// TODO: Implement this function to properly wrap the input data in a Wrapper struct
// This function should:
// 1. Create and initialize a Wrapper struct with data from the input
// 2. Set appropriate fields in the Wrapper (productCode, traceID, customerID, etc.)
// 3. Set the payload field with the raw data
// 4. Marshal the Wrapper to protobuf format
// 5. Send the data to Kinesis stream
func (r receiver) PutData(ctx context.Context, input entity.PutDataInput) error {

	// This is a minimal implementation that should be expanded by new employees
	// Create a basic wrapper - new employees should properly initialize this
	dataWrap := &protobuf.Wrapper{
		Payload:     input.RawData,
		ProductCode: aws.String(input.ProductCode),
		TraceId:     aws.String(input.TraceID),
		CustomerId:  aws.String(input.CustomerID),
		Encoding:    aws.String(input.Encoding),
		SourceId:    aws.String(input.SourceID),
		DataSchema:  protobuf.Wrapper_COMMON.Enum(), // 若你有其他 schema，也可改這邊
		ReceivedTime: func() *uint64 {
			ms := uint64(time.Now().UnixMilli())
			return &ms
		}(),
	}

	// Marshal the wrapper to protobuf format
	data, err := proto.Marshal(dataWrap)
	if err != nil {
		r.logger.ErrorF("Encode Protobuf Failed: %v", err)
		return err
	}

	// Send the data to Kinesis stream
	err = r.publisher.PutDataToStream(ctx, data)
	if err != nil {
		r.logger.ErrorF("PutDataToStream Failed: %v", err)
		return err
	}

	return err
}
