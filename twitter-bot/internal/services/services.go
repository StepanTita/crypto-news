package services

import (
	"context"

	"github.com/sirupsen/logrus"

	"twitter-bot/internal/config"
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

	return nil
}
