// data.go
// author: hankbao

package scrape

import (
	"time"

	"github.com/mmcdole/gofeed/rss"
)

type Feed struct {
	Id            string `json:"_id"`
	CreatedAt     int    `json:"createdAt,omitempty"`
	UpdatedAt     int    `json:"updatedAt,omitempty"`
	Title         string `json:"title,omitempty"`
	Description   string `json:"description,omitempty"`
	Language      string `json:"language,omitempty"`
	Link          string `json:"link,omitempty"`
	FeedLink      string `json:"feedLink,omitempty"`
	LastBuildDate string `json:"lastBuildDate,omitempty"`
	LastBuildAt   int    `json:"lastBuildAt,omitempty"`
	LastFetchAt   int    `json:"lastFetchAt,omitempty"`
	LastModified  string `json:"lastModified,omitempty"`
	ETag          string `json:"eTag,omitempty"`
}

func (feed *Feed) UpdateFrom(rssFeed *rss.Feed, lastModified string, eTag string) {
	feed.Title = rssFeed.Title
	feed.Description = rssFeed.Description
	feed.Language = rssFeed.Language
	feed.Link = rssFeed.Link
	feed.LastBuildDate = rssFeed.LastBuildDate

	if rssFeed.LastBuildDateParsed == nil {
		feed.LastBuildAt = 0
	} else {
		feed.LastBuildAt = int(rssFeed.LastBuildDateParsed.Unix() * 1000) // as timestamp in js
	}

	feed.LastFetchAt = int(time.Now().Unix() * 1000)
	feed.LastModified = lastModified
	feed.ETag = eTag
}

type Article struct {
	Id          string    `json:"_id"`
	CreatedAt   int       `json:"createdAt,omitempty"`
	UpdatedAt   int       `json:"updatedAt,omitempty"`
	FeedID      string    `json:"feedId"`
	Title       string    `json:"title,omitempty"`
	Link        string    `json:"link,omitempty"`
	Description string    `json:"description,omitempty"`
	Content     string    `json:"content,omitempty"`
	Author      string    `json:"author,omitempty"`
	Categories  []string  `json:"categories,omitempty"`
	Comments    string    `json:"comments,omitempty"`
	Enclosure   Enclosure `json:"enclosure,omitempty"`
	GUID        string    `json:"guid,omitempty"`
	PubDate     string    `json:"pubDate,omitempty"`
	PubAt       int       `json:"pubAt,omitempty"`
	Source      Source    `json:"source,omitempty"`
}

func NewArticleFrom(item *rss.Item, feedID string) *Article {
	article := &Article{
		FeedID:      feedID,
		Title:       item.Title,
		Link:        item.Link,
		Description: item.Description,
		Content:     item.Content,
		Author:      item.Author,
		Comments:    item.Comments,
		PubDate:     item.PubDate,
	}

	if item.GUID != nil {
		article.GUID = item.GUID.Value
	}

	if item.Source != nil {
		article.Source = Source{Title: item.Source.Title, URL: item.Source.URL}
	}

	if item.Enclosure != nil {
		article.Enclosure = Enclosure{URL: item.Enclosure.URL, Length: item.Enclosure.Length, Type: item.Enclosure.Type}
	}

	if item.PubDateParsed == nil {
		article.PubAt = 0
	} else {
		article.PubAt = int(item.PubDateParsed.Unix() * 1000) // as timestamp in js
	}

	var categories []string
	for _, c := range item.Categories {
		categories = append(categories, c.Value)
	}
	article.Categories = categories

	return article
}

type Enclosure struct {
	URL    string `json:"url,omitempty"`
	Length string `json:"length,omitempty"`
	Type   string `json:"type,omitempty"`
}

type Source struct {
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
}
