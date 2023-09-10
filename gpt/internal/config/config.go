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
	Generator
	BotConfig
}

type config struct {
	commoncfg.Config
	Generator
	BotConfig
}

type yamlConfig struct {
	LogLevel  string                       `yaml:"log_level"`
	Database  commoncfg.YamlDatabaseConfig `yaml:"database"`
	KVStore   commoncfg.YamlKVStoreConfig  `yaml:"kv_store"`
	Runtime   commoncfg.YamlRuntimeConfig  `yaml:"runtime"`
	GPTConfig struct {
		AuthToken     string        `yaml:"auth_token"`
		GenerateEvery time.Duration `yaml:"generate_every"`
		QueryContext  string        `yaml:"query_context"`
		Prompt        string        `yaml:"prompt"`
		ImagesPrompt  string        `yaml:"images_prompt"`
	} `yaml:"gpt"`
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
		Config:    commoncfg.New(cfg.LogLevel, cfg.Runtime, cfg.Database, cfg.KVStore),
		BotConfig: NewBotConfig(cfg.GPTConfig.AuthToken),
		Generator: NewGenerator(cfg.GPTConfig.GenerateEvery, cfg.GPTConfig.ImagesPrompt, cfg.GPTConfig.QueryContext),
	}
}
