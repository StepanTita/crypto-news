package kv_provider

import (
	"context"
	"encoding/json"
	"time"

	rediscli "github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common/data"
	"common/data/queriers"
)

type kv struct {
	log     *logrus.Entry
	kvStore *rediscli.Client
}

func New(kvStore *rediscli.Client, log *logrus.Entry) queriers.KVProvider {
	return &kv{
		log:     log,
		kvStore: kvStore,
	}
}

func (k kv) Get(ctx context.Context, key string) (string, error) {
	res, err := k.kvStore.WithContext(ctx).Get(key).Result()
	if err != nil {
		if errors.Is(err, rediscli.Nil) {
			return "", data.ErrNotFound
		}
		return "", errors.Wrapf(err, "failed to get by key: %s from redis", key)
	}

	return res, nil
}

func (k kv) GetStruct(ctx context.Context, key string, out any) error {
	body, err := k.kvStore.WithContext(ctx).Get(key).Result()
	if err != nil {
		if errors.Is(err, rediscli.Nil) {
			return data.ErrNotFound
		}
		return errors.Wrapf(err, "failed to get by key: %s from redis", key)
	}

	if err := json.Unmarshal([]byte(body), out); err != nil {
		return errors.Wrap(err, "failed to unmarshal struct from redis")
	}

	return nil
}

func (k kv) SetValue(ctx context.Context, key, value string, exp time.Duration) (string, error) {
	id, err := k.kvStore.WithContext(ctx).Set(key, value, exp).Result()
	if err != nil {
		return "", errors.Wrapf(err, "failed to set by key: %s to redis", key)
	}
	return id, nil
}

func (k kv) SetStruct(ctx context.Context, key string, value any, exp time.Duration) (string, error) {
	body, err := json.Marshal(value)
	if err != nil {
		return "", errors.Wrapf(err, "failed to marshal struct: %T, with key: %s", value, key)
	}
	id, err := k.kvStore.WithContext(ctx).Set(key, body, exp).Result()
	if err != nil {
		return "", errors.Wrapf(err, "failed to set by key: %s to redis", key)
	}
	return id, nil
}

func (k kv) Remove(ctx context.Context, key string) error {
	if err := k.kvStore.WithContext(ctx).Del(key).Err(); err != nil {
		return errors.Wrapf(err, "failed to remove entity by key: %s in redis", key)
	}
	return nil
}
