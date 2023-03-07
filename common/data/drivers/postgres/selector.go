package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common"
	"common/data"
	"common/data/model"
	"common/data/queriers"
)

type Selector[T model.Model] interface {
	queriers.Selector[T]

	WithExpr(expr sq.Sqlizer) Selector[T]

	Limit(l uint64) Selector[T]
	Order(by, order string) Selector[T]
}

type selector[T model.Model] struct {
	log *logrus.Entry
	ext sqlx.ExtContext

	//TODO: can use that instead of db to support transactions
	//ext sqlx.Ext
	expr sq.Sqlizer
	sql  sq.SelectBuilder
}

func NewSelector[T model.Model](ext sqlx.ExtContext, log *logrus.Entry) Selector[T] {
	var entity T

	return &selector[T]{
		log: log.WithField("service", "[selector]"),
		ext: ext,

		expr: common.BasicSqlizer,
		sql:  sq.Select(model.Columns(entity, false)...).From(entity.TableName()),
	}
}

func (s selector[T]) Select(ctx context.Context) ([]T, error) {
	var entity T

	sql, args, err := s.sql.Where(s.expr).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build sql insert query")
	}

	rows, err := s.ext.QueryxContext(ctx, s.ext.Rebind(sql), args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to insert entities into table: %selector", entity.TableName())
	}

	entities := make([]T, 0, 10)
	for rows.Next() {
		var emptyEntity T
		i := len(entities)
		entities = append(entities, emptyEntity)

		if err := rows.StructScan(&entities[i]); err != nil {
			return nil, errors.Wrap(err, "failed to scan selected entity")
		}
	}

	if len(entities) == 0 {
		return nil, data.ErrNotFound
	}
	return entities, nil
}

func (s selector[T]) WithExpr(expr sq.Sqlizer) Selector[T] {
	s.expr = expr
	return s
}

func (s selector[T]) Limit(l uint64) Selector[T] {
	s.sql = s.sql.Limit(l)
	return s
}

func (s selector[T]) Order(by, order string) Selector[T] {
	s.sql = s.sql.OrderBy(fmt.Sprintf("%s %s", by, order))
	return s
}
