package news_coins

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"common/data/drivers/postgres"
	"common/data/model"
	"common/data/queriers"
)

type newsCoins struct {
	log *logrus.Entry
	ext sqlx.ExtContext

	postgres.Inserter[model.NewsCoin]
}

func New(ext sqlx.ExtContext, log *logrus.Entry) queriers.NewsCoinsProvider {
	return &newsCoins{
		log: log.WithField("provider", "news_coins"),
		ext: ext,

		Inserter: postgres.NewInserter[model.NewsCoin](ext, log),
	}
}
