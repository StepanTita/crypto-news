package news

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"common"
	"common/convert"
	"common/data"
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
	var entity model.News
	newsColumns := model.PrependTableName(entity.TableName(), model.Columns(entity, false))
	return &news{
		log: log.WithField("provider", "news"),
		ext: ext,

		Inserter: postgres.NewInserter[model.News](ext, log),
		Getter:   postgres.NewGetter[model.News](ext, log, newsColumns),
		Selector: postgres.NewSelector[model.News](ext, log, newsColumns),
		Updater:  postgres.NewUpdater[model.UpdateNewsParams, model.News](ext, log),

		expr: data.BasicSqlizer,
	}
}

func (n news) BySources(sources ...string) queriers.NewsProvider {
	n.expr = sq.And{n.expr, sq.Eq{"news.source": sources}}
	return n
}

func (n news) ByStatus(status ...string) queriers.NewsProvider {
	n.expr = sq.And{n.expr, sq.Eq{"news.status": status}}
	return n
}

func (n news) ByIDs(ids []uuid.UUID) queriers.NewsProvider {
	n.expr = sq.And{n.expr, sq.Eq{"news.id": ids}}
	return n
}

func (n news) ByCoins(codes []string) queriers.NewsProvider {
	sql, args, _ := sq.Eq{"news_coins.code": codes}.ToSql()
	n.Selector = n.Selector.Join("news_coins", sql, args...)
	return n
}

func (n news) GetLatest(ctx context.Context) (*model.News, error) {
	n.Getter = n.Getter.Order("news.published_at", data.OrderDesc)
	return n.Get(ctx)
}

func (n news) Get(ctx context.Context) (*model.News, error) {
	n.Getter = n.Getter.WithExpr(n.expr)
	return n.Getter.Get(ctx)
}

func (n news) Select(ctx context.Context) ([]model.News, error) {
	n.Selector = n.Selector.WithExpr(n.expr)
	return n.Selector.Select(ctx)
}

func (n news) Update(ctx context.Context, news model.UpdateNewsParams) ([]model.News, error) {
	n.Updater = n.Updater.WithExpr(n.expr)

	news.UpdatedAt = convert.ToPtr(common.CurrentTimestamp())
	return n.Updater.Update(ctx, news)
}
