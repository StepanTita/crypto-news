package services

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common"
	"common/data"
	"common/data/model"
	"common/data/store"
	"crypto-news/internal/config"
	"crypto-news/internal/services/crawler"
	crypto_panic_crawler "crypto-news/internal/services/crypto-panic-crawler"
	"crypto-news/internal/services/worker"
)

// TODO: add more when new sources added
const workersNum = 1

type Service interface {
	Run(ctx context.Context) error
}

type service struct {
	cfg config.Config
	log *logrus.Entry

	crawlers []crawler.Crawler

	dataProvider store.DataProvider
}

func NewService(cfg config.Config) Service {
	return &service{
		cfg:          cfg,
		log:          cfg.Logging().WithField("service", "[PARSER]"),
		crawlers:     []crawler.Crawler{crypto_panic_crawler.NewCrawler(cfg)},
		dataProvider: store.New(cfg),
	}
}

func (s *service) Run(ctx context.Context) error {
	s.log.Infof("Staring crawling every %v...", s.cfg.CrawlEvery())
	return common.RunEvery(s.cfg.CrawlEvery(), func() error {
		s.log.Debugf("Crawling %d...", len(s.crawlers))

		channels, err := s.dataProvider.ChannelsProvider().Select(ctx)
		if err != nil {
			if !errors.Is(err, data.ErrNotFound) {
				return errors.Wrap(err, "failed to select channels")
			}
		}

		wrk := worker.New(workersNum, s.cfg)

		wrk.Produce(ctx, s.crawlers)
		for _, t := range wrk.Work(ctx) {
			if t.Err != nil {
				return errors.Wrap(t.Err, "failed to crawl")
			}

			if t.StatusCode != http.StatusOK {
				s.log.WithField("status-code", t.StatusCode).Warn("request returned unsuccessful status code...")
				continue
			}

			bodiesBatch := t.Body

			if len(bodiesBatch) == 0 {
				s.log.Debug("early stopping, no new posts...")
				continue
			}

			newsBatch := toNewsBatch(bodiesBatch)
			s.log.Debugf("Adding new batch to the database: %d", len(bodiesBatch))
			err := s.dataProvider.NewsProvider().InsertBatch(ctx, newsBatch)
			if err != nil {
				return errors.Wrap(err, "failed to insert batch of news")
			}

			newsChannelsBatch := toNewsChannelsBatch(newsBatch, channels)
			if err = s.dataProvider.NewsChannelsProvider().InsertBatch(ctx, newsChannelsBatch); err != nil {
				return errors.Wrap(err, "failed to insert batch of news-channels")
			}

			coinsBatch, newsCoinsBatch := splitCoinsBatch(newsBatch)
			if err = s.dataProvider.CoinsProvider().UpsertCoinsBatch(ctx, coinsBatch); err != nil {
				return errors.Wrap(err, "failed to insert batch of coins")
			}

			if err = s.dataProvider.NewsCoinsProvider().InsertBatch(ctx, newsCoinsBatch); err != nil {
				return errors.Wrap(err, "failed to insert batch of news")
			}
		}
		return nil
	})
}

func toNewsBatch(bodies []crawler.ParsedBody) []model.News {
	newsBatch := make([]model.News, len(bodies))
	for i := range bodies {
		newsBatch[i] = bodies[i].ToNews()
	}
	return newsBatch
}

func splitCoinsBatch(newsBatch []model.News) ([]model.Coin, []model.NewsCoin) {
	uniqueCoinsMap := make(map[string]model.Coin)
	newsCoins := make([]model.NewsCoin, 0, 10)
	for _, n := range newsBatch {
		for _, c := range n.Coins {
			uniqueCoinsMap[c.Code] = c
			newsCoins = append(newsCoins, model.NewsCoin{
				Code:   c.Code,
				NewsID: n.ID,
			})
		}
	}

	uniqueCoins := make([]model.Coin, 0, len(uniqueCoinsMap))
	for _, coin := range uniqueCoinsMap {
		uniqueCoins = append(uniqueCoins, coin)
	}
	return uniqueCoins, newsCoins
}

func toNewsChannelsBatch(newsBatch []model.News, channels []model.Channel) []model.NewsChannel {
	newsChannels := make([]model.NewsChannel, len(newsBatch)*len(channels))
	i := 0
	for _, n := range newsBatch {
		for _, c := range channels {
			newsChannels[i] = model.NewsChannel{
				ChannelID: c.ChannelID,
				NewsID:    n.ID,
			}
			i++
		}
	}
	return newsChannels
}
