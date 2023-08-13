package crawler

import (
	"context"
)

type ParsedBody interface {
	ToModel() any
}

type Crawler interface {
	Crawl(ctx context.Context) (any, int, error)
}

type CrawlFunc func(ctx context.Context) (any, int, error)
