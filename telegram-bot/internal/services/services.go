package services

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common"
	"telegram-bot/internal/config"
	"telegram-bot/internal/services/listener"
	"telegram-bot/internal/services/poster"
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
	s.log.Info("Staring telegram bot service...")

	bot, err := tgbotapi.NewBotAPI(s.cfg.TelegramApiToken())
	if err != nil {
		return errors.Wrap(err, "failed to initialize bot API")
	}

	go func() {
		s.log.Info("Staring listening bot service...")
		lst := listener.New(s.cfg, bot)
		if err := lst.Listen(ctx); err != nil {
			// for now we are not stopping everything
			// stopping just listener would not affect the poster
			s.log.WithError(err).Error("Failed to listen to bot events!")
		}
	}()

	s.log.Info("Staring poster bot service...")
	pst := poster.New(s.cfg, bot)

	err = common.RunEvery(15*time.Second, func() error {
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
