// scraper.go
// author: hankbao

package scrape

import (
	"log"

	"github.com/mmcdole/gofeed/rss"
)

type Scraper struct {
	fetcher *Fetcher
	parser  *RSSParser
}

func NewScraper(timeout int) *Scraper {
	return &Scraper{
		fetcher: NetFetcher(timeout),
		parser:  NewRSSParser(),
	}
}

func (sc *Scraper) ScrapeFeed(link string) (*Feed, error) {
	log.Printf("start to fetch feed: %s", link)

	resp, err := sc.fetcher.Fetch(link, "", "")
	if err != nil {
		log.Printf("failed to fetch feed: %s as %v", link, err)
		return nil, err
	}

	if resp == nil {
		log.Printf("feed not modified since last fetch: %s", link)
		return nil, nil
	}

	log.Printf("feed fetched: %s", link)
	defer func() {
		ce := resp.Body.Close()
		if ce != nil {
			log.Printf("failed to close response body: %v", ce)
		}
	}()

	feed, err := sc.parser.Parse(resp.Body)
	if err != nil {
		log.Printf("failed to parse feed: %s as %v", link, err)
		return nil, err
	}

	log.Printf("feed response parsed: %s", link)

	var f = &Feed{}
	lastModified := resp.Header.Get("Last-Modified")
	eTag := resp.Header.Get("ETag")
	f.UpdateFrom(feed, lastModified, eTag)

	return f, nil
}

func (sc *Scraper) ScrapeArticles(reqFeed *Feed, linksFetched []string) (*Feed, []*Article, error) {
	log.Printf("start to fetch feed: %s", reqFeed.FeedLink)

	resp, err := sc.fetcher.Fetch(reqFeed.FeedLink, reqFeed.LastModified, reqFeed.ETag)
	if err != nil {
		log.Printf("failed to fetch feed: %v as %v", reqFeed, err)
		return nil, nil, err
	}

	if resp == nil {
		log.Printf("feed not modified since last fetch: %v", reqFeed)
		return nil, nil, nil
	}

	log.Printf("feed fetched: %v", reqFeed)
	defer func() {
		ce := resp.Body.Close()
		if ce != nil {
			log.Printf("failed to close response body: %v", ce)
		}
	}()

	feed, err := sc.parser.Parse(resp.Body)
	if err != nil {
		log.Printf("failed to parse feed: %v as %v", reqFeed, err)
		return nil, nil, err
	}

	log.Printf("feed response parsed: %v", reqFeed)

	// Extract new articles from feed
	articles := sc.extractArticles(feed, reqFeed, linksFetched)

	// Update feed metadata
	lastModified := resp.Header.Get("Last-Modified")
	eTag := resp.Header.Get("ETag")
	reqFeed.UpdateFrom(feed, lastModified, eTag)

	return reqFeed, articles, nil
}

func (sc *Scraper) extractArticles(feed *rss.Feed, reqFeed *Feed, linksFetched []string) []*Article {
	// we don't have a Set data structure in Go, so we use a map instead
	linkMap := make(map[string]string, len(linksFetched))
	for _, link := range linksFetched {
		linkMap[link] = link
	}

	articles := make([]*Article, 0, len(feed.Items))
	for _, item := range feed.Items {
		if item.Link == "" {
			log.Printf("fetched item got no link: %v", item)
			continue
		}

		if _, ok := linkMap[item.Link]; ok {
			continue
		}

		article := NewArticleFrom(item, reqFeed.Id)
		articles = append(articles, article)
	}

	return articles
}
