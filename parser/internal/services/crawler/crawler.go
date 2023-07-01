package crawler

import (
	"context"

	"common/data/model"
)

type ParsedBody interface {
	ToNews() model.News
}

type Crawler interface {
	Crawl(ctx context.Context) (any, int, error)
}

type CrawlFunc func(ctx context.Context) (any, int, error)
