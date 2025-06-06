package main

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

const (
	streamName = "jason-orientation-test"
	region     = "us-west-2"
)

func main() {
	// 建立 session 與 client
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	client := kinesis.New(sess)

	// 列出 shards
	shardsOutput, err := client.DescribeStream(&kinesis.DescribeStreamInput{
		StreamName: aws.String(streamName),
	})
	if err != nil {
		log.Fatalf("DescribeStream error: %v", err)
	}

	for _, shard := range shardsOutput.StreamDescription.Shards {
		log.Printf("Reading from Shard: %s", *shard.ShardId)

		// 取得 shard iterator
		iterOutput, err := client.GetShardIterator(&kinesis.GetShardIteratorInput{
			StreamName:        aws.String(streamName),
			ShardId:           shard.ShardId,
			ShardIteratorType: aws.String("TRIM_HORIZON"), // 可用 LATEST 或 AT_TIMESTAMP
		})
		if err != nil {
			log.Fatalf("GetShardIterator error: %v", err)
		}

		shardIterator := iterOutput.ShardIterator

		for {
			out, err := client.GetRecords(&kinesis.GetRecordsInput{
				ShardIterator: shardIterator,
				Limit:         aws.Int64(10),
			})
			if err != nil {
				log.Fatalf("GetRecords error: %v", err)
			}

			for _, record := range out.Records {
				log.Printf("Record: PartitionKey=%v, Data=%v", *record.PartitionKey, string(record.Data))
			}

			if out.NextShardIterator == nil {
				log.Println("Shard closed")
				break
			}

			shardIterator = out.NextShardIterator
			time.Sleep(1 * time.Second)
		}
	}
}
