package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"

	"github.com/hankbao/reader-replica/scrape"
)

type MyEvent struct {
	Feed  scrape.Feed `json:"Feed"`
	Links []string    `json:"Links"`
}

type MyResponse struct {
	Feed     scrape.Feed       `json:"Feed"`
	Articles []*scrape.Article `json:"Articles"`
}

func HandleRequest(ctx context.Context, event *MyEvent) (*[]byte, error) {
	// Get the Lambda context
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		log.Panic("Failed to get Lambda context")
	}

	requestID := lc.AwsRequestID

	log.Printf("[%v] Processing message: %s\n", requestID, event.Feed.FeedLink)

	// Process the message
	data, err := fetchRSS(requestID, &event.Feed, event.Links)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// Fetch the RSS feed and return the processed data
func fetchRSS(requestID string, reqFeed *scrape.Feed, links []string) ([]byte, error) {
	scraper := scrape.NewScraper(30)
	feed, articles, err := scraper.ScrapeArticles(reqFeed, links)
	if err != nil {
		log.Printf("[%v] Failed to scrape articles: %v", requestID, err)
		return nil, err
	}

	log.Printf("[%v] Scraped %d articles from feed: %s", requestID, len(articles), reqFeed.FeedLink)

	// Serialize the processed data to JSON
	return json.Marshal(MyResponse{
		Feed:     *feed,
		Articles: articles,
	})
}

func main() {
	// Setup logging
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	// Start the Lambda handler
	lambda.Start(HandleRequest)
}
