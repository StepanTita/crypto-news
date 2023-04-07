package preferences_channel_coins

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"common"
	"common/data/drivers/postgres"
	"common/data/model"
	"common/data/queriers"
)

type preferencesChannelCoins struct {
	log *logrus.Entry
	db  sqlx.ExtContext

	expr sq.Sqlizer

	postgres.Inserter[model.PreferencesChannelCoin]
	postgres.Selector[model.PreferencesChannelCoin]
	postgres.Remover[model.PreferencesChannelCoin]
}

func New(ext sqlx.ExtContext, log *logrus.Entry) queriers.PreferencesChannelCoinsProvider {
	return &preferencesChannelCoins{
		log: log.WithField("provider", "preferences-channel-coins"),
		db:  ext,

		Inserter: postgres.NewInserter[model.PreferencesChannelCoin](ext, log),
		Selector: postgres.NewSelector[model.PreferencesChannelCoin](ext, log),
		Remover:  postgres.NewRemover[model.PreferencesChannelCoin](ext, log),

		expr: common.BasicSqlizer,
	}
}

func (p preferencesChannelCoins) ByChannel(channelID int64) queriers.PreferencesChannelCoinsProvider {
	p.expr = sq.And{p.expr, sq.Eq{"channel_id": channelID}}
	return p
}
