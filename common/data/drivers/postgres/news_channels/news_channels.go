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
	var entity model.NewsChannel
	newsChannelColumns := model.PrependTableName(entity.TableName(), model.Columns(entity, false))
	return &newsChannels{
		log: log.WithField("provider", "news-channels"),
		db:  ext,

		Inserter: postgres.NewInserter[model.NewsChannel](ext, log),
		Selector: postgres.NewSelector[model.NewsChannel](ext, log, newsChannelColumns),
		Remover:  postgres.NewRemover[model.NewsChannel](ext, log),

		expr: common.BasicSqlizer,
	}
}

func (n newsChannels) Ordered() queriers.NewsChannelsProvider {
	n.Selector = n.Order("news_channels.priority", common.OrderAsc)
	return n
}

func (n newsChannels) BySources(source []string) queriers.NewsChannelsProvider {
	n.expr = sq.And{n.expr, sq.Eq{"news.source": source}}
	return n
}

func (n newsChannels) Select(ctx context.Context) ([]model.NewsChannel, error) {
	return n.Selector.Join("news", "news.id=news_channels.news_id").WithExpr(n.expr).Select(ctx)
}

func (n newsChannels) Remove(ctx context.Context, entity model.NewsChannel) error {
	return n.Remover.WithExpr(n.expr).Join([]string{"news"}, "news.id=news_channels.news_id").Remove(ctx, entity)
}
