package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common"
	"common/data/model"
	"common/data/queriers"
)

type Updater[U, T model.Model] interface {
	queriers.Updater[U, T]

	WithExpr(expr sq.Sqlizer) Updater[U, T]
}

type updater[U, T model.Model] struct {
	log *logrus.Entry
	ext sqlx.ExtContext

	sql  sq.UpdateBuilder
	expr sq.Sqlizer
}

func NewUpdater[U, T model.Model](ext sqlx.ExtContext, log *logrus.Entry) Updater[U, T] {
	var entity U
	return &updater[U, T]{
		log: log.WithField("service", "[updater]"),
		ext: ext,

		sql: sq.Update(entity.TableName()),

		expr: common.BasicSqlizer,
	}
}

func (u updater[U, T]) Update(ctx context.Context, entity U) ([]T, error) {
	u.log.Debug(sq.DebugSqlizer(u.sql.Where(u.expr)))

	sql, args, err := u.sql.SetMap(model.ToMap(entity)).Where(u.expr).Suffix("RETURNING *").ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build sql update query")
	}

	rows, err := u.ext.QueryxContext(ctx, u.ext.Rebind(sql), args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update entity: %s", entity.TableName())
	}

	updatedEntities := make([]T, 0, 10)
	for rows.Next() {
		var emptyEntity T
		i := len(updatedEntities)
		updatedEntities = append(updatedEntities, emptyEntity)

		if err := rows.StructScan(&updatedEntities[i]); err != nil {
			return nil, errors.Wrap(err, "failed to scan updated entity")
		}
	}
	return updatedEntities, nil
}

func (u updater[U, T]) WithExpr(expr sq.Sqlizer) Updater[U, T] {
	u.expr = expr
	return u
}
