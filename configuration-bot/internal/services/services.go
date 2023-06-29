package services

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"configuration-bot/internal/config"
	"configuration-bot/internal/services/listener"
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

	s.log.Info("Staring listening bot service...")
	lst := listener.New(s.cfg, bot)
	if err := lst.Listen(ctx); err != nil {
		// for now we are not stopping everything
		// stopping just listener would not affect the poster
		s.log.WithError(err).Error("Failed to listen to bot events!")
	}

	return nil
}
