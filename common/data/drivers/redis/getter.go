package redis

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common/data"
	"common/data/model"
	"common/data/queriers"
)

type Getter[T model.Model] interface {
	queriers.Getter[T]

	ByKey(key string) Getter[T]
}

type getter[T model.Model] struct {
	log *logrus.Entry

	kvStore *redis.Client

	key string
}

func NewGetter[T model.Model](kvStore *redis.Client, log *logrus.Entry) Getter[T] {
	return &getter[T]{
		log:     log.WithField("service", "[nosql-getter]"),
		kvStore: kvStore,
	}
}

func (g getter[T]) Get(ctx context.Context) (*T, error) {
	var res T
	if err := g.kvStore.WithContext(ctx).Get(g.key).Scan(&res); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, data.ErrNotFound
		}
		return nil, errors.Wrapf(err, "failed to get by key: %s from redis", g.key)
	}

	return &res, nil
}

func (g getter[T]) ByKey(key string) Getter[T] {
	g.key = key
	return g
}
