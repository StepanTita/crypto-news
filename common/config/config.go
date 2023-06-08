package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config interface {
	Logger
	Databaser
	KVStorer
	Runtime
	Templator
	Localizer
}

type config struct {
	Logger
	Databaser
	KVStorer
	Runtime
	Templator
	Localizer
}

type yamlConfig struct {
	Environment string   `yaml:"environment"`
	Version     string   `yaml:"version"`
	Locales     []string `yaml:"locales"`

	LogLevel string             `yaml:"log_level"`
	Database YamlDatabaseConfig `yaml:"database"`
	KVStore  YamlKVStoreConfig  `yaml:"kv_store"`
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
		KVStorer:  NewKVStorer(cfg.KVStore),
		Runtime:   NewRuntime(cfg.Environment, cfg.Version, cfg.Locales),
	}
}

func New(logLevel string, runtime YamlRuntimeConfig, database YamlDatabaseConfig, kvStore YamlKVStoreConfig) Config {
	return &config{
		Logger:    NewLogger(logLevel),
		Runtime:   NewRuntime(runtime.Environment, runtime.Version, runtime.Locales),
		Databaser: NewDatabaser(database),
		KVStorer:  NewKVStorer(kvStore),
		Templator: NewTemplator(),
		Localizer: NewLocalizer(),
	}
}
