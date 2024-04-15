// fetcher.go
// author: hankbao

package scrape

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	retryable "github.com/hashicorp/go-retryablehttp"
)

type Fetcher struct {
	client *http.Client
}

func NetFetcher(timeout int) *Fetcher {
	// enable http retry
	client := retryable.NewClient()

	tp := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client.HTTPClient = &http.Client{Transport: tp}

	client.RetryMax = 3
	client.CheckRetry = retryable.ErrorPropagatedRetryPolicy

	cr := client.StandardClient()
	if timeout == 0 {
		timeout = 30 // default timeout 30s
	}
	cr.Timeout = time.Duration(timeout) * time.Second
	f := &Fetcher{
		client: cr,
	}

	return f
}

func (cr *Fetcher) Fetch(url string, lastModified string, eTag string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "rrscraper/0.1")
	if lastModified != "" {
		req.Header.Set("If-Modified-Since", lastModified)
	}
	if eTag != "" {
		req.Header.Set("If-None-Match", eTag)
	}

	resp, err := cr.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("response is nil without error")
	}

	// check if response is 304
	if resp.StatusCode == http.StatusNotModified {
		resp.Body.Close()
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	return resp, nil
}
