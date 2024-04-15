package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/config"
)

type MyEvent struct {
	FeedInfo struct {
		FeedLink string `json:"feedLink"`
	} `json:"FeedInfo"`
}

func HandleRequest(ctx context.Context, event *MyEvent) (*string, error) {
	_, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Panicf("configuration error, " + err.Error())
	}

	// Get the Lambda context
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		log.Panic("failed to get Lambda context")
	}

	requestID := lc.AwsRequestID

	fmt.Printf("[%v] Processing message: %s\n", requestID, event.FeedInfo.FeedLink)

	// Process the message
	data, err := fetchRSS(event.FeedInfo.FeedLink)
	if err != nil {
		ret := ""
		return &ret, err
	}

	fmt.Printf("[%v] Processed data: %s\n", requestID, string(data))

	// // The URL of the second SQS queue
	// secondQueueURL := "https://sqs.<region>.amazonaws.com/<account_id>/<queue_name>"

	// // Send the message to the second SQS queue
	// sqsClient := sqs.NewFromConfig(cfg)
	// _, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
	// 	QueueUrl:    aws.String(secondQueueURL),
	// 	MessageBody: aws.String(string(data)),
	// })

	// if err != nil {
	// 	return "", err
	// }

	ret := "Messages processed and sent to the second queue successfully!"
	return &ret, nil
}

func fetchRSS(link string) ([]byte, error) {
	// Your message processing logic here
	// This is just a dummy example that converts the message to uppercase
	data := map[string]string{
		"processed_data": strings.ToUpper(link),
	}

	// Serialize the processed data to JSON
	return json.Marshal(data)
}

func main() {
	// Setup logging
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	// Start the Lambda handler
	lambda.Start(HandleRequest)
}
