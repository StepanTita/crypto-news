package users

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"common/data/drivers/postgres"
	"common/data/model"
	"common/data/queriers"
)

type users struct {
	postgres.Inserter[model.User]
}

func New(ext sqlx.ExtContext, log *logrus.Entry) queriers.UserProvider {
	return &users{
		Inserter: postgres.NewInserter[model.User](ext, log),
	}
}
