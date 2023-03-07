package coins

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common/data/model"
	"common/data/queriers"
)

type coins struct {
	log *logrus.Entry
	ext sqlx.ExtContext
}

func New(ext sqlx.ExtContext, log *logrus.Entry) queriers.CoinsProvider {
	return &coins{
		log: log.WithField("provider", "coins"),
		ext: ext,
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
