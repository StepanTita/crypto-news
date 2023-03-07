package config

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Logging() *logrus.Entry
}

type logger struct {
	log *logrus.Entry
}

func NewLogger(lvl string) Logger {
	log := logrus.New()
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		panic(errors.Wrapf(err, "failed to parse log level: %s", lvl))
	}
	log.SetLevel(level)

	timestampFormatter := new(logrus.TextFormatter)
	timestampFormatter.FullTimestamp = true
	log.SetFormatter(timestampFormatter)

	return &logger{
		log: logrus.NewEntry(log),
	}
}

func (l *logger) Logging() *logrus.Entry {
	return l.log
}
