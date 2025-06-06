package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"

	"log-receiver/internal/domain/entity"
	repoMock "log-receiver/mock/internal_/repo"
	mockLogger "log-receiver/mock/pkg/logger"
	"log-receiver/pkg/protobuf"
)

func TestReceiverPutDataSuccess(t *testing.T) {
	ctx := context.Background()
	input := entity.PutDataInput{
		RawData:     []byte("data"),
		ProductCode: "sao",
		TraceID:     "trace",
		CustomerID:  "cust",
		Encoding:    "enc",
		SourceID:    "src",
	}

	mPub := repoMock.NewPublisher(t)
	mLogger := mockLogger.NewLogger(t)

	r := NewReceiver(mLogger, mPub)

	mPub.On("PutDataToStream", ctx, testifymock.AnythingOfType("[]uint8")).Run(func(args testifymock.Arguments) {
		b := args.Get(1).([]byte)
		var w protobuf.Wrapper
		err := proto.Unmarshal(b, &w)
		assert.NoError(t, err)
		assert.Equal(t, input.RawData, w.Payload)
		assert.Equal(t, input.ProductCode, aws.StringValue(w.ProductCode))
		assert.Equal(t, input.TraceID, aws.StringValue(w.TraceId))
		assert.Equal(t, input.CustomerID, aws.StringValue(w.CustomerId))
		assert.Equal(t, input.Encoding, aws.StringValue(w.Encoding))
		assert.Equal(t, input.SourceID, aws.StringValue(w.SourceId))
		assert.NotNil(t, w.ReceivedTime)
	}).Return(nil).Once()

	err := r.PutData(ctx, input)
	assert.NoError(t, err)
	mPub.AssertExpectations(t)
}

func TestReceiverPutDataPublisherError(t *testing.T) {
	ctx := context.Background()
	input := entity.PutDataInput{}
	retErr := errors.New("fail")

	mPub := repoMock.NewPublisher(t)
	mLogger := mockLogger.NewLogger(t)

	r := NewReceiver(mLogger, mPub)

	mPub.On("PutDataToStream", ctx, testifymock.Anything).Return(retErr).Once()
	mLogger.On("ErrorF", "PutDataToStream Failed: %v", retErr).Once()

	err := r.PutData(ctx, input)
	assert.ErrorIs(t, err, retErr)
	mPub.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}
