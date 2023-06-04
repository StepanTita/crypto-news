package queriers

import (
	"context"
	"time"

	"github.com/google/uuid"

	"common/data/model"
)

type Inserter[T model.Model] interface {
	Insert(ctx context.Context, entity T) (*T, error)
	InsertBatch(ctx context.Context, entities []T) error
}

type Getter[T model.Model] interface {
	Get(ctx context.Context) (*T, error)
}

type Selector[T model.Model] interface {
	Select(ctx context.Context) ([]T, error)
}

type Updater[U model.Model, T model.Model] interface {
	Update(ctx context.Context, entity U) ([]T, error)
}

type Remover[T model.Model] interface {
	Remove(ctx context.Context, entity T) error
}

type NewsProvider interface {
	Inserter[model.News]
	Getter[model.News]
	Selector[model.News]
	Updater[model.UpdateNewsParams, model.News]

	ByStatus(status ...string) NewsProvider

	BySources(sources ...string) NewsProvider
	ByIDs(ids []uuid.UUID) NewsProvider

	// ByCoins TODO: maybe implement results filtering by coins later (but this might slow down,
	// since then we need to get news for each channel independently
	ByCoins(codes []string) NewsProvider

	GetLatest(ctx context.Context) (*model.News, error)
}

type CoinsProvider interface {
	UpsertCoinsBatch(ctx context.Context, coins []model.Coin) error
}

type ChannelsProvider interface {
	Inserter[model.Channel]
	Selector[model.Channel]
}

type NewsCoinsProvider interface {
	Inserter[model.NewsCoin]
}

type NewsChannelsProvider interface {
	Inserter[model.NewsChannel]
	Selector[model.NewsChannel]
	Remover[model.NewsChannel]

	// Ordered orders by priority
	Ordered() NewsChannelsProvider

	BySources(source []string) NewsChannelsProvider
	ByIDs(ids []uuid.UUID) NewsChannelsProvider
}

type PreferencesChannelCoinsProvider interface {
	Inserter[model.PreferencesChannelCoin]
	Selector[model.PreferencesChannelCoin]
	Remover[model.PreferencesChannelCoin]

	ByChannel(channelID int64) PreferencesChannelCoinsProvider
}

// No-SQL

type KVProvider interface {
	Get(ctx context.Context, key string) (string, error)
	GetStruct(ctx context.Context, key string, out any) error

	SetValue(ctx context.Context, key, value string, exp time.Duration) (string, error)
	SetStruct(ctx context.Context, key string, value any, exp time.Duration) (string, error)

	Remove(ctx context.Context, key string) error
}
