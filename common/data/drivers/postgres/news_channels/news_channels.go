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
	postgres.Selector[model.NewsChannel]
	postgres.Remover[model.NewsChannel]
}

func New(ext sqlx.ExtContext, log *logrus.Entry) queriers.NewsChannelsProvider {
	return &newsChannels{
		log: log.WithField("provider", "news-channels"),
		db:  ext,

		Inserter: postgres.NewInserter[model.NewsChannel](ext, log),
		Selector: postgres.NewSelector[model.NewsChannel](ext, log),
		Remover:  postgres.NewRemover[model.NewsChannel](ext, log),

		expr: common.BasicSqlizer,
	}
}

func (n newsChannels) Ordered(ctx context.Context) queriers.NewsChannelsProvider {
	n.Selector = n.Order("priority", common.OrderAsc)
	return n
}
