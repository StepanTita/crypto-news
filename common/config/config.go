package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config interface {
	Logger
	Databaser
	Runtime
}

type config struct {
	Logger
	Databaser
	Runtime
}

type yamlConfig struct {
	LogLevel string             `yaml:"log_level"`
	Database YamlDatabaseConfig `yaml:"database"`
}

func NewFromFile(path string) Config {
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
		Logger:    NewLogger(cfg.LogLevel),
		Databaser: NewDatabaser(cfg.Database),
	}
}

func New(logLevel string, runtime YamlRuntimeConfig, database YamlDatabaseConfig) Config {
	return &config{
		Logger:    NewLogger(logLevel),
		Runtime:   NewRuntime(runtime.Environment, runtime.Version),
		Databaser: NewDatabaser(database),
	}
}
