package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	commoncfg "common/config"
)

type Config interface {
	commoncfg.Config
	Templator
	Listener
}

type config struct {
	commoncfg.Config
	Templator
	Listener
}

type yamlConfig struct {
	LogLevel         string                       `yaml:"log_level"`
	TemplatesDir     string                       `yaml:"templates_dir"`
	TelegramApiToken string                       `yaml:"telegram_api_token"`
	Database         commoncfg.YamlDatabaseConfig `yaml:"database"`
	Runtime          commoncfg.YamlRuntimeConfig  `yaml:"runtime"`
}

func New(path string) Config {
	cfg := yamlConfig{}

	yamlConfig, err := os.ReadFile(path)
	if err != nil {
		panic(errors.Wrapf(err, "failed to read config %s", path))
	}

	err = yaml.Unmarshal(yamlConfig, &cfg)
	if err != nil {
		panic(errors.Wrapf(err, "failed to unmarshal config %s", path))
	}

	return &config{
		Config:    commoncfg.New(cfg.LogLevel, cfg.Runtime, cfg.Database),
		Templator: NewTemplator(cfg.TemplatesDir),
		Listener:  NewListener(cfg.TelegramApiToken),
	}
}
