package services

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common"

	"twitter-bot/internal/config"
	"twitter-bot/internal/services/authenticator"
	"twitter-bot/internal/services/poster"
)

type Service interface {
	Run(ctx context.Context) error
}

type service struct {
	cfg config.Config
	log *logrus.Entry
}

func New(cfg config.Config) Service {
	return &service{
		cfg: cfg,
		log: cfg.Logging(),
	}
}

func (s service) Run(ctx context.Context) error {
	s.log.Info("Staring twitter bot service...")

	go func() {
		s.log.Info("Staring authenticator bot service...")

		auth := authenticator.New(s.cfg)
		if err := auth.Listen(ctx); err != nil {
			// for now we are not stopping everything
			// stopping just listener would not affect the poster
			s.log.WithError(err).Error("Failed to listen to bot events!")
		}
	}()

	s.log.Info("Staring twitter poster bot service...")
	pst := poster.New(s.cfg)

	err := common.RunEvery(15*time.Minute, func() error {
		s.log.Debug("Posting news...")

		n, err := pst.Post(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to post")
		}

		s.log.WithField("news-posted", n).Debug("Finished posting")
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to run poster")
	}

	return nil
}
