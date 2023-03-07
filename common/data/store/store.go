package store

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common/data/drivers/postgres/channels"
	"common/data/drivers/postgres/news_channels"

	"common/config"
	"common/data/drivers/postgres/coins"
	"common/data/drivers/postgres/news"
	"common/data/drivers/postgres/news_coins"
	"common/data/drivers/postgres/users"
	"common/data/queriers"
)

type DataProvider interface {
	NewsProvider() queriers.NewsProvider
	UsersProvider() queriers.UserProvider
	CoinsProvider() queriers.CoinsProvider
	ChannelsProvider() queriers.ChannelsProvider
	NewsCoinsProvider() queriers.NewsCoinsProvider
	NewsChannelsProvider() queriers.NewsChannelsProvider

	InTx(ctx context.Context, fn func(dp DataProvider) error) error
}

func (d dataProvider) UsersProvider() queriers.UserProvider {
	return users.New(d.ext(), d.log)
}

func (d dataProvider) NewsProvider() queriers.NewsProvider {
	return news.New(d.ext(), d.log)
}

func (d dataProvider) CoinsProvider() queriers.CoinsProvider {
	return coins.New(d.ext(), d.log)
}

func (d dataProvider) ChannelsProvider() queriers.ChannelsProvider {
	return channels.New(d.ext(), d.log)
}

func (d dataProvider) NewsCoinsProvider() queriers.NewsCoinsProvider {
	return news_coins.New(d.ext(), d.log)
}

func (d dataProvider) NewsChannelsProvider() queriers.NewsChannelsProvider {
	return news_channels.New(d.ext(), d.log)
}

func (d dataProvider) InTx(ctx context.Context, fn func(dp DataProvider) error) error {
	tx, err := d.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	})
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	if err = fn(d.new(tx, d.db, d.log.WithField("tx", "[TRANSACTION]"))); err != nil {
		return errors.Wrap(err, "failed to run transaction")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit tx")
	}

	return nil
}

type dataProvider struct {
	log  *logrus.Entry
	db   *sqlx.DB
	tx   *sqlx.Tx
	inTx bool
}

func New(cfg config.Config) DataProvider {
	logging := cfg.Logging().WithField("[SQL]", cfg.Driver())
	logging.Info("Data provider connecting...")

	return &dataProvider{
		log:  logging,
		db:   cfg.DB(),
		tx:   nil,
		inTx: false,
	}
}

func (d dataProvider) new(tx *sqlx.Tx, db *sqlx.DB, log *logrus.Entry) DataProvider {
	return &dataProvider{
		db:   db,
		tx:   tx,
		inTx: true,
		log:  log,
	}
}

func (d dataProvider) ext() sqlx.ExtContext {
	if d.inTx {
		return d.tx
	}
	return d.db
}
