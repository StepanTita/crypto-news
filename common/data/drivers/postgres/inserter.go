package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	sq "github.com/Masterminds/squirrel"

	"common/data"
	"common/data/model"
	"common/data/queriers"
)

type Inserter[T model.Model] interface {
	queriers.Inserter[T]
}

type inserter[T model.Model] struct {
	log *logrus.Entry
	ext sqlx.ExtContext

	sql sq.InsertBuilder
}

func NewInserter[T model.Model](ext sqlx.ExtContext, log *logrus.Entry) Inserter[T] {
	var entity T
	return &inserter[T]{
		log: log.WithField("service", "[inserter]"),
		ext: ext,
		sql: sq.Insert(entity.TableName()),
	}
}

func (i inserter[T]) Insert(ctx context.Context, entity T) (*T, error) {
	sql, args, err := i.sql.SetMap(model.ToMap(entity)).Suffix("RETURNING *").ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build sql insert query")
	}

	rows, err := i.ext.QueryxContext(ctx, i.ext.Rebind(sql), args...)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == ErrCodeUniqueViolation {
			// Here err is of type *pq.Error, you may inspect all its fields, e.g.:
			return nil, data.ErrDuplicateRecord
		}
		return nil, errors.Wrapf(err, "failed to insert entity into table: %s", entity.TableName())
	}

	for rows.Next() {
		err := rows.StructScan(&entity)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan entity")
		}
	}
	return &entity, nil
}

func (i inserter[T]) InsertBatch(ctx context.Context, entities []T) error {
	if len(entities) == 0 {
		return nil
	}

	tableName := entities[0].TableName()

	columns, namedBindings := model.NamedBinding(entities[0])
	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s) RETURNING *`, tableName, strings.Join(columns, ","), strings.Join(namedBindings, ","))

	i.log.Debug(sql)

	rows, err := sqlx.NamedQueryContext(ctx, i.ext, i.ext.Rebind(sql), entities)
	if err != nil {
		return errors.Wrapf(err, "failed to insert entity into table: %s", tableName)
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
