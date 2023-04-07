package config

import (
	rediscli "github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type KVStorer interface {
	KVStore() *rediscli.Client
}

type YamlKVStoreConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

type kvStorer struct {
	kvStore *rediscli.Client
}

func NewKVStorer(kvStoreConfig YamlKVStoreConfig) KVStorer {
	kvStore := rediscli.NewClient(&rediscli.Options{
		Addr:     kvStoreConfig.Address,
		Password: kvStoreConfig.Password,
		DB:       0,
	})

	if err := kvStore.Ping().Err(); err != nil {
		panic(errors.Errorf("couldn't ping kv store: %s", kvStoreConfig.Address))
	}
	return &kvStorer{
		kvStore: kvStore,
	}
}

func (s *kvStorer) KVStore() *rediscli.Client {
	return s.kvStore
}
