// parser.go
// author: hankbao

package scrape

import (
	"fmt"
	"io"

	"github.com/mmcdole/gofeed/rss"
)

type RSSParser struct {
	inner *rss.Parser
}

func NewRSSParser() *RSSParser {
	return &RSSParser{
		inner: &rss.Parser{},
	}
}

func (rp *RSSParser) Parse(r io.Reader) (*rss.Feed, error) {
	feed, err := rp.inner.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	return feed, nil
}
