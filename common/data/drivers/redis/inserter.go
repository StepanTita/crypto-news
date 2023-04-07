package redis

import (
	"context"
	"encoding/json"

	rediscli "github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common/data/model"
	"common/data/queriers"
)

type Inserter[T model.Model] interface {
	queriers.Inserter[T]
}

type inserter[T model.Model] struct {
	log *logrus.Entry

	kvStore *rediscli.Client
}

func NewInserter[T model.Model](kvStore *rediscli.Client, log *logrus.Entry) Inserter[T] {
	return &inserter[T]{
		log:     log.WithField("service", "[nosql-inserter]"),
		kvStore: kvStore,
	}
}

func (i inserter[T]) Insert(ctx context.Context, entity T) (*T, error) {
	b, err := json.Marshal(entity)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal entity for insertion")
	}

	if err := i.kvStore.WithContext(ctx).Set(model.ToKey(entity, false), b, 0).Err(); err != nil {
		return nil, errors.Wrap(err, "failed to set entity to the redis store")
	}
	return &entity, nil
}

// TODO: what if we need these keys in the future?
func (i inserter[T]) InsertBatch(ctx context.Context, entities []T) error {
	pairs := make([]interface{}, len(entities))
	for _, e := range entities {
		b, err := json.Marshal(e)
		if err != nil {
			return errors.Wrap(err, "failed to marshal entity for batch insert")
		}
		pairs = append(pairs, model.ToKey(e, true), b)
	}
	if err := i.kvStore.WithContext(ctx).MSet(pairs).Err(); err != nil {
		return errors.Wrap(err, "failed to set pairs in redis")
	}
	return nil
}
