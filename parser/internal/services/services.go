package services

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common"
	"common/containers/set"
	"common/convert"
	"common/data"
	"common/data/model"
	"common/data/store"
	"common/iteration"
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
	newsCrawler    crawler.MultiCrawler[model.Title]

	dataProvider store.DataProvider
}

func NewService(cfg config.Config) Service {
	return &service{
		cfg: cfg,
		log: cfg.Logging().WithField("service", "[PARSER]"),
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
			s.log.WithError(err).Error("failed to run titles crawlers")
		}

		return
	}()

	return common.RunEvery(s.cfg.CrawlEvery()+10*time.Minute, func() error {
		// TODO: process this in batches to reduce RAM load
		pendingTitles, err := s.dataProvider.TitlesProvider().ByStatus(model.StatusPending, model.StatusFailed).Select(ctx)
		if err != nil {
			if !errors.Is(err, data.ErrNotFound) {
				return errors.Wrap(err, "failed to select pending titles")
			}

			s.log.Debug("Have not found pending titles records")
			return nil
		}

		body, statusCodes, errs := s.newsCrawler.Crawl(ctx, pendingTitles)
		if err != nil {
			if !errors.Is(err, data.ErrNotFound) {
				return errors.Wrap(err, "failed to crawl url")
			}
			return nil
		}

		failedIDs := set.NewSet[uuid.UUID]()
		for i, statusCode := range statusCodes {
			if statusCode != http.StatusOK {
				s.log.WithFields(logrus.Fields{
					"title-id":    pendingTitles[i].ID,
					"title-url":   pendingTitles[i].URL,
					"status-code": statusCode,
					"info":        body,
				}).Warn("request returned unsuccessful status code...")

				failedIDs.Put(pendingTitles[i].ID)
			}
		}

		for i, err := range errs {
			if err != nil {
				s.log.WithFields(logrus.Fields{
					"title-id":  pendingTitles[i].ID,
					"title-url": pendingTitles[i].URL,
					"info":      body,
				}).WithError(err).Error("failed to run request...")

				failedIDs.Put(pendingTitles[i].ID)
			}
		}

		successIDs := set.NewSet[uuid.UUID]()
		for i := range pendingTitles {
			if errs[i] == nil && statusCodes[i] == http.StatusOK {
				successIDs.Put(pendingTitles[i].ID)
			}
		}

		if len(body) == 0 {
			s.log.Debug("early stopping, no new titles...")
			return nil
		}

		rawNewsBatch := crawler.ToModelBatch[model.RawNews](body)
		s.log.Debugf("Adding new batch to the database: %d", len(body))
		err = s.dataProvider.RawNewsProvider().InsertBatch(ctx, rawNewsBatch)
		if err != nil {
			return errors.Wrap(err, "failed to insert batch of titles")
		}

		err = s.updateStatusForProcessed(ctx, mapFilter(successIDs, rawNewsBatch), model.StatusProcessed)
		if err != nil {
			return errors.Wrap(err, "failed to update titles status to processed")
		}

		err = s.updateStatusForProcessed(ctx, mapFilter(failedIDs, rawNewsBatch), model.StatusFailed)
		if err != nil {
			return errors.Wrap(err, "failed to update titles status to failed")
		}

		return nil
	})
}

func mapFilter(s set.Set[uuid.UUID], rawNews []model.RawNews) []uuid.UUID {
	return iteration.Map(
		iteration.Filter(rawNews,
			func(news model.RawNews) bool {
				return s.Has(news.TitleID)
			},
		),
		func(t model.RawNews) uuid.UUID {
			return t.TitleID
		},
	)
}

func (s *service) updateStatusForProcessed(ctx context.Context, processedIDs []uuid.UUID, status string) error {
	if _, err := s.dataProvider.TitlesProvider().ByIDs(processedIDs).Update(ctx, model.UpdateTitleParams{
		Status: convert.ToPtr(status),
	}); err != nil {
		return errors.Wrap(err, "failed to update titles status")
	}
	return nil
}
