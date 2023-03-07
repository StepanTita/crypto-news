package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common/data/model"
	"common/data/queriers"
)

type Getter[T model.Model] interface {
	queriers.Getter[T]

	WithExpr(expr sq.Sqlizer) Getter[T]

	Order(by, order string) Getter[T]
}

type getter[T model.Model] struct {
	log *logrus.Entry

	selector Selector[T]

	expr sq.Sqlizer
}

func NewGetter[T model.Model](ext sqlx.ExtContext, log *logrus.Entry) Getter[T] {
	return &getter[T]{
		log:      log.WithField("service", "[getter]"),
		selector: NewSelector[T](ext, log).Limit(1),
	}
}

func (g getter[T]) WithExpr(expr sq.Sqlizer) Getter[T] {
	g.expr = expr
	return g
}

func (g getter[T]) Get(ctx context.Context) (*T, error) {
	out, err := g.selector.WithExpr(g.expr).Select(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get record")
	}

	if len(out) == 0 {
		return nil, nil
	}

	if len(out) > 1 {
		return nil, errors.New("get method returned more than one row")
	}

	return &out[0], nil
}

func (g getter[T]) Order(by, order string) Getter[T] {
	g.selector = g.selector.Order(by, order)
	return g
}
