package poster

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common/convert"
	"common/data/model"
	"common/data/store"
	"twitter-bot/internal/config"
	"twitter-bot/internal/services/twitter"
)

type Poster interface {
	Post(ctx context.Context) (int, error)
}

type poster struct {
	cfg config.Config

	log *logrus.Entry

	dataProvider store.DataProvider

	twitter twitter.Client
}

func New(cfg config.Config) Poster {
	return &poster{
		cfg: cfg,

		log: cfg.Logging().WithField("service", "[TWITTER-POSTER]"),

		dataProvider: store.New(cfg),

		twitter: twitter.New(cfg),
	}
}

func (p poster) Post(ctx context.Context) (int, error) {
	news, err := p.dataProvider.NewsProvider().ByStatus(model.StatusPending, model.StatusFailed).BySources(p.cfg.Sources()...).Select(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to select news by status")
	}

	count := 0
	processedIDs := make([]uuid.UUID, 0, 10)
	failedIDs := make([]uuid.UUID, 0, 10)
	for _, n := range news {
		if err := p.twitter.PostTweet(ctx, n); err != nil {
			failedIDs = append(failedIDs, n.ID)
			p.log.WithError(err).Warn("failed to post tweet, adding to failed list...")
		} else {
			processedIDs = append(processedIDs, n.ID)
			count++
		}
	}

	if err := p.updateStatusForProcessed(ctx, processedIDs, model.StatusProcessed); err != nil {
		return count, errors.Wrap(err, "failed to update news status to 'processed'")
	}

	if err := p.updateStatusForProcessed(ctx, failedIDs, model.StatusFailed); err != nil {
		return count, errors.Wrap(err, "failed to update news status to 'failed'")
	}

	return count, nil
}

func (p poster) updateStatusForProcessed(ctx context.Context, processedIDs []uuid.UUID, status string) error {
	if _, err := p.dataProvider.NewsProvider().ByIDs(processedIDs).Update(ctx, model.UpdateNewsParams{
		Status: convert.ToPtr(status),
	}); err != nil {
		return errors.Wrap(err, "failed to update news status")
	}
	return nil
}
