package news_channels

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"common"
	"common/data/drivers/postgres"
	"common/data/model"
	"common/data/queriers"
)

type newsChannels struct {
	log *logrus.Entry
	db  sqlx.ExtContext

	expr sq.Sqlizer

	postgres.Inserter[model.NewsChannel]
	postgres.Updater[model.UpdateNewsChannelParams, model.NewsChannel]
	postgres.Selector[model.NewsChannel]
}

func New(ext sqlx.ExtContext, log *logrus.Entry) queriers.NewsChannelsProvider {
	return &newsChannels{
		log: log.WithField("provider", "news"),
		db:  ext,

		Inserter: postgres.NewInserter[model.NewsChannel](ext, log),
		Updater:  postgres.NewUpdater[model.UpdateNewsChannelParams, model.NewsChannel](ext, log),
		Selector: postgres.NewSelector[model.NewsChannel](ext, log),

		expr: common.BasicSqlizer,
	}
}

func (n newsChannels) ByStatus(_ context.Context, status string) queriers.NewsChannelsProvider {
	n.expr = sq.And{n.expr, sq.Eq{"status": status}}
	return n
}

func (n newsChannels) Ordered(ctx context.Context) queriers.NewsChannelsProvider {
	n.Selector = n.Order("priority", common.OrderAsc)
	return n
}

func (n newsChannels) Update(ctx context.Context, entity model.UpdateNewsChannelParams) ([]model.NewsChannel, error) {
	return n.Updater.WithExpr(n.expr).Update(ctx, entity)
}
