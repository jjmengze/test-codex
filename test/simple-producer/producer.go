package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

func main() {

	const (
		streamName = "jason-orientation-test"
		region     = "us-west-2"
	)

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create Kinesis client
	client := kinesis.NewFromConfig(cfg)

	partitionKey := fmt.Sprintf("partition-%d", time.Now().UnixNano())
	data := []byte("Hello, Kinesis from Go!")

	input := &kinesis.PutRecordInput{
		Data:         data,
		StreamName:   aws.String(streamName),
		PartitionKey: aws.String(partitionKey),
	}

	// Send record
	result, err := client.PutRecord(ctx, input)
	if err != nil {
		log.Fatalf("failed to put record: %v", err)
	}

	if result == nil {
		log.Fatalf("failed to put record, result is nil")
	}

	fmt.Printf("Successfully put record. ShardID: %s, SequenceNumber: %s\n",
		*result.ShardId, *result.SequenceNumber)
}
