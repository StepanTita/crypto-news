package coins

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common"
	"common/data/drivers/postgres"
	"common/data/model"
	"common/data/queriers"
)

type coins struct {
	log *logrus.Entry
	ext sqlx.ExtContext

	postgres.Selector[model.Coin]
	expr sq.Sqlizer
}

func New(ext sqlx.ExtContext, log *logrus.Entry) queriers.CoinsProvider {
	var entity model.Coin
	coinsColumns := model.PrependTableName(entity.TableName(), model.Columns(entity, false))
	return &coins{
		log: log.WithField("provider", "coins"),
		ext: ext,

		Selector: postgres.NewSelector[model.Coin](ext, log, coinsColumns),

		expr: common.BasicSqlizer,
	}
}

func (c coins) UpsertCoinsBatch(ctx context.Context, entities []model.Coin) error {
	if len(entities) == 0 {
		return nil
	}

	sql := `
		INSERT INTO coins (code, title, slug) VALUES (:code, :title, :slug) 
        	ON CONFLICT (code) DO UPDATE 
                SET code=excluded.code, title=excluded.title, slug=excluded.slug RETURNING *`

	rows, err := sqlx.NamedQueryContext(ctx, c.ext, c.ext.Rebind(sql), entities)
	if err != nil {
		return errors.Wrap(err, "failed to insert entity into table: coins")
	}

	idx := 0
	for rows.Next() {
		err := rows.StructScan(&entities[idx])
		if err != nil {
			return errors.Wrap(err, "failed to scan entity")
		}
		idx++
	}
	return nil
}

func (c coins) ByNewsID(id uuid.UUID) queriers.CoinsProvider {
	c.expr = sq.And{c.expr, sq.Eq{"news_coins.news_id": id}}
	return c
}

func (c coins) Select(ctx context.Context) ([]model.Coin, error) {
	return c.Selector.Join("news_coins", "coins.code=news_coins.code").WithExpr(c.expr).Select(ctx)
}
