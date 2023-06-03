package channels

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"common/data/drivers/postgres"
	"common/data/model"
	"common/data/queriers"
)

type channels struct {
	postgres.Inserter[model.Channel]
	postgres.Selector[model.Channel]
}

func New(ext sqlx.ExtContext, log *logrus.Entry) queriers.ChannelsProvider {
	var entity model.Channel
	channelColumns := model.PrependTableName(entity.TableName(), model.Columns(entity, false))
	return &channels{
		Inserter: postgres.NewInserter[model.Channel](ext, log),
		Selector: postgres.NewSelector[model.Channel](ext, log, channelColumns),
	}
}
