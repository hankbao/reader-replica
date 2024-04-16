// main.go
// author: hankbao

package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"main/models"

	"github.com/hankbao/reader-replica/scrape"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"gorm.io/gorm"
)

func main() {
	log.Print("Taskmaster started")

	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_username := os.Getenv("DB_USERNAME")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := "reader"

	db, err := models.ConnectDatabase(db_host, db_port, db_username, db_password, db_name)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Set up AWS configuration and create an SQS client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	sqsClient := sqs.NewFromConfig(cfg)

	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		updateFeeds(db, sqsClient)
	}
}

func updateFeeds(db *gorm.DB, sqsClient *sqs.Client) {
	var feeds []scrape.Feed

	// Query to find feeds where lastBuildDate is older than 30 minutes ago
	result := db.Where("last_build_date < ?", time.Now().Add(-30*time.Minute)).Find(&feeds)
	if result.Error != nil {
		log.Printf("failed to query feeds: %v", result.Error)
		return
	}

	queueURL := "https://sqs.us-west-2.amazonaws.com/440824912727/rr-tasks"

	for _, feed := range feeds {
		msgBody, err := json.Marshal(feed)
		if err != nil {
			log.Printf("failed to marshal message body: %v", err)
			continue
		}

		input := &sqs.SendMessageInput{
			QueueUrl:    &queueURL,
			MessageBody: aws.String(string(msgBody)),
		}

		// Send message to SQS queue
		a, err := sqsClient.SendMessage(context.TODO(), input)
		if err != nil {
			log.Printf("failed to send message: %v", err)
			continue
		} else {
			log.Printf("message sent: %v", a)
		}
	}
}
