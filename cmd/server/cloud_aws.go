package main

//type cloudAWS struct{}
//
//func cloudImpl() CloudProvider {
//	return &cloudAWS{}
//}

// TODO: Implement this function to put data into a Kinesis stream
// This function should:
// 1. Create a PutRecord request for the Kinesis stream
// 2. Set the stream name, partition key, and data
// 3. Send the request to the Kinesis service
// 4. Handle any errors and return them
//func (cloudAWS) PutDataToStream(data []byte) error {
//
//	input := &kinesis.PutRecordInput{
//		Data:         data,
//		StreamName:   aws.String("jason-orientation-test"),
//		PartitionKey: aws.String("test"), // 可用 userID、UUID 等作為 key
//	}
//	output, err := kinesisClient.PutRecord(input)
//	if err != nil {
//		err = fmt.Errorf("failed to put record to Kinesis: %w", err)
//		log.Printf(err.Error())
//		return err
//	}
//
//	For now, just log that we would send data and return nil
//log.Printf("put record to Kinesis succeeded %v", *output)
//return nil
//}
