package crypto_panic_crawler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common/convert"
	"common/data"
	"common/data/model"
	"common/data/store"
	"parser/internal/config"
	"parser/internal/services/connector"
	"parser/internal/services/crawler"
	"parser/internal/utils"
)

type CryptoPanicCrawler struct {
	log *logrus.Entry

	authToken string
	url       string
	path      string
	conn      connector.Connector

	dataProvider store.DataProvider
}

func NewCrawler(cfg config.Config) crawler.Crawler {
	return &CryptoPanicCrawler{
		log: cfg.Logging().WithField("service", "[CRYPTO-PANIC-CRAWLER]"),

		authToken: cfg.Credentials(utils.CryptoPanic)["auth_token"],
		url:       cfg.Credentials(utils.CryptoPanic)["url"],
		path:      cfg.Credentials(utils.CryptoPanic)["path"],

		dataProvider: store.New(cfg),

		conn: connector.New(cfg),
	}
}

func (c CryptoPanicCrawler) Crawl(ctx context.Context) ([]crawler.ParsedBody, int, error) {
	latestNews, err := c.dataProvider.NewsProvider().BySources(utils.CryptoPanic).GetLatest(ctx)
	if err != nil {
		if !errors.Is(err, data.ErrNotFound) {
			return nil, 0, errors.Wrap(err, "failed to query news provider by source")
		}
		// if there is no previous record, just set prev published time as 0
		latestNews = &model.News{PublishedAt: &time.Time{}}
	}

	body, statusCode, err := c.conn.Request(ctx, connector.RequestParams{
		Url:  c.url,
		Path: c.path,
		Params: url.Values{
			"auth_token": []string{c.authToken},
			"kind":       []string{"news"},
			"public":     []string{"true"},
			"metadata":   []string{"true"},
			//"approved": []string{"true"},
		},
	})
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to request crypto-panic API")
	}

	if statusCode != http.StatusOK {
		return nil, statusCode, nil
	}

	var b map[string]json.RawMessage
	if err := json.NewDecoder(body).Decode(&b); err != nil {
		return nil, statusCode, errors.Wrap(err, "failed to decode response body")
	}

	var out []Body
	if err = json.Unmarshal(b["results"], &out); err != nil {
		return nil, statusCode, errors.Wrap(err, "failed to decode response body results")
	}

	return utils.Map(utils.Filter(out, filterOld(convert.FromPtr(latestNews.PublishedAt))), toParsedBody), statusCode, nil
}

func filterOld(latestPrev time.Time) func(b Body) bool {
	return func(b Body) bool {
		return b.PublishedAt.After(latestPrev)
	}
}

func toParsedBody(b Body) crawler.ParsedBody {
	return b
}
