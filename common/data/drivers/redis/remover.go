package redis

import (
	"context"

	rediscli "github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common/data/model"
	"common/data/queriers"
)

type Remover[T model.Model] interface {
	queriers.Remover[T]
}

type remover[T model.Model] struct {
	log *logrus.Entry

	kvStore *rediscli.Client
}

func NewRemover[T model.Model](kvStore *rediscli.Client, log *logrus.Entry) Remover[T] {
	return &remover[T]{
		log: log.WithField("service", "[nosql-remover]"),

		kvStore: kvStore,
	}
}

func (r remover[T]) Remove(ctx context.Context, entity T) error {
	if err := r.kvStore.WithContext(ctx).Del(model.ToKey(entity, false)).Err(); err != nil {
		return errors.Wrap(err, "failed to remove entity in redis")
	}
	return nil
}
