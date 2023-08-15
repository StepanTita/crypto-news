package url_crawler

import (
	"context"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common/convert"
	"common/data/model"
	"common/data/store"
	"parser/internal/config"
	"parser/internal/services/connector"
	"parser/internal/services/crawler"
)

const (
	maxLimit = 100
)

type UrlCrawler struct {
	log *logrus.Entry

	conn connector.Connector

	dataProvider store.DataProvider
}

func NewCrawler(cfg config.Config) crawler.Crawler {
	return UrlCrawler{
		log:          cfg.Logging().WithField("service", "[URL-CRAWLER]"),
		conn:         connector.New(cfg),
		dataProvider: store.New(cfg),
	}
}

func (u UrlCrawler) Crawl(ctx context.Context) (any, int, error) {
	// TODO: process this in batches to reduce RAM load
	pendingTitles, err := u.dataProvider.TitlesProvider().ByStatus(model.StatusPending).Select(ctx)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to select pending titles")
	}

	outBodies := make([]crawler.ParsedBody, 0, 10)

	for _, t := range pendingTitles {
		respBody, statusCode, err := u.conn.Poll(ctx, connector.PollParams{
			Url: convert.FromPtr(t.URL),
		})

		if err != nil {
			return nil, 0, errors.Wrapf(err, "failed to poll url %s", convert.FromPtr(t.URL))
		}

		if statusCode != http.StatusOK {
			return nil, statusCode, nil
		}

		rawResp, err := io.ReadAll(respBody)
		if err != nil {
			return nil, statusCode, errors.Wrap(err, "failed to read response body")
		}

		outBodies = append(outBodies, body{text: string(rawResp)})
	}

	return outBodies, http.StatusOK, nil
}
