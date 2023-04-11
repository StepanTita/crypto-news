package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	commoncfg "common/config"
)

type Config interface {
	commoncfg.Config
	Twitter
}

type config struct {
	commoncfg.Config
	Twitter
}

type yamlConfig struct {
	LogLevel string                       `yaml:"log_level"`
	Twitter  YamlTwitterConfig            `yaml:"twitter"`
	Database commoncfg.YamlDatabaseConfig `yaml:"database"`
	KVStore  commoncfg.YamlKVStoreConfig  `yaml:"kv_store"`
	Runtime  commoncfg.YamlRuntimeConfig  `yaml:"runtime"`
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
		Config:  commoncfg.New(cfg.LogLevel, cfg.Runtime, cfg.Database, cfg.KVStore),
		Twitter: NewTwitter(cfg.Twitter),
	}
}
