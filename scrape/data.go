// data.go
// author: hankbao

package scrape

import (
	"time"

	"github.com/mmcdole/gofeed/rss"
)

type Feed struct {
	Id            string `json:"id" gorm:"primary_key;column:id"`
	CreatedAt     int    `json:"createdAt,omitempty" gorm:"column:createdAt"`
	UpdatedAt     int    `json:"updatedAt,omitempty" gorm:"column:updatedAt"`
	Title         string `json:"title,omitempty" gorm:"column:title"`
	Description   string `json:"description,omitempty" gorm:"column:description"`
	Language      string `json:"language,omitempty" gorm:"column:language"`
	Link          string `json:"link,omitempty" gorm:"column:link"`
	FeedLink      string `json:"feedLink,omitempty" gorm:"column:feedLink"`
	LastBuildDate string `json:"lastBuildDate,omitempty" gorm:"column:lastBuildDate"`
	LastBuildAt   int    `json:"lastBuildAt,omitempty" gorm:"column:lastBuildAt"`
	LastFetchAt   int    `json:"lastFetchAt,omitempty" gorm:"column:lastFetchAt"`
	LastModified  string `json:"lastModified,omitempty" gorm:"column:lastModified"`
	ETag          string `json:"eTag,omitempty" gorm:"column:eTag"`
}

func (Feed) TableName() string {
	return "feeds"
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
	Id          string    `json:"id" gorm:"primary_key;column:id"`
	CreatedAt   int       `json:"createdAt,omitempty" gorm:"column:createdAt"`
	UpdatedAt   int       `json:"updatedAt,omitempty" gorm:"column:updatedAt"`
	FeedID      string    `json:"feedId" gorm:"column:feedId"`
	Title       string    `json:"title,omitempty" gorm:"column:title"`
	Link        string    `json:"link,omitempty" gorm:"column:link"`
	Description string    `json:"description,omitempty" gorm:"column:description"`
	Content     string    `json:"content,omitempty" gorm:"column:content"`
	Author      string    `json:"author,omitempty" gorm:"column:author"`
	Categories  []string  `json:"categories,omitempty" gorm:"column:categories"`
	Comments    string    `json:"comments,omitempty" gorm:"column:comments"`
	Enclosure   Enclosure `json:"enclosure,omitempty" gorm:"column:enclosure"`
	GUID        string    `json:"guid,omitempty" gorm:"column:guid"`
	PubDate     string    `json:"pubDate,omitempty" gorm:"column:pubDate"`
	PubAt       int       `json:"pubAt,omitempty" gorm:"column:pubAt"`
	Source      Source    `json:"source,omitempty" gorm:"column:source"`
}

func (Article) TableName() string {
	return "articles"
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
