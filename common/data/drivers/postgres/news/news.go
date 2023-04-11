package news

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"common"
	"common/data/drivers/postgres"
	"common/data/model"
	"common/data/queriers"
)

type news struct {
	log *logrus.Entry
	ext sqlx.ExtContext

	expr sq.Sqlizer

	postgres.Inserter[model.News]
	postgres.Getter[model.News]
	postgres.Selector[model.News]
	postgres.Updater[model.UpdateNewsParams, model.News]
}

func New(ext sqlx.ExtContext, log *logrus.Entry) queriers.NewsProvider {
	return &news{
		log: log.WithField("provider", "news"),
		ext: ext,

		Inserter: postgres.NewInserter[model.News](ext, log),
		Getter:   postgres.NewGetter[model.News](ext, log),
		Selector: postgres.NewSelector[model.News](ext, log),
		Updater:  postgres.NewUpdater[model.UpdateNewsParams, model.News](ext, log),

		expr: common.BasicSqlizer,
	}
}

func (n news) BySource(source string) queriers.NewsProvider {
	n.expr = sq.And{n.expr, sq.Eq{"source": source}}
	return n
}

func (n news) ByStatus(status ...string) queriers.NewsProvider {
	n.expr = sq.And{n.expr, sq.Eq{"status": status}}
	return n
}

func (n news) ByIDs(ids []uuid.UUID) queriers.NewsProvider {
	n.expr = sq.And{n.expr, sq.Eq{"id": ids}}
	return n
}

func (n news) ByCoins(codes []string) queriers.NewsProvider {
	n.Selector = n.Selector.Join("news_coins", sq.Eq{"news_coins.code": codes})
	return n
}

func (n news) GetLatest(ctx context.Context) (*model.News, error) {
	n.Getter = n.Getter.Order("published_at", common.OrderDesc)
	return n.Get(ctx)
}

func (n news) Get(ctx context.Context) (*model.News, error) {
	n.Getter = n.Getter.WithExpr(n.expr)
	return n.Getter.Get(ctx)
}
