package crawler

import (
	"context"
	"time"
)

// Crawler defines settings for a crawl job.
type Crawler struct {
	StartURL string
	// Additional configuration (concurrency, delays, etc.) can be added here.
}

// NewCrawler creates a new crawler instance.
func NewCrawler(startURL string) *Crawler {
	return &Crawler{
		StartURL: startURL,
	}
}

// Crawl simulates a crawl job. In a real scenario, this method would run the crawler,
// collect metadata (titles, IP, headers, etc.), and return the results.
func (c *Crawler) Crawl(ctx context.Context) ([]string, error) {
	// Simulate work by waiting a few seconds.
	select {
	case <-time.After(5 * time.Second):
		// Simulated results.
		return []string{c.StartURL + "/page1", c.StartURL + "/page2"}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
