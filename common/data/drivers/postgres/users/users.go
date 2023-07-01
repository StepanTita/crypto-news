package users

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

type users struct {
	log *logrus.Entry
	ext sqlx.ExtContext

	expr sq.Sqlizer

	postgres.Inserter[model.User]
	postgres.Getter[model.User]
	postgres.Selector[model.User]
}

func New(ext sqlx.ExtContext, log *logrus.Entry) queriers.UsersProvider {
	var entity model.User
	usersColumns := model.PrependTableName(entity.TableName(), model.Columns(entity, false))
	return &users{
		log: log.WithField("provider", "users"),
		ext: ext,

		Inserter: postgres.NewInserter[model.User](ext, log),
		Getter:   postgres.NewGetter[model.User](ext, log, usersColumns),
		Selector: postgres.NewSelector[model.User](ext, log, usersColumns),

		expr: common.BasicSqlizer,
	}
}

func (u users) Get(ctx context.Context) (*model.User, error) {
	u.Getter = u.Getter.WithExpr(u.expr)
	return u.Getter.Get(ctx)
}

func (u users) Select(ctx context.Context) ([]model.User, error) {
	u.Selector = u.Selector.WithExpr(u.expr)
	return u.Selector.Select(ctx)
}

func (u users) ByUsername(username string) queriers.UsersProvider {
	u.expr = sq.And{u.expr, sq.Eq{"users.username": username}}
	return u
}
