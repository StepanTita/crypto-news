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
}

type config struct {
	Logger
	Databaser
	KVStorer
	Runtime
	Templator
}

type yamlConfig struct {
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
	}
}

func New(logLevel string, runtime YamlRuntimeConfig, database YamlDatabaseConfig, kvStore YamlKVStoreConfig) Config {
	return &config{
		Logger:    NewLogger(logLevel),
		Runtime:   NewRuntime(runtime.Environment, runtime.Version),
		Databaser: NewDatabaser(database),
		KVStorer:  NewKVStorer(kvStore),
		Templator: NewTemplator(),
	}
}
