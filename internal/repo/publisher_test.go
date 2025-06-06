package repo

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	awsKinesis "github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"

	mockKinesis "log-receiver/mock/pkg/aws/kinesis"
	mockLogger "log-receiver/mock/pkg/logger"
)

func TestPublisherPutDataToStreamSuccess(t *testing.T) {
	ctx := context.Background()
	data := []byte("hello")

	mKinesis := mockKinesis.NewKinesis(t)
	mLogger := mockLogger.NewLogger(t)

	p := NewPublisher(mLogger, mKinesis)

	mKinesis.On("PutRecord", testifymock.MatchedBy(func(input *awsKinesis.PutRecordInput) bool {
		return string(input.Data) == string(data) &&
			aws.StringValue(input.StreamName) == "jason-orientation-test" &&
			aws.StringValue(input.PartitionKey) == "test"
	})).Return(&awsKinesis.PutRecordOutput{SequenceNumber: aws.String("1")}, nil).Once()

	mLogger.On("InfoF", testifymock.Anything, testifymock.Anything).Once()

	err := p.PutDataToStream(ctx, data)
	assert.NoError(t, err)
	mKinesis.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestPublisherPutDataToStreamError(t *testing.T) {
	ctx := context.Background()
	data := []byte("hello")
	retErr := errors.New("fail")

	mKinesis := mockKinesis.NewKinesis(t)
	mLogger := mockLogger.NewLogger(t)

	p := NewPublisher(mLogger, mKinesis)

	mKinesis.On("PutRecord", testifymock.Anything).Return((*awsKinesis.PutRecordOutput)(nil), retErr).Once()
	mLogger.On("ErrorF", "failed to put record to Kinesis: %w", retErr).Once()

	err := p.PutDataToStream(ctx, data)
	assert.ErrorIs(t, err, retErr)
	mKinesis.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}
