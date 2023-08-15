package raw_news_webpages

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

type rawNewsWebpages struct {
	log *logrus.Entry
	ext sqlx.ExtContext

	expr sq.Sqlizer

	postgres.Inserter[model.RawNewsWebpage]
	postgres.Selector[model.RawNewsWebpage]
	postgres.Remover[model.RawNewsWebpage]
}

func New(ext sqlx.ExtContext, log *logrus.Entry) queriers.RawNewsWebpagesProvider {
	var entity model.RawNewsWebpage
	whitelistColumns := model.PrependTableName(entity.TableName(), model.Columns(entity, false))
	return &rawNewsWebpages{
		log: log.WithField("provider", "raw_news_webpages"),
		ext: ext,

		Inserter: postgres.NewInserter[model.RawNewsWebpage](ext, log),
		Selector: postgres.NewSelector[model.RawNewsWebpage](ext, log, whitelistColumns),
		Remover:  postgres.NewRemover[model.RawNewsWebpage](ext, log),

		expr: common.BasicSqlizer,
	}
}

func (w rawNewsWebpages) Remove(ctx context.Context, entity model.RawNewsWebpage) error {
	w.Remover = w.Remover.WithExpr(w.expr)
	return w.Remover.Remove(ctx, entity)
}

func (w rawNewsWebpages) Select(ctx context.Context) ([]model.RawNewsWebpage, error) {
	w.Selector = w.Selector.WithExpr(w.expr)
	return w.Selector.Select(ctx)
}
