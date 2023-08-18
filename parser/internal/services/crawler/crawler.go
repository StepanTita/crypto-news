package crawler

import (
	"context"

	"common/data/model"
)

type ParsedBody interface {
	ToModel() any
}

type Crawler interface {
	Crawl(ctx context.Context) ([]ParsedBody, int, error)
}

type MultiCrawler[T model.Model] interface {
	Crawl(ctx context.Context, entities []T) ([]ParsedBody, []int, []error)
}

type CrawlFunc func(ctx context.Context) ([]ParsedBody, int, error)
