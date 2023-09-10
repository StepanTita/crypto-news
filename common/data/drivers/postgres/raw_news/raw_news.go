package raw_news

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"common/data"
	"common/data/drivers/postgres"
	"common/data/model"
	"common/data/queriers"
)

type rawNews struct {
	log *logrus.Entry
	ext sqlx.ExtContext

	expr sq.Sqlizer

	postgres.Inserter[model.RawNews]
	postgres.Selector[model.RawNews]
	postgres.Remover[model.RawNews]
}

func New(ext sqlx.ExtContext, log *logrus.Entry) queriers.RawNewsProvider {
	var entity model.RawNews
	whitelistColumns := model.PrependTableName(entity.TableName(), model.Columns(entity, false))
	return &rawNews{
		log: log.WithField("provider", "raw_news"),
		ext: ext,

		Inserter: postgres.NewInserter[model.RawNews](ext, log),
		Selector: postgres.NewSelector[model.RawNews](ext, log, whitelistColumns),
		Remover:  postgres.NewRemover[model.RawNews](ext, log),

		expr: data.BasicSqlizer,
	}
}

func (w rawNews) ByIDs(ids []uuid.UUID) queriers.RawNewsProvider {
	w.expr = sq.And{w.expr, sq.Eq{"raw_news.id": ids}}
	return w
}

func (w rawNews) Limit(l uint64) queriers.RawNewsProvider {
	w.Selector = w.Selector.Limit(l)
	return w
}

func (w rawNews) Offset(o uint64) queriers.RawNewsProvider {
	w.Selector = w.Selector.Offset(o)
	return w
}

func (w rawNews) Order(by, order string) queriers.RawNewsProvider {
	w.Selector = w.Selector.Order(by, order)
	return w
}

func (w rawNews) Remove(ctx context.Context, entity model.RawNews) error {
	w.Remover = w.Remover.WithExpr(w.expr)
	return w.Remover.Remove(ctx, entity)
}

func (w rawNews) Select(ctx context.Context) ([]model.RawNews, error) {
	w.Selector = w.Selector.WithExpr(w.expr)
	return w.Selector.Select(ctx)
}

func (w rawNews) Count(ctx context.Context) (uint64, error) {
	return w.Selector.Count(ctx)
}
