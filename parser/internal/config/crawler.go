package config

import "time"

type Crawler interface {
	// TODO: think of making this service dependent
	RateLimit() int
	CrawlEvery() time.Duration
}

type crawler struct {
	rateLimit  int
	crawlEvery time.Duration
}

func NewCrawler(rateLimit int, crawlEvery time.Duration) Crawler {
	return &crawler{
		rateLimit:  rateLimit,
		crawlEvery: crawlEvery,
	}
}

func (l *crawler) CrawlEvery() time.Duration {
	return l.crawlEvery
}

func (l *crawler) RateLimit() int {
	return l.rateLimit
}
