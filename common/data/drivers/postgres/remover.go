package postgres

import (
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common"
	"common/data"
	"common/data/model"
	"common/data/queriers"
)

type Remover[T model.Model] interface {
	queriers.Remover[T]

	WithExpr(expr sq.Sqlizer) Remover[T]

	Join(to []string, on string, args ...interface{}) Remover[T]
}

type remover[T model.Model] struct {
	log *logrus.Entry
	ext sqlx.ExtContext

	expr sq.Sqlizer
	sql  sq.DeleteBuilder
}

func NewRemover[T model.Model](ext sqlx.ExtContext, log *logrus.Entry) Remover[T] {
	var entity T

	return &remover[T]{
		log: log.WithField("service", "[remover]"),
		ext: ext,

		expr: common.BasicSqlizer,
		sql:  sq.Delete(entity.TableName()),
	}
}

// TODO: remove by fields that are set on entity
func (r remover[T]) Remove(ctx context.Context, entity T) error {
	r.log.Debug(sq.DebugSqlizer(r.sql.Where(r.expr)))

	sql, args, err := r.sql.Where(r.expr).ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build sql remove query")
	}

	result, err := r.ext.ExecContext(ctx, r.ext.Rebind(sql), args...)
	if err != nil {
		return errors.Wrapf(err, "failed to remove entities from table: %s", entity.TableName())
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to count affected rows")
	}

	if n == 0 {
		return data.ErrNotFound
	}

	r.log.WithField("rows", n).Debug("removed rows")
	return nil
}

func (r remover[T]) WithExpr(expr sq.Sqlizer) Remover[T] {
	r.expr = expr
	return r
}

func (r remover[T]) Join(to []string, on string, args ...interface{}) Remover[T] {
	var entity T
	r.sql = r.sql.From(fmt.Sprintf("%s USING %s", entity.TableName(), strings.Join(to, ","))).Where(on, args...)
	return r
}
