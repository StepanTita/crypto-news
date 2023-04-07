package authorization_keys

import (
	"github.com/sirupsen/logrus"

	rediscli "github.com/go-redis/redis"

	"common/data/drivers/redis"
	"common/data/model"
	"common/data/queriers"
)

type authorizationKeys struct {
	log     *logrus.Entry
	kvStore *rediscli.Client

	redis.Getter[model.AuthorizationKeys]
	redis.Inserter[model.AuthorizationKeys]
	redis.Remover[model.AuthorizationKeys]
}

func New(kvStore *rediscli.Client, log *logrus.Entry) queriers.AuthorizationKeysProvider {
	return &authorizationKeys{
		log: log,

		kvStore:  kvStore,
		Getter:   redis.NewGetter[model.AuthorizationKeys](kvStore, log),
		Inserter: redis.NewInserter[model.AuthorizationKeys](kvStore, log),
		Remover:  redis.NewRemover[model.AuthorizationKeys](kvStore, log),
	}
}
