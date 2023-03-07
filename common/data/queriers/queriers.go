package queriers

import (
	"context"

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

type NewsProvider interface {
	Inserter[model.News]
	Getter[model.News]
	Selector[model.News]

	BySource(ctx context.Context, source string) NewsProvider
	ByIDs(ctx context.Context, ids []uuid.UUID) NewsProvider

	GetLatest(ctx context.Context) (*model.News, error)
}

type UserProvider interface {
	Inserter[model.User]
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
	Updater[model.UpdateNewsChannelParams, model.NewsChannel]

	ByStatus(ctx context.Context, source string) NewsChannelsProvider

	// orders by priority
	Ordered(ctx context.Context) NewsChannelsProvider
}
