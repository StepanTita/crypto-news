package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	commoncfg "common/config"
)

type Config interface {
	commoncfg.Config
	Crawler
	ServiceProvider
}

type config struct {
	commoncfg.Config
	Crawler
	ServiceProvider
}

type yamlConfig struct {
	LogLevel         string                       `yaml:"log_level"`
	RateLimit        int                          `yaml:"rate_limit"`
	CrawlEvery       time.Duration                `yaml:"crawl_every"`
	Database         commoncfg.YamlDatabaseConfig `yaml:"database"`
	KVStore          commoncfg.YamlKVStoreConfig  `yaml:"kv_store"`
	ServiceProviders yamlServiceProviderConfig    `yaml:"service_providers"`
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
		Config:          commoncfg.New(cfg.LogLevel, cfg.Runtime, cfg.Database, cfg.KVStore),
		Crawler:         NewCrawler(cfg.RateLimit, cfg.CrawlEvery),
		ServiceProvider: NewServiceProvider(cfg.ServiceProviders),
	}
}
