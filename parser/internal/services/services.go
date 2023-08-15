package services

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common"
	"common/convert"
	"common/data"
	"common/data/model"
	"common/data/store"
	"parser/internal/config"
	browse_ai_crawler "parser/internal/services/browse-ai-crawler"
	"parser/internal/services/crawler"
	url_crawler "parser/internal/services/url-crawler"
	"parser/internal/services/worker"
)

type Service interface {
	Run(ctx context.Context) error
}

const workersNum = 1

type service struct {
	cfg config.Config
	log *logrus.Entry

	titlesCrawlers []crawler.Crawler
	newsCrawler    crawler.Crawler

	dataProvider store.DataProvider
}

func NewService(cfg config.Config) Service {
	return &service{
		cfg: cfg,
		log: cfg.Logging().WithField("service", "[TITLES-PARSER]"),
		titlesCrawlers: []crawler.Crawler{
			browse_ai_crawler.NewCrawler(cfg, cfg.Credentials(browse_ai_crawler.BrowseAI, "robots", "coin_telegraph")),
		},
		newsCrawler:  url_crawler.NewCrawler(cfg),
		dataProvider: store.New(cfg),
	}
}

func (s *service) Run(ctx context.Context) error {
	s.log.Infof("Staring crawling every %v...", s.cfg.CrawlEvery())
	go func() {
		err := common.RunEvery(s.cfg.CrawlEvery(), func() error {
			s.log.Debugf("Crawling %d...", len(s.titlesCrawlers))

			wrk := worker.New(workersNum, s.cfg)

			wrk.Produce(ctx, s.titlesCrawlers)
			for _, t := range wrk.Work(ctx) {
				if t.Err != nil {
					return errors.Wrap(t.Err, "failed to crawl")
				}

				if t.StatusCode != http.StatusOK {
					s.log.WithFields(logrus.Fields{
						"status-code": t.StatusCode,
						"info":        t.StatusBody,
					}).Warn("request returned unsuccessful status code...")
					continue
				}

				bodiesBatch := t.Body

				if len(bodiesBatch) == 0 {
					s.log.Debug("early stopping, no new titles...")
					continue
				}

				titlesBatch := crawler.ToModelBatch[model.Title](bodiesBatch)
				s.log.Debugf("Adding new batch to the database: %d", len(bodiesBatch))
				err := s.dataProvider.TitlesProvider().InsertUniqueBatch(ctx, titlesBatch)
				if err != nil {
					return errors.Wrap(err, "failed to insert batch of titles")
				}
			}
			return nil
		})

		if err != nil {

		}
	}()

	return common.RunEvery(s.cfg.CrawlEvery()+15*time.Second, func() error {
		body, statusCode, err := s.newsCrawler.Crawl(ctx)
		if err != nil {
			if !errors.Is(err, data.ErrNotFound) {
				return errors.Wrap(err, "failed to crawl url")
			}
			return nil
		}

		if statusCode != http.StatusOK {
			s.log.WithFields(logrus.Fields{
				"status-code": statusCode,
				"info":        body,
			}).Warn("request returned unsuccessful status code...")
			return nil
		}

		bodiesBatch, ok := body.([]crawler.ParsedBody)
		if !ok {
			return errors.New("could not cast output to parsed body")
		}

		if len(bodiesBatch) == 0 {
			s.log.Debug("early stopping, no new titles...")
			return nil
		}

		rawNewsWebpagesBatch := crawler.ToModelBatch[model.RawNewsWebpage](bodiesBatch)
		s.log.Debugf("Adding new batch to the database: %d", len(bodiesBatch))
		err = s.dataProvider.RawNewsWebpagesProvider().InsertBatch(ctx, rawNewsWebpagesBatch)
		if err != nil {
			return errors.Wrap(err, "failed to insert batch of titles")
		}

		_, err = s.dataProvider.TitlesProvider().Update(ctx, model.UpdateTitleParams{
			Status: convert.ToPtr(model.StatusProcessed),
		})
		if err != nil {
			return errors.Wrap(err, "failed to update titles status to processed")
		}

		return nil
	})
}
