package whitelist

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

type whitelist struct {
	log *logrus.Entry
	ext sqlx.ExtContext

	expr sq.Sqlizer

	postgres.Inserter[model.Whitelist]
	postgres.Getter[model.Whitelist]
	postgres.Remover[model.Whitelist]
}

func New(ext sqlx.ExtContext, log *logrus.Entry) queriers.WhitelistProvider {
	var entity model.Whitelist
	whitelistColumns := model.PrependTableName(entity.TableName(), model.Columns(entity, false))
	return &whitelist{
		log: log.WithField("provider", "whitelist"),
		ext: ext,

		Inserter: postgres.NewInserter[model.Whitelist](ext, log),
		Getter:   postgres.NewGetter[model.Whitelist](ext, log, whitelistColumns),
		Remover:  postgres.NewRemover[model.Whitelist](ext, log),

		expr: common.BasicSqlizer,
	}
}

func (l whitelist) ByUsername(username string) queriers.WhitelistProvider {
	l.expr = sq.And{l.expr, sq.Eq{"whitelist.username": username}}
	return l
}

func (l whitelist) Get(ctx context.Context) (*model.Whitelist, error) {
	l.Getter = l.Getter.WithExpr(l.expr)
	return l.Getter.Get(ctx)
}

func (l whitelist) Remove(ctx context.Context, entity model.Whitelist) error {
	return l.Remover.WithExpr(l.expr).Remove(ctx, entity)
}

func (l whitelist) ExtractToken(ctx context.Context, token uuid.UUID) error {
	return l.Remover.WithExpr(sq.And{l.expr, sq.Eq{"whitelist.token": token}}).Remove(ctx, model.Whitelist{})
}
