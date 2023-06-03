package config

import (
	"os"

	gptconfig "github.com/StepanTita/go-EdgeGPT/config"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	commoncfg "common/config"
)

type Config interface {
	commoncfg.Config
	GPTConfig() gptconfig.Config
}

type config struct {
	commoncfg.Config
	gptCfg gptconfig.Config
}

func (c config) GPTConfig() gptconfig.Config {
	return c.gptCfg
}

type yamlConfig struct {
	LogLevel  string                       `yaml:"log_level"`
	Database  commoncfg.YamlDatabaseConfig `yaml:"database"`
	KVStore   commoncfg.YamlKVStoreConfig  `yaml:"kv_store"`
	Runtime   commoncfg.YamlRuntimeConfig  `yaml:"runtime"`
	GPTConfig gptconfig.YamlGPTConfig      `yaml:"gpt"`
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
		Config: commoncfg.New(cfg.LogLevel, cfg.Runtime, cfg.Database, cfg.KVStore),
		gptCfg: gptconfig.NewFromGPTConfig(cfg.GPTConfig),
	}
}
