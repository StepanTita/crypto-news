package authenticator

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"twitter-bot/internal/config"
)

type Service interface {
	Listen(ctx context.Context) error
}

type service struct {
	log    *logrus.Entry
	cfg    config.Config
	router chi.Router
}

func New(cfg config.Config) Service {
	l := &service{
		cfg: cfg,
		log: cfg.Logging().WithField("service", "[authenticator]"),
	}
	l.setupRouter()

	return l
}

func (l *service) Listen(ctx context.Context) error {
	l.log.WithField("port", l.cfg.Port()).Info("Starting service...")

	if err := http.ListenAndServe(l.cfg.Port(), l.router); err != nil {
		return errors.Wrap(err, "service failed")
	}
	return nil
}
